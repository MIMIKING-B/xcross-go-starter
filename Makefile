# =============================================================================
# 项目根 Makefile
# =============================================================================
# 作用：
#   这是整个项目的构建入口文件。Makefile 是 Linux/macOS/Windows(需make工具)
#   下的自动化构建脚本，通过简单的命令即可执行复杂的编译、生成、部署等操作。
#
# 使用场景：
#   1. 开发阶段：快速编译、热重载、生成代码（dao/service/controller）
#   2. 构建阶段：编译二进制、构建 Docker 镜像
#   3. 部署阶段：推送镜像、部署到 Kubernetes
#   4. 初始化阶段：基于本模板创建新项目
#
# 使用方法：
#   make <命令>        例如：make build
#   make init-project name=github.com/company/myproject
# =============================================================================

# 项目根目录（自动获取当前路径）
ROOT_DIR    = $(shell pwd)

# Kubernetes 命名空间，用于 make deploy 命令部署时的目标命名空间
NAMESPACE   = "default"

# Kubernetes Deployment 名称，用于滚动更新时 patch 的部署对象名称
DEPLOY_NAME = "template-single"

# Docker 镜像名称，用于 make image 命令构建镜像时的镜像名
DOCKER_NAME = "template-single"

# =============================================================================
# 引入子 Makefile
# =============================================================================
# hack-cli.mk：GoFrame CLI 工具的安装与管理（gf 命令的自动安装/更新）
# hack.mk    ：项目核心构建命令（编译、生成代码、Docker、K8s 部署等）
# =============================================================================
include ./hack/hack-cli.mk
include ./hack/hack.mk

# =============================================================================
# 项目模板初始化
# =============================================================================
# 作用：
#   基于当前脚手架模板创建一个新项目，自动完成以下操作：
#   - 复制项目文件（排除 .git、日志、缓存等无关文件）
#   - 替换 go.mod 中的模块名
#   - 替换所有 Go 文件中的 import 路径
#   - 替换配置文件中的应用名称、标题等
#   - 清理模板痕迹，执行 go mod tidy
#   - 初始化新的 Git 仓库
#
# 跨平台支持：
#   - macOS / Linux：自动使用 bash 脚本 (hack/init-project.sh)
#   - Windows     ：自动使用 PowerShell 脚本 (hack/init-project.ps1)
#
# 使用方法：
#   make init-project name=<新模块名> [out=<输出目录>]
#
# 示例：
#   make init-project name=github.com/company/mynewproject
#   make init-project name=github.com/company/mynewproject out=./mynewproject
# =============================================================================
.PHONY: init-project
init-project:
	@if [ -z "$(name)" ]; then \
		echo ""; \
		echo "错误：缺少 'name' 参数"; \
		echo ""; \
		echo "用法：make init-project name=<新模块名> [out=<输出目录>]"; \
		echo ""; \
		echo "示例："; \
		echo "  make init-project name=github.com/company/mynewproject"; \
		echo "  make init-project name=github.com/company/mynewproject out=./mynewproject"; \
		echo ""; \
		exit 1; \
	fi
	@echo "正在检测操作系统..."
	@if [ "$$(uname)" = "Darwin" ] || [ "$$(uname)" = "Linux" ]; then \
		echo "检测到 Unix 系统，使用 bash 脚本..."; \
		bash ./hack/init-project.sh $(name) $(out); \
	else \
		echo "检测到 Windows 系统，使用 PowerShell 脚本..."; \
		powershell -ExecutionPolicy Bypass -File ./hack/init-project.ps1 -ModuleName $(name) -OutputDir $(out); \
	fi