@echo off
setlocal EnableDelayedExpansion

cls
echo ========================================================
echo   GoFiber V3 Starter Pack Wizard
echo ========================================================
echo.

set "OLD_MODULE=gofiber-starterkit"

:: Check if module name is provided as argument
if "%~1"=="" (
    echo Please enter your new module name (e.g., github.com/username/project):
    set /p NEW_MODULE="> "
) else (
    set "NEW_MODULE=%~1"
)

if "%NEW_MODULE%"=="" (
    echo Error: Module name cannot be empty.
    goto :End
)

echo.
echo You are about to rename the module from:
echo [ %OLD_MODULE% ] -^> [ %NEW_MODULE% ]
echo.
set /p CONFIRM="Are you sure? (y/n): "
if /i not "!CONFIRM!"=="y" (
    echo Operation cancelled.
    exit /b 1
)

echo.
echo Renaming module...

for /r %%f in (*.go) do (
    powershell -Command "(Get-Content '%%f') -replace '%OLD_MODULE%', '%NEW_MODULE%' | Set-Content '%%f'"
)

powershell -Command "(Get-Content 'go.mod') -replace 'module %OLD_MODULE%', 'module %NEW_MODULE%' | Set-Content 'go.mod'"

powershell -Command "(Get-Content 'rename-module.sh') -replace 'OLD_MODULE=\"%OLD_MODULE%\"', 'OLD_MODULE=\"%NEW_MODULE%\"' | Set-Content 'rename-module.sh'"
powershell -Command "(Get-Content 'rename-module.bat') -replace 'OLD_MODULE=%OLD_MODULE%', 'OLD_MODULE=%NEW_MODULE%' | Set-Content 'rename-module.bat'"

echo.
echo Module renamed successfully!
echo.
echo Next steps:
echo 1. Run 'go mod tidy' to update dependencies
echo 2. Run 'go build' to verify the build
echo 3. Copy .env.example to .env and configure your environment
echo 4. Run 'go run .' to start the server

:End
endlocal
