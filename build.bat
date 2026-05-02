@echo off
echo Building WinGopher with UPX compression...
wails build -platform windows/amd64 -upx -upxflags="--best --lzma"
if %errorlevel% equ 0 (
    echo.
    echo Build complete! Binary: build\bin\wingopher.exe
) else (
    echo.
    echo Build failed!
    exit /b %errorlevel%
)
