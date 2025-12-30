# Antigravity Usage Checker - Install Script
# Run: iwr https://raw.githubusercontent.com/TungCorn/antigravity-usage-checker/main/install.ps1 -useb | iex

$ErrorActionPreference = "Stop"

Write-Host "Installing Antigravity Usage Checker..." -ForegroundColor Cyan

# Define paths
$installDir = "$env:LOCALAPPDATA\agusage"
$exePath = "$installDir\agusage.exe"
$repo = "TungCorn/antigravity-usage-checker"

# Create install directory
if (!(Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir -Force | Out-Null
}

# Check if running locally with exe
$localExe = Join-Path $PSScriptRoot "agusage.exe"

if (Test-Path $localExe) {
    # Local install mode
    Write-Host "Found local usage.exe, installing..." -ForegroundColor Cyan
    Copy-Item $localExe $exePath -Force
} else {
    # Remote install mode (download from GitHub)
    Write-Host "Fetching latest release..." -ForegroundColor Cyan
    
    try {
        $latestUrl = "https://github.com/$repo/releases/latest/download/antigravity-usage-checker-windows.zip"
        $zipPath = "$env:TEMP\agusage.zip"
        
        Write-Host "Downloading..." -ForegroundColor Cyan
        Invoke-WebRequest -Uri $latestUrl -OutFile $zipPath
        
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
