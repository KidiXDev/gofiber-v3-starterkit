@echo off
setlocal EnableDelayedExpansion

if "%~1"=="" (
    echo Usage: rename-module.bat ^<new-module-name^>
    echo Example: rename-module.bat github.com/username/my-project
    exit /b 1
)

set "OLD_MODULE=gofiber-starterkit"
set "NEW_MODULE=%~1"

echo Renaming module from '%OLD_MODULE%' to '%NEW_MODULE%'...

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
echo 4. Run the SQL migrations in migrations/ folder
echo 5. Run 'go run .' to start the server

endlocal
