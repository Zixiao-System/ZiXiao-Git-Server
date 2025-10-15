# Installation script for Windows
Write-Host "======================================"
Write-Host "ZiXiao Git Server Installation"
Write-Host "======================================"
Write-Host ""

$ErrorActionPreference = "Stop"

# Check for Go
Write-Host "Checking for Go..."
if (!(Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "ERROR: Go is not installed" -ForegroundColor Red
    Write-Host "Please install Go 1.21+ from https://golang.org/dl/" -ForegroundColor Yellow
    exit 1
}
$goVersion = go version
Write-Host "√ Go found: $goVersion" -ForegroundColor Green

# Check for CMake
Write-Host "Checking for CMake..."
if (!(Get-Command cmake -ErrorAction SilentlyContinue)) {
    Write-Host "ERROR: CMake is not installed" -ForegroundColor Red
    Write-Host "Please install Visual Studio 2022 with 'Desktop development with C++'" -ForegroundColor Yellow
    exit 1
}
$cmakeVersion = cmake --version | Select-Object -First 1
Write-Host "√ CMake found: $cmakeVersion" -ForegroundColor Green

# Check for vcpkg
Write-Host "Checking for vcpkg..."
$vcpkgPath = "C:\vcpkg"
if (!(Test-Path $vcpkgPath)) {
    Write-Host "WARNING: vcpkg not found at $vcpkgPath" -ForegroundColor Yellow
    Write-Host "Installing vcpkg..." -ForegroundColor Cyan

    Push-Location C:\
    git clone https://github.com/Microsoft/vcpkg.git
    Push-Location vcpkg
    .\bootstrap-vcpkg.bat
    .\vcpkg integrate install
    Pop-Location
    Pop-Location

    Write-Host "√ vcpkg installed" -ForegroundColor Green
} else {
    Write-Host "√ vcpkg found" -ForegroundColor Green
}

# Install OpenSSL and zlib via vcpkg
Write-Host ""
Write-Host "Installing dependencies via vcpkg..."
Push-Location $vcpkgPath

$packages = @("openssl:x64-windows", "zlib:x64-windows")
foreach ($package in $packages) {
    Write-Host "Installing $package..." -ForegroundColor Cyan
    .\vcpkg install $package
}

Pop-Location
Write-Host "√ Dependencies installed" -ForegroundColor Green

# Install Go dependencies
Write-Host ""
Write-Host "Installing Go dependencies..."
go mod download
Write-Host "√ Go dependencies installed" -ForegroundColor Green

# Create directories
Write-Host ""
Write-Host "Creating directories..."
New-Item -ItemType Directory -Force -Path "data\repositories" | Out-Null
New-Item -ItemType Directory -Force -Path "logs" | Out-Null
New-Item -ItemType Directory -Force -Path "bin" | Out-Null
New-Item -ItemType Directory -Force -Path "git-core\lib" | Out-Null
Write-Host "√ Directories created" -ForegroundColor Green

# Create default config if it doesn't exist
if (!(Test-Path "configs\server.yaml")) {
    Write-Host ""
    Write-Host "Creating default configuration..."

    # Generate random JWT secret
    $bytes = New-Object byte[] 32
    $rng = [System.Security.Cryptography.RandomNumberGenerator]::Create()
    $rng.GetBytes($bytes)
    $jwtSecret = [System.BitConverter]::ToString($bytes).Replace("-", "").ToLower()

    @"
server:
  host: 0.0.0.0
  port: 8080
  mode: release

database:
  type: sqlite
  path: ./data/gitserver.db

git:
  repo_path: ./data/repositories
  max_repo_size: 1024
  max_file_size: 100

security:
  jwt_secret: $jwtSecret
  jwt_expiration: 24
  password_min: 8
  enable_ssh: false
  ssh_port: 2222
"@ | Out-File -FilePath "configs\server.yaml" -Encoding UTF8

    Write-Host "√ Default configuration created: configs\server.yaml" -ForegroundColor Green
    Write-Host "  IMPORTANT: Update the jwt_secret in production!" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "======================================"
Write-Host "Installation completed successfully!"
Write-Host "======================================"
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Cyan
Write-Host "  1. Review and update configs\server.yaml"
Write-Host "  2. Build the project: .\scripts\build.ps1"
Write-Host "  3. Run the server: .\bin\zixiao-git-server.exe -config .\configs\server.yaml"
Write-Host ""
