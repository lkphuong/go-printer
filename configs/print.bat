@echo off
setlocal

REM === Check arguments ===
if "%~1"=="" (
    echo Usage: print.bat "PrinterName" "FilePath"
    exit /b 1
)

if "%~2"=="" (
    echo Usage: print.bat "PrinterName" "FilePath"
    exit /b 1
)

set PRINTER=%~1
set FILE=%~2

REM === Check file exists ===
if not exist "%FILE%" (
    echo File not found: %FILE%
    exit /b 1
)

REM === Get file extension in lowercase ===
for %%I in ("%FILE%") do set EXT=%%~xI
set EXT=%EXT:~1%
set EXT=%EXT:~0,3%
set EXT=%EXT:~0,1%%EXT:~1%%EXT:~2%
set EXT=%EXT:~0,3%
set EXT=%EXT:~0,3%
set EXT=%EXT:~0,3%
set EXT=%EXT:~0,3%
set EXT=%EXT:~0,3%
set EXT=%EXT:~0,3%
set EXT=%EXT:~0,3%
set EXT=%EXT:~0,3%

REM --- Decide how to print ---
if /I "%EXT%"=="PDF" (
    REM Use Edge headless to print PDF
    powershell -NoProfile -Command "Start-Process 'msedge' -ArgumentList '--headless','--disable-gpu','--print-to-printer=""%PRINTER%""','""%FILE%""' -Wait"
) else (
    REM Use PowerShell .NET PrintDocument for images
    powershell -NoProfile -Command ^
        "Add-Type -AssemblyName System.Drawing; " ^
        "$printer='%PRINTER%'; " ^
        "$file='%FILE%'; " ^
        "if (-not (Test-Path $file)) { Write-Host 'File not found'; exit 1 }; " ^
        "$doc = New-Object System.Drawing.Printing.PrintDocument; " ^
        "$doc.PrinterSettings.PrinterName = $printer; " ^
        "if (-not $doc.PrinterSettings.IsValid) { Write-Host 'Invalid printer name'; exit 1 }; " ^
        "$img = [System.Drawing.Image]::FromFile($file); " ^
        "$doc.add_PrintPage( { param($sender,$e) $e.Graphics.DrawImage($img,0,0); $e.HasMorePages=$false } ); " ^
        "$doc.Print(); " ^
        "$img.Dispose();"
)

endlocal
