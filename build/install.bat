@echo off

set SERVICE_NAME=PrinterAMD64
set EXE_PATH=%~dp0printer-amd64.exe

echo Installing service...

sc stop %SERVICE_NAME% >nul 2>&1
sc delete %SERVICE_NAME% >nul 2>&1

timeout /t 2 >nul

sc create %SERVICE_NAME% binPath= "%EXE_PATH%" start= auto

sc qc %SERVICE_NAME%

echo Starting service...

sc start %SERVICE_NAME%

pause