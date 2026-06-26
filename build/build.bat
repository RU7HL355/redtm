@echo off
echo ========================================
echo   BUILDING GO REDTEAM TOOLKIT
echo ========================================
echo.

echo [1/4] Cleaning previous build...
if exist "redteam.exe" del redteam.exe
if exist "redteam" rmdir /s /q redteam

echo [2/4] Downloading dependencies...
go mod tidy
go mod download

echo [3/4] Building Windows executable...
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-s -w -H windowsgui" -o redteam.exe main.go

if errorlevel 1 (
    echo ERROR: Build failed!
    pause
    exit /b 1
)

echo [4/4] Building Linux executable...
set GOOS=linux
set GOARCH=amd64
go build -ldflags="-s -w" -o redteam_linux main.go

echo.
echo ========================================
echo   BUILD COMPLETE!
echo ========================================
echo.
echo Windows: redteam.exe
echo Linux:   redteam_linux
echo.
pause