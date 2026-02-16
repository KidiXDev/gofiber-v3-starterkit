@echo off
setlocal EnableDelayedExpansion

cls
echo ========================================================
echo   GoFiber V3 Starter Pack Wizard
echo ========================================================
echo.


:: Check if module name is provided as argument
if "%~1"=="" (
    echo Please enter your new module name ^(e.g., github.com/username/project^):
    set /p NEW_MODULE="> "
) else (
    set "NEW_MODULE=%~1"
)

if "%NEW_MODULE%"=="" (
    echo Error: Module name cannot be empty.
    goto :End
)

if not exist go.mod (
    echo Error: go.mod not found.
    goto :End
)

for /f "tokens=2" %%i in ('findstr /B "module" go.mod') do set OLD_MODULE=%%i

if "%OLD_MODULE%"=="" (
    echo Error: Could not determine module name from go.mod.
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

powershell -Command "$old = '%OLD_MODULE%'; $new = '%NEW_MODULE%'; $utf8 = New-Object System.Text.UTF8Encoding $False; Get-ChildItem -Recurse -Include *.go,*.mod,*.md,*.yaml,*.yml,*.json -Exclude .git | ForEach-Object { $c = [System.IO.File]::ReadAllText($_.FullName, $utf8); if ($c -match $old) { $c = $c -replace $old, $new; [System.IO.File]::WriteAllText($_.FullName, $c, $utf8); Write-Host 'Updated ' $_.FullName } }"

powershell -Command "$old = '%OLD_MODULE%'; $new = '%NEW_MODULE%'; $utf8 = New-Object System.Text.UTF8Encoding $False; $c = [System.IO.File]::ReadAllText('rename-module.sh', $utf8); $c = $c -replace 'OLD_MODULE=\"' + $old + '\"', 'OLD_MODULE=\"' + $new + '\"'; [System.IO.File]::WriteAllText('rename-module.sh', $c, $utf8)"


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
