@echo off
echo ========================================
echo   BUILDING GO REDTEAM TOOLKIT (WINDOWS)
echo ========================================
echo.

echo [1/3] Cleaning previous build...
if exist "redteam.exe" del redteam.exe

echo [2/3] Downloading dependencies...
go mod tidy
go mod download

echo [3/3] Building Windows executable...
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-s -w -H windowsgui" -o redteam.exe main.go

if errorlevel 1 (
    echo ERROR: Build failed!
    pause
    exit /b 1
)

echo.
echo ========================================
echo   BUILD COMPLETE!
echo ========================================
echo.
echo Windows: redteam.exe
echo.
pause