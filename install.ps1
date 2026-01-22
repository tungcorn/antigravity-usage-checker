# ============================================================================
# Antigravity Usage Checker - Install Script for Windows
# ============================================================================
#
# WHAT THIS SCRIPT DOES:
# 1. Downloads agusage.exe from GitHub Releases (or uses local copy if present)
# 2. Installs to: %LOCALAPPDATA%\agusage\agusage.exe
# 3. Adds the install directory to your PATH (user-level, no admin required)
#
# WHAT THIS SCRIPT DOES NOT DO:
# - Does NOT require administrator privileges
# - Does NOT modify system files
# - Does NOT install any dependencies
# - Does NOT send any data to external servers
#
# You can review this script before running. Source code:
# https://github.com/tungcorn/antigravity-usage-checker
#
# Run: iwr https://raw.githubusercontent.com/tungcorn/antigravity-usage-checker/main/install.ps1 -OutFile $env:TEMP\install.ps1; . $env:TEMP\install.ps1
# Install specific version: ... $env:TEMP\install.ps1 -Version 0.5.0
# ============================================================================

param(
    [string]$Version = "latest"  # Default: download latest release
)

# Stop on any error
$ErrorActionPreference = "Stop"

if ($Version -eq "latest") {
    Write-Host "Installing Antigravity Usage Checker (latest)..." -ForegroundColor Cyan
} else {
    Write-Host "Installing Antigravity Usage Checker v$Version..." -ForegroundColor Cyan
}

# Define paths
$installDir = "$env:LOCALAPPDATA\agusage"
$exePath = "$installDir\agusage.exe"
$repo = "tungcorn/antigravity-usage-checker"

# Create install directory
if (!(Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir -Force | Out-Null
}

# Check if running locally with exe
$localExe = Join-Path $PSScriptRoot "agusage.exe"

if (Test-Path $localExe) {
    # Local install mode
    Write-Host "Found local agusage.exe, installing..." -ForegroundColor Green
    
    # Check if updating
    if (Test-Path $exePath) {
        $oldVersion = (& $exePath --version 2>$null) -replace '.*v([0-9.]+).*', '$1'
        Write-Host "Updating from version $oldVersion..." -ForegroundColor Yellow
    }
    
    Copy-Item $localExe $exePath -Force
    
    # Get new version
    $newVersion = (& $exePath --version 2>$null) -replace '.*v([0-9.]+).*', '$1'
    Write-Host "Installed version $newVersion" -ForegroundColor Green
} else {
    # Remote install mode (download from GitHub)
    Write-Host "Fetching release..." -ForegroundColor Cyan
    
    try {
        if ($Version -eq "latest") {
            $downloadUrl = "https://github.com/$repo/releases/latest/download/antigravity-usage-checker-windows-amd64.zip"
        } else {
            $downloadUrl = "https://github.com/$repo/releases/download/v$Version/antigravity-usage-checker-windows-amd64.zip"
        }
        $zipPath = "$env:TEMP\agusage.zip"
        
        Write-Host "Downloading..." -ForegroundColor Cyan
        Invoke-WebRequest -Uri $downloadUrl -OutFile $zipPath
        
        Write-Host "Extracting..." -ForegroundColor Cyan
        Expand-Archive -Path $zipPath -DestinationPath $env:TEMP -Force
        
        # Move exe to install dir
        $extractedExe = "$env:TEMP\agusage.exe"
        Move-Item $extractedExe $exePath -Force
        
        # Cleanup
        Remove-Item $zipPath -Force
        if (Test-Path "$env:TEMP\install.ps1") { Remove-Item "$env:TEMP\install.ps1" -Force }
        
    } catch {
        Write-Host "Error downloading/installing: $_" -ForegroundColor Red
        exit 1
    }
}

Write-Host "Installed to: $installDir" -ForegroundColor Green

# Create alias 'agu' for convenience
$aguPath = "$installDir\agu.exe"
Copy-Item $exePath $aguPath -Force
Write-Host "Created alias 'agu' for quick access" -ForegroundColor Green

# Add to PATH
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -notlike "*$installDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$userPath;$installDir", "User")
    Write-Host "Added to PATH" -ForegroundColor Green
    $env:Path += ";$installDir" # Update current session for immediate use
} else {
    Write-Host "Already in PATH" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Installation complete! Run 'agusage' to start." -ForegroundColor Cyan
