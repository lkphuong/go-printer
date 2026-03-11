@echo off
setlocal enabledelayedexpansion

set SERVICE_NAME=PrinterAMD64
set DISPLAY_NAME=Go Printer Service (AMD64)
set DESCRIPTION=Go Printer Windows Service - AMD64 Edition

cd /d "%~dp0"
cd ..
set EXE_PATH=%CD%\build\printer-amd64.exe

echo Installing Windows Service...
echo.

net session >nul 2>&1
if %errorLevel% neq 0 (
    echo Error: Administrator privileges required
    pause
    exit /b 1
)

if not exist "%EXE_PATH%" (
    echo Error: %EXE_PATH% not found
    echo Run: make build-windows
    pause
    exit /b 1
)

echo Service: %SERVICE_NAME%
echo Path: %EXE_PATH%
echo.

sc query %SERVICE_NAME% >nul 2>&1
if %errorLevel% equ 0 (
    echo Removing existing service...
    net stop %SERVICE_NAME% >nul 2>&1
    timeout /t 1 /nobreak > nul
    sc delete %SERVICE_NAME% >nul 2>&1
    timeout /t 1 /nobreak > nul
)

echo Creating service...
sc create %SERVICE_NAME% binPath= "%EXE_PATH%" DisplayName= "%DISPLAY_NAME%" >nul 2>&1

if %errorLevel% neq 0 (
    echo Error: Failed to create service
    pause
    exit /b 1
)

sc config %SERVICE_NAME% start= auto >nul 2>&1
sc description %SERVICE_NAME% "%DESCRIPTION%" >nul 2>&1

echo Starting service...
net start %SERVICE_NAME% >nul 2>&1
timeout /t 1 /nobreak > nul

echo Done.
echo.
echo Commands:
echo   Start:  net start %SERVICE_NAME%
echo   Stop:   net stop %SERVICE_NAME%
echo   Delete: sc delete %SERVICE_NAME%
echo.
pause
exit /b 0
