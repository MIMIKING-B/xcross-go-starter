# =============================================================================
# 项目核心构建命令
# =============================================================================
# 作用：
#   本文件定义了项目开发、构建、部署的全生命周期命令。
#   这些命令都基于 GoFrame CLI (gf) 实现，因此每个目标都依赖 cli.install
#   来确保 gf 工具已安装。
#
# 为什么单独放在 hack.mk？
#   1. 与根 Makefile 解耦，根 Makefile 只负责变量定义和文件引入
#   2. 便于在不同项目中复用相同的构建逻辑
#   3. 使命令分类清晰（开发命令 vs CLI 管理命令）
#
# 命令速查表：
#   开发命令  ：make build / make ctrl / make dao / make service / make enums
#   构建命令  ：make image / make image.push
#   部署命令  ：make deploy
#   协议命令  ：make pb / make pbentity
#   升级命令  ：make up
# =============================================================================

# 默认目标：直接执行 make 等同于 make build
.DEFAULT_GOAL := build


# -----------------------------------------------------------------------------
# 升级 GoFrame 框架及其 CLI 到最新稳定版
# 原理：gf up -a 会自动检测并更新 go.mod 中的 gf 依赖到最新版本
# 使用场景：定期同步 GoFrame 框架的安全补丁和新功能
# -----------------------------------------------------------------------------
.PHONY: up
up: cli.install
	@gf up -a


# -----------------------------------------------------------------------------
# 编译项目二进制文件
# 原理：读取 hack/config.yaml 中的 build 配置，交叉编译为可执行文件
# 参数 -ew：启用额外功能（如自动打包 resource 目录到可执行文件中）
#
# 使用场景：
#   1. 本地开发时编译测试
#   2. CI/CD 流水线中编译生产版本
#   3. 打包静态资源（HTML/CSS/JS）到二进制中，实现单文件部署
# -----------------------------------------------------------------------------
.PHONY: build
build: cli.install
	@gf build -ew


# -----------------------------------------------------------------------------
# 根据 api 目录下的接口定义生成 Controller 代码
# 原理：
#   1. 扫描 api/ 目录下所有带 g.Meta 标签的结构体（定义了路由、方法、参数）
#   2. 自动在 internal/controller/ 下生成对应的控制器文件和接口方法
#
# 使用场景：
#   定义了新的 API 接口后，运行此命令自动生成 Controller 骨架代码，
#   开发者只需在生成的文件中填充业务逻辑即可。
# -----------------------------------------------------------------------------
.PHONY: ctrl
ctrl: cli.install
	@gf gen ctrl


# -----------------------------------------------------------------------------
# 根据数据库表结构生成 DAO/DO/Entity 代码
# 原理：
#   1. 连接数据库（配置在 hack/config.yaml 的 gen.dao.link）
#   2. 读取所有表结构，为每张表生成对应的 DAO（数据访问对象）
#   3. 同时生成 DO（数据库操作对象）和 Entity（实体结构体）
#
# 使用场景：
#   数据库表结构变更后（新增表、修改字段），运行此命令同步更新 Go 代码，
#   无需手写繁琐的数据库模型代码。
# -----------------------------------------------------------------------------
.PHONY: dao
dao: cli.install
	@gf gen dao


# -----------------------------------------------------------------------------
# 解析项目中的 Go 文件并生成枚举定义
# 原理：扫描代码中带有特定注释的常量定义，自动生成枚举相关代码
#
# 使用场景：
#   定义了状态码、类型常量等枚举值后，自动生成枚举的字符串映射、
#   验证方法等辅助代码。
# -----------------------------------------------------------------------------
.PHONY: enums
enums: cli.install
	@gf gen enums


# -----------------------------------------------------------------------------
# 根据 internal/logic/ 下的实现自动生成 Service 接口
# 原理：
#   1. 扫描 internal/logic/ 目录下所有 service.RegisterXxx() 调用
#   2. 在 internal/service/ 目录下生成对应的接口定义文件
#   3. 自动处理接口方法的参数和返回值
#
# 使用场景：
#   在 logic 层实现了新的业务方法后，运行此命令自动同步更新 service 接口，
#   保持接口与实现的一致性。这是 GoFrame 服务层规范的核心支撑。
# -----------------------------------------------------------------------------
.PHONY: service
service: cli.install
	@gf gen service


# -----------------------------------------------------------------------------
# 构建 Docker 镜像
# 原理：
#   1. 读取 hack/config.yaml 中的 docker 配置
#   2. 自动编译 Linux amd64 架构的二进制
#   3. 基于 manifest/docker/Dockerfile 构建镜像
#
# 镜像标签规则：
#   默认使用当前 git commit 的短哈希作为标签（如 abc1234）
#   如果有未提交的修改，标签会加上 .dirty 后缀（如 abc1234.dirty）
#   可通过 TAG 参数自定义标签：make image TAG=v1.0.0
#
# 使用场景：本地构建镜像后推送到镜像仓库，供 K8s 部署使用
# -----------------------------------------------------------------------------
.PHONY: image
image: cli.install
	$(eval _TAG  = $(shell git rev-parse --short HEAD))
ifneq (, $(shell git status --porcelain 2>/dev/null))
	$(eval _TAG  = $(_TAG).dirty)
endif
	$(eval _TAG  = $(if ${TAG},  ${TAG}, $(_TAG)))
	$(eval _PUSH = $(if ${PUSH}, ${PUSH}, ))
	@gf docker ${_PUSH} -tn $(DOCKER_NAME):${_TAG};


# -----------------------------------------------------------------------------
# 构建 Docker 镜像并自动推送到镜像仓库
# 原理：调用上方的 make image，并传入 PUSH=-p 参数触发自动推送
#
# 使用场景：CI/CD 流水线中一键构建并推送镜像
# -----------------------------------------------------------------------------
.PHONY: image.push
image.push: cli.install
	@make image PUSH=-p;


# -----------------------------------------------------------------------------
# 部署到当前 kubectl 环境的 Kubernetes 集群
# 原理：
#   1. 使用 kustomize 合并基础配置和覆盖层配置（如 develop/production）
#   2. 生成最终的 K8s YAML 并应用到集群
#   3. 通过 kubectl patch 触发 Deployment 滚动更新（更新时间戳标签）
#
# 环境变量 _ENV：
#   指定使用哪个覆盖层配置，如 develop、staging、product
#   对应目录：manifest/deploy/kustomize/overlays/${_ENV}/
#
# 使用场景：
#   本地联调时部署到开发集群，或 CI/CD 中自动部署到生产集群
# -----------------------------------------------------------------------------
.PHONY: deploy
deploy: cli.install
	$(eval _TAG = $(if ${TAG},  ${TAG}, develop))

	@set -e; \
	mkdir -p $(ROOT_DIR)/temp/kustomize;\
	cd $(ROOT_DIR)/manifest/deploy/kustomize/overlays/${_ENV};\
	kustomize build > $(ROOT_DIR)/temp/kustomize.yaml;\
	kubectl   apply -f $(ROOT_DIR)/temp/kustomize.yaml; \
	if [ $(DEPLOY_NAME) != "" ]; then \
		kubectl patch -n $(NAMESPACE) deployment/$(DEPLOY_NAME) -p "{\"spec\":{\"template\":{\"metadata\":{\"labels\":{\"date\":\"$(shell date +%s)\"}}}}}"; \
	fi;


# -----------------------------------------------------------------------------
# 解析 protobuf 文件并生成 Go 代码
# 原理：扫描 manifest/protobuf/ 目录下的 .proto 文件，生成对应的 .pb.go
#
# 使用场景：微服务间使用 gRPC 通信时，根据 proto 定义生成客户端/服务端代码
# -----------------------------------------------------------------------------
.PHONY: pb
pb: cli.install
	@gf gen pb


# -----------------------------------------------------------------------------
# 根据数据库表结构生成 protobuf 文件
# 原理：连接数据库，读取表结构并转换为 .proto 消息定义
#
# 使用场景：
#   需要基于现有数据库表生成 gRPC 的 Protocol Buffers 定义时，
#   避免手动编写与表结构对应的 proto 消息。
# -----------------------------------------------------------------------------
.PHONY: pbentity
pbentity: cli.install
	@gf gen pbentity