# ============================================================================
# GoFrame Project Template Initializer (Windows PowerShell 版本)
# ============================================================================
# 作用：
#   基于当前脚手架模板创建一个新项目，自动完成模块名替换、import 路径更新、
#   配置项替换、Git 初始化等操作。
#
# 这是 init-project.sh 的 Windows 版本，功能完全一致。
#
# 使用方法：
#   .\hack\init-project.ps1 -ModuleName <新模块名> [-OutputDir <输出目录>]
#
# 示例：
#   .\hack\init-project.ps1 -ModuleName github.com/company/mynewproject
#   .\hack\init-project.ps1 -ModuleName github.com/company/mynewproject -OutputDir .\mynewproject
# ============================================================================

param(
    [Parameter(Mandatory = $true, HelpMessage = "新项目的 Go 模块名，例如 github.com/company/app")]
    [string]$ModuleName,

    [Parameter(Mandatory = $false, HelpMessage = "输出目录，默认为模块名的最后一段")]
    [string]$OutputDir = ""
)

# ----------------------------------------------------------------------------
# 配置
# ----------------------------------------------------------------------------
$OLD_MODULE = "xcross-go-starter"
$NEW_MODULE = $ModuleName

# 推断输出目录
if ([string]::IsNullOrWhiteSpace($OutputDir)) {
    $OutputDir = Split-Path -Leaf $ModuleName
}

# 推断应用名（取模块名最后一段）
$APP_NAME = Split-Path -Leaf $ModuleName

# 推断模板根目录（脚本在 hack/ 下，模板根目录是父目录）
$SCRIPT_DIR = Split-Path -Parent $MyInvocation.MyCommand.Path
$TEMPLATE_DIR = Split-Path -Parent $SCRIPT_DIR

# ----------------------------------------------------------------------------
# 预检查
# ----------------------------------------------------------------------------
if (Test-Path $OutputDir) {
    Write-Host "错误：目录 '$OutputDir' 已存在！" -ForegroundColor Red
    exit 1
}

$goCmd = Get-Command go -ErrorAction SilentlyContinue
if (-not $goCmd) {
    Write-Host "错误：未找到 Go，请确保 Go 已安装并添加到 PATH。" -ForegroundColor Red
    exit 1
}

# ----------------------------------------------------------------------------
# 开始初始化
# ----------------------------------------------------------------------------
Write-Host ""
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "  GoFrame Project Template Initializer" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "  Old Module : $OLD_MODULE"
Write-Host "  New Module : $NEW_MODULE"
Write-Host "  App Name   : $APP_NAME"
Write-Host "  Output Dir : $OutputDir"
Write-Host "=========================================="
Write-Host ""

# [1/7] 复制模板（排除不需要的文件）
Write-Host "[1/7] Copying template..."
New-Item -ItemType Directory -Path $OutputDir | Out-Null

# 使用 robocopy 或 Copy-Item 复制，然后删除不需要的文件
Copy-Item -Path "$TEMPLATE_DIR\*" -Destination $OutputDir -Recurse -Force

# 删除不需要的文件和目录
$excludeItems = @(
    "$OutputDir\.git",
    "$OutputDir\.idea",
    "$OutputDir\logs",
    "$OutputDir\temp",
    "$OutputDir\go.sum"
)
foreach ($item in $excludeItems) {
    if (Test-Path $item) {
        Remove-Item -Path $item -Recurse -Force -ErrorAction SilentlyContinue
    }
}

# 删除 storage/cache 和 storage/dev 下的内容（保留目录）
$storageDirs = @("$OutputDir\storage\cache", "$OutputDir\storage\dev")
foreach ($dir in $storageDirs) {
    if (Test-Path $dir) {
        Get-ChildItem -Path $dir -Recurse | Remove-Item -Recurse -Force -ErrorAction SilentlyContinue
    }
}

# 删除 .iml 文件
Get-ChildItem -Path $OutputDir -Filter "*.iml" -Recurse | Remove-Item -Force -ErrorAction SilentlyContinue

# 删除 init-project 脚本本身（它们不应该出现在新项目中）
Remove-Item -Path "$OutputDir\hack\init-project.sh" -Force -ErrorAction SilentlyContinue
Remove-Item -Path "$OutputDir\hack\init-project.ps1" -Force -ErrorAction SilentlyContinue

# 确保 storage/cache 和 storage/dev 目录结构保留
New-Item -ItemType Directory -Path "$OutputDir\storage\cache" -Force | Out-Null
New-Item -ItemType Directory -Path "$OutputDir\storage\dev" -Force | Out-Null

Set-Location $OutputDir

# [2/7] 替换 go.mod 模块名
Write-Host "[2/7] Updating go.mod..."
$goModContent = Get-Content go.mod -Raw
$goModContent = $goModContent -replace "^module $([regex]::Escape($OLD_MODULE))$", "module $NEW_MODULE"
Set-Content go.mod $goModContent -NoNewline

# [3/7] 替换所有 Go 文件中的 import 路径
Write-Host "[3/7] Replacing import paths in Go files..."
$goFiles = Get-ChildItem -Path . -Filter "*.go" -Recurse
foreach ($file in $goFiles) {
    $content = Get-Content $file.FullName -Raw
    $content = $content -replace "`"$([regex]::Escape($OLD_MODULE))\/", "`"$NEW_MODULE/"
    Set-Content $file.FullName $content -NoNewline
}

# [4/7] 替换配置文件中应用名称相关的内容
Write-Host "[4/7] Updating configuration files..."

# manifest/config/config.yaml
$configFile = "manifest\config\config.yaml"
if (Test-Path $configFile) {
    $content = Get-Content $configFile -Raw
    $content = $content -replace "appName: `"$([regex]::Escape($OLD_MODULE))`"", "appName: `"$APP_NAME`""
    $content = $content -replace "title: `"$([regex]::Escape($OLD_MODULE))`"", "title: `"$APP_NAME`""
    $content = $content -replace "keywords: `"$([regex]::Escape($OLD_MODULE))`"", "keywords: `"$APP_NAME`""
    $content = $content -replace "description: `"$([regex]::Escape($OLD_MODULE))`"", "description: `"$APP_NAME`""
    Set-Content $configFile $content -NoNewline
}

# hack/config.yaml
$hackConfigFile = "hack\config.yaml"
if (Test-Path $hackConfigFile) {
    $content = Get-Content $hackConfigFile -Raw
    $content = $content -replace "name: `"$([regex]::Escape($OLD_MODULE))`"", "name: `"$APP_NAME`""
    $content = $content -replace "output: `"\.\/temp\/$([regex]::Escape($OLD_MODULE))`"", "output: `"./temp/$APP_NAME`""
    Set-Content $hackConfigFile $content -NoNewline
}

# README.md
if (Test-Path "README.md") {
    $content = Get-Content "README.md" -Raw
    $content = $content -replace [regex]::Escape($OLD_MODULE), $APP_NAME
    Set-Content "README.md" $content -NoNewline
}

# Makefile
if (Test-Path "Makefile") {
    $content = Get-Content "Makefile" -Raw
    $content = $content -replace [regex]::Escape($OLD_MODULE), $APP_NAME
    Set-Content "Makefile" $content -NoNewline
}

# Dockerfile
$dockerFile = "manifest\docker\Dockerfile"
if (Test-Path $dockerFile) {
    $content = Get-Content $dockerFile -Raw
    $content = $content -replace [regex]::Escape($OLD_MODULE), $APP_NAME
    Set-Content $dockerFile $content -NoNewline
}

# [5/7] 清理模板痕迹
Write-Host "[5/7] Cleaning template artifacts..."
Remove-Item go.sum -Force -ErrorAction SilentlyContinue

# 清理日志文件（保留目录）
$logDirs = @("logs", "temp")
foreach ($dir in $logDirs) {
    if (Test-Path $dir) {
        Get-ChildItem -Path $dir -Recurse -File | Remove-Item -Force -ErrorAction SilentlyContinue
    }
}

# [6/7] 重新初始化 Go module
Write-Host "[6/7] Running go mod tidy..."
& go mod tidy
if ($LASTEXITCODE -ne 0) {
    Write-Host "警告：go mod tidy 执行失败，可能需要手动运行。" -ForegroundColor Yellow
}

# [7/7] 初始化 Git
Write-Host "[7/7] Initializing git..."
if (Test-Path .git) {
    Remove-Item .git -Recurse -Force
}
& git init -q
& git add .
& git commit -q -m "chore: init project from template $OLD_MODULE -> $NEW_MODULE"

# ----------------------------------------------------------------------------
# 完成
# ----------------------------------------------------------------------------
Write-Host ""
Write-Host "==========================================" -ForegroundColor Green
Write-Host "  Project initialized successfully!" -ForegroundColor Green
Write-Host "==========================================" -ForegroundColor Green
Write-Host ""
Write-Host "Next steps:"
Write-Host "  cd $OutputDir"
Write-Host "  go run main.go"
Write-Host "  # or: gf run main.go"
Write-Host ""
Write-Host "You may also want to:"
Write-Host "  1. Update manifest/config/config.yaml with your DB/Redis settings"
Write-Host "  2. Update manifest/docker/Dockerfile if needed"
Write-Host "  3. Update manifest/deploy/kustomize/ for your K8s deployment"
Write-Host "  4. Edit README.md to match your project description"
Write-Host ""
