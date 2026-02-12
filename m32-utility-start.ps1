# --- CHURCH MIXER SYNC INSTALLER (PODMAN VERSION) ---

# 1. Ensure we are in the script directory
$PSScriptRoot = Split-Path -Parent $MyInvocation.MyCommand.Definition
Set-Location $PSScriptRoot

Write-Host "--- Starting Mixer Sync Utility Setup (Podman) ---" -ForegroundColor Cyan

# 2. Check if Podman is running
if (!(podman machine inspect | ConvertFrom-Json).State -eq "running") {
    Write-Host "Podman machine is not running. Starting it now..." -ForegroundColor Yellow
    podman machine start
}

# 3. Create/Overwrite Dockerfile & docker-compose.yml
# (Content remains the same as the Docker version)

# 4. Create the Data folder
if (!(Test-Path "node-red-data")) {
    New-Item -ItemType Directory -Path "node-red-data"
}

# 5. Build and Run using Podman Compose
Write-Host "Building container via Podman..." -ForegroundColor Yellow
# Podman now supports 'podman compose' directly in newer versions
podman compose up -d --build

# 6. Wait for services
Write-Host "Waiting for Node-RED to initialize..." -ForegroundColor Yellow
Start-Sleep -Seconds 10

# 7. Launch the Dashboard
Write-Host "Setup Complete!" -ForegroundColor Green
Start-Process "http://localhost:1880/dashboard/mixer"