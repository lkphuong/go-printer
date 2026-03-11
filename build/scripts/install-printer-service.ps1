# PowerShell Script to register printer-amd64.exe as Windows Service
# Usage: powershell -NoProfile -ExecutionPolicy Bypass -File .\install-printer-service.ps1

param(
    [string]$ServiceName = "PrinterAMD64",
    [string]$DisplayName = "Go Printer Service (AMD64)",
    [string]$Description = "Go Printer Windows Service - AMD64 Edition"
)

# Check admin privileges
if (-not ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    Write-Host "Error: Administrator privileges required"
    exit 1
}

# Get executable path
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$exePath = Join-Path (Split-Path -Parent $scriptDir) "printer-amd64.exe"

Write-Host "Installing Windows Service..."
Write-Host ""

# Check if executable exists
if (-not (Test-Path $exePath)) {
    Write-Host "Error: $exePath not found"
    Write-Host "Run: make build-windows"
    exit 1
}

Write-Host "Service: $ServiceName"
Write-Host "Path: $exePath"
Write-Host ""

# Remove existing service if it exists
$existingService = Get-Service -Name $ServiceName -ErrorAction SilentlyContinue

if ($existingService) {
    Write-Host "Removing existing service..."
    if ($existingService.Status -eq "Running") {
        Stop-Service -Name $ServiceName -Force -ErrorAction SilentlyContinue
        Start-Sleep -Seconds 1
    }
    $null = sc.exe delete $ServiceName
    Start-Sleep -Seconds 1
}

# Create service
Write-Host "Creating service..."

try {
    $null = New-Service `
        -Name $ServiceName `
        -DisplayName $DisplayName `
        -Description $Description `
        -BinaryPathName $exePath `
        -StartupType Automatic `
        -ErrorAction Stop
    
    # Start service
    Write-Host "Starting service..."
    Start-Service -Name $ServiceName -ErrorAction Stop
    Start-Sleep -Seconds 1
    
    Write-Host "Done."
    Write-Host ""
    Write-Host "Commands:"
    Write-Host "  Start:  Start-Service -Name $ServiceName"
    Write-Host "  Stop:   Stop-Service -Name $ServiceName"
    Write-Host "  Delete: Remove-Service -Name $ServiceName"
    Write-Host ""
    
    exit 0
} catch {
    Write-Host "Error: $_"
    exit 1
}
