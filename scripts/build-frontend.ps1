# Frontend build script for Windows

Write-Host "======================================"
Write-Host "Building Frontend"
Write-Host "======================================"
Write-Host ""

# Check if node_modules exists
if (!(Test-Path "frontend\node_modules")) {
    Write-Host "Installing dependencies..."
    Push-Location frontend
    npm install
    Pop-Location
}

# Build frontend
Write-Host "Building Vue application..."
Push-Location frontend
npm run build
Pop-Location

Write-Host ""
Write-Host "√ Frontend build completed successfully!" -ForegroundColor Green
Write-Host "Built files are in: web\dist\"
Write-Host ""
