#!/bin/bash

# ============================================================================
# GoFrame Project Template Initializer
# ============================================================================
# Usage:
#   ./hack/init-project.sh <new-module-name> [output-directory]
#
# Examples:
#   ./hack/init-project.sh github.com/company/mynewproject
#   ./hack/init-project.sh github.com/company/mynewproject ./mynewproject
#   make init-project name=github.com/company/mynewproject
# ============================================================================

set -e

OLD_MODULE="xcross-go-starter"
NEW_MODULE="${1:-}"
OUTPUT_DIR="${2:-}"

# ----------------------------------------------------------------------------
# 定位模板根目录（脚本在 hack/ 下，模板根目录是脚本所在目录的父目录）
# ----------------------------------------------------------------------------
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
TEMPLATE_DIR="$(dirname "$SCRIPT_DIR")"

# ----------------------------------------------------------------------------
# 参数校验
# ----------------------------------------------------------------------------
if [ -z "$NEW_MODULE" ]; then
    echo ""
    echo "Error: Missing new module name."
    echo ""
    echo "Usage: $0 <new-module-name> [output-directory]"
    echo ""
    echo "Examples:"
    echo "  $0 github.com/company/mynewproject"
    echo "  $0 github.com/company/mynewproject ./mynewproject"
    echo ""
    exit 1
fi

# 推断输出目录
if [ -z "$OUTPUT_DIR" ]; then
    OUTPUT_DIR=$(basename "$NEW_MODULE")
fi

# 推断应用名（取模块名最后一段）
APP_NAME=$(basename "$NEW_MODULE")

# ----------------------------------------------------------------------------
# 预检查
# ----------------------------------------------------------------------------
if [ -e "$OUTPUT_DIR" ]; then
    echo "Error: Directory '$OUTPUT_DIR' already exists!"
    exit 1
fi

if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed or not in PATH."
    exit 1
fi

# 跨平台 sed 替换函数（兼容 macOS 和 Linux）
sed_replace() {
    local pattern="$1"
    local file="$2"
    if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -i '' "$pattern" "$file"
    else
        sed -i "$pattern" "$file"
    fi
}

# ----------------------------------------------------------------------------
# 开始初始化
# ----------------------------------------------------------------------------
echo ""
echo "=========================================="
echo "  GoFrame Project Template Initializer"
echo "=========================================="
echo "  Old Module : $OLD_MODULE"
echo "  New Module : $NEW_MODULE"
echo "  App Name   : $APP_NAME"
echo "  Output Dir : $OUTPUT_DIR"
echo "=========================================="
echo ""

# [1/7] 复制模板（排除不需要的文件）
echo "[1/7] Copying template..."
mkdir -p "$OUTPUT_DIR"

cp -r "$TEMPLATE_DIR"/. "$OUTPUT_DIR/"
cd "$OUTPUT_DIR"

# 删除不需要的文件和目录
rm -rf .git .idea logs temp go.sum
find storage/cache -mindepth 1 -delete 2>/dev/null || true
find storage/dev -mindepth 1 -delete 2>/dev/null || true
find . -name '*.iml' -delete 2>/dev/null || true
rm -f hack/init-project.sh

# 确保 storage/cache 和 storage/dev 目录结构保留
mkdir -p storage/cache storage/dev

# [2/7] 替换 go.mod 模块名
echo "[2/7] Updating go.mod..."
sed_replace "s|^module ${OLD_MODULE}$|module ${NEW_MODULE}|g" go.mod

# [3/7] 替换所有 Go 文件中的 import 路径
echo "[3/7] Replacing import paths in Go files..."
while IFS= read -r -d '' file; do
    sed_replace "s|\"${OLD_MODULE}/|\"${NEW_MODULE}/|g" "$file"
done < <(find . -type f -name "*.go" -print0)

# [4/7] 替换配置文件中应用名称相关的内容
echo "[4/7] Updating configuration files..."
sed_replace "s|appName: \"${OLD_MODULE}\"|appName: \"${APP_NAME}\"|g" manifest/config/config.yaml
sed_replace "s|title: \"${OLD_MODULE}\"|title: \"${APP_NAME}\"|g" manifest/config/config.yaml
sed_replace "s|keywords: \"${OLD_MODULE}\"|keywords: \"${APP_NAME}\"|g" manifest/config/config.yaml
sed_replace "s|description: \"${OLD_MODULE}\"|description: \"${APP_NAME}\"|g" manifest/config/config.yaml

if [ -f "README.md" ]; then
    sed_replace "s|${OLD_MODULE}|${APP_NAME}|g" README.md
fi

if [ -f "Makefile" ]; then
    sed_replace "s|${OLD_MODULE}|${APP_NAME}|g" Makefile
fi

if [ -f "manifest/docker/Dockerfile" ]; then
    sed_replace "s|${OLD_MODULE}|${APP_NAME}|g" manifest/docker/Dockerfile
fi

# hack/config.yaml 中的构建配置名和输出路径
if [ -f "hack/config.yaml" ]; then
    sed_replace "s|name: \"${OLD_MODULE}\"|name: \"${APP_NAME}\"|g" hack/config.yaml
    sed_replace "s|output: \"./temp/${OLD_MODULE}\"|output: \"./temp/${APP_NAME}\"|g" hack/config.yaml
fi

# [5/7] 清理模板痕迹
echo "[5/7] Cleaning template artifacts..."
rm -f go.sum
find logs -type f -delete 2>/dev/null || true
find temp -type f -delete 2>/dev/null || true

# [6/7] 重新初始化 Go module
echo "[6/7] Running go mod tidy..."
go mod tidy

# [7/7] 初始化 Git
echo "[7/7] Initializing git..."
rm -rf .git
git init -q
git add .
git commit -q -m "chore: init project from template ${OLD_MODULE} -> ${NEW_MODULE}"

# ----------------------------------------------------------------------------
# 完成
# ----------------------------------------------------------------------------
echo ""
echo "=========================================="
echo "  Project initialized successfully!"
echo "=========================================="
echo ""
echo "Next steps:"
echo "  cd ${OUTPUT_DIR}"
echo "  go run main.go"
echo "  # or: gf run main.go"
echo ""
echo "You may also want to:"
echo "  1. Update manifest/config/config.yaml with your DB/Redis settings"
echo "  2. Update manifest/docker/Dockerfile if needed"
echo "  3. Update manifest/deploy/kustomize/ for your K8s deployment"
echo "  4. Edit README.md to match your project description"
echo ""
