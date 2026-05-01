# WinGopher Installation Script
# This script downloads and runs the latest version of WinGopher from GitHub.

$repo = "fereshetyan/wingopher"
$url = "https://api.github.com/repos/$repo/releases/latest"

# Check for Administrator privileges
if (-NOT ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")) {
    Write-Host "Requesting Administrator privileges..." -ForegroundColor Cyan
    Start-Process powershell.exe "-NoProfile -ExecutionPolicy Bypass -Command `"iwr -useb https://raw.githubusercontent.com/$repo/main/install.ps1 | iex`"" -Verb RunAs
    exit
}

Write-Host "Checking for latest WinGopher release..." -ForegroundColor Cyan

try {
    $release = Invoke-RestMethod -Uri $url
    $asset = $release.assets | Where-Object { $_.name -like "*wingopher*.exe" } | Select-Object -First 1
    
    if (-not $asset) {
        # Fallback to search for any exe if naming convention is different
        $asset = $release.assets | Where-Object { $_.name -like "*.exe" } | Select-Object -First 1
    }

    if (-not $asset) {
        Write-Error "Could not find a .exe file in the latest release ($($release.tag_name))."
        exit
    }
    
    $downloadUrl = $asset.browser_download_url
    $destDir = Join-Path $env:LOCALAPPDATA "WinGopher"
    if (-not (Test-Path $destDir)) {
        New-Item -Path $destDir -ItemType Directory | Out-Null
    }
    
    $exePath = Join-Path $destDir "wingopher.exe"
    
    Write-Host "Downloading WinGopher $($release.tag_name)..." -ForegroundColor Cyan
    Invoke-WebRequest -Uri $downloadUrl -OutFile $exePath
    
    Write-Host "Launching WinGopher..." -ForegroundColor Green
    Start-Process $exePath
} catch {
    Write-Error "An error occurred: $_"
}
