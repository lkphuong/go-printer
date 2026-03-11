@echo off

set SERVICE_NAME=PrinterAMD64
set EXE_PATH=%CD%\printer-amd64.exe

echo Checking service...

sc query %SERVICE_NAME% >nul 2>&1
if %errorlevel%==0 (
    echo Service exists. Stopping...

    sc stop %SERVICE_NAME% >nul 2>&1
    timeout /t 2 /nobreak >nul

    echo Deleting old service...
    sc delete %SERVICE_NAME%
    timeout /t 2 /nobreak >nul
)

echo Creating new service...

sc create %SERVICE_NAME% binPath= "%EXE_PATH%" start= auto

echo Starting service...

sc start %SERVICE_NAME%

echo Done.
pause