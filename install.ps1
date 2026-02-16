$ErrorActionPreference = "Stop"

$REPO_URL = "https://github.com/KidiXDev/gofiber-v3-starterkit.git"
$BRANCH = "main"

Clear-Host
Write-Host "GoFiber V3 Starter Pack Installer" -ForegroundColor Cyan
Write-Host "---------------------------------" -ForegroundColor Gray

if (-not (Get-Command git -ErrorAction SilentlyContinue)) {
    Write-Host "Error: git is not installed." -ForegroundColor Red
    exit 1
}

$PROJECT_NAME = Read-Host "Enter project name (default: my-gofiber-app)"
if ([string]::IsNullOrWhiteSpace($PROJECT_NAME)) {
    $PROJECT_NAME = "my-gofiber-app"
}

if (Test-Path $PROJECT_NAME) {
    Write-Host "Error: Directory '$PROJECT_NAME' already exists." -ForegroundColor Red
    exit 1
}

Write-Host "Cloning repository into '$PROJECT_NAME'..." -ForegroundColor Cyan
git clone --depth 1 $REPO_URL $PROJECT_NAME

if ($LASTEXITCODE -ne 0) {
    Write-Host "Error: Failed to clone repository." -ForegroundColor Red
    exit 1
}

Set-Location $PROJECT_NAME

Remove-Item -Path .git -Recurse -Force
git init

Remove-Item install.sh, install.ps1 -ErrorAction SilentlyContinue

Write-Host "Repository cloned successfully!" -ForegroundColor Green
Write-Host ""

.\rename-module.bat

if ($LASTEXITCODE -eq 0) {
    Remove-Item rename-module.bat, rename-module.sh -ErrorAction SilentlyContinue
}

git add .
git commit -m "initial commit"
