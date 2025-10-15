# Check if running as Administrator
$isAdmin = ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)

Write-Host "======================================"
Write-Host "ZiXiao Git Server - Windows Build"
Write-Host "======================================"
Write-Host ""

# Check for required tools
Write-Host "Checking dependencies..."

# Check for CMake
if (!(Get-Command cmake -ErrorAction SilentlyContinue)) {
    Write-Host "ERROR: CMake not found. Please install Visual Studio with CMake tools." -ForegroundColor Red
    exit 1
}
Write-Host "√ CMake found" -ForegroundColor Green

# Check for Go
if (!(Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "ERROR: Go not found. Please install Go 1.21+" -ForegroundColor Red
    exit 1
}
$goVersion = go version
Write-Host "√ Go found: $goVersion" -ForegroundColor Green

# Create directories
Write-Host ""
Write-Host "Creating directories..."
New-Item -ItemType Directory -Force -Path "git-core\lib" | Out-Null
New-Item -ItemType Directory -Force -Path "bin" | Out-Null
New-Item -ItemType Directory -Force -Path "data\repositories" | Out-Null
New-Item -ItemType Directory -Force -Path "logs" | Out-Null
Write-Host "√ Directories created" -ForegroundColor Green

# Build C++ library
Write-Host ""
Write-Host "Building C++ Git core library..."

Push-Location git-core

if (Test-Path "build") {
    Remove-Item -Recurse -Force build
}

New-Item -ItemType Directory -Force -Path "build" | Out-Null
Push-Location build

# Detect vcpkg
$vcpkgPath = "C:\vcpkg\scripts\buildsystems\vcpkg.cmake"
$cmakeArgs = @("-G", "Visual Studio 17 2022", "-A", "x64")

if (Test-Path $vcpkgPath) {
    Write-Host "Using vcpkg toolchain" -ForegroundColor Cyan
    $cmakeArgs += "-DCMAKE_TOOLCHAIN_FILE=$vcpkgPath"
}

cmake $cmakeArgs ..

if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: CMake configuration failed" -ForegroundColor Red
    Pop-Location
    Pop-Location
    exit 1
}

cmake --build . --config Release

if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: C++ build failed" -ForegroundColor Red
    Pop-Location
    Pop-Location
    exit 1
}

# Copy DLL
if (Test-Path "Release\gitcore.dll") {
    Copy-Item "Release\gitcore.dll" "..\lib\" -Force
    Write-Host "√ C++ library built successfully" -ForegroundColor Green
} else {
    Write-Host "ERROR: gitcore.dll not found" -ForegroundColor Red
    Pop-Location
    Pop-Location
    exit 1
}

Pop-Location
Pop-Location

# Build Go server
Write-Host ""
Write-Host "Building Go server..."

$env:CGO_ENABLED = "1"
$pwd = Get-Location
$env:CGO_CFLAGS = "-I$pwd\git-core\include"
$env:CGO_LDFLAGS = "-L$pwd\git-core\lib -lgitcore"

go build -o bin\zixiao-git-server.exe .\cmd\server

if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: Go build failed" -ForegroundColor Red
    exit 1
}

Write-Host "√ Go server built successfully" -ForegroundColor Green

# Copy DLL dependencies to bin
Write-Host ""
Write-Host "Copying dependencies..."

if (Test-Path "git-core\lib\gitcore.dll") {
    Copy-Item "git-core\lib\gitcore.dll" "bin\" -Force
}

# Try to copy vcpkg DLLs if available
$vcpkgBin = "C:\vcpkg\installed\x64-windows\bin"
if (Test-Path $vcpkgBin) {
    Get-ChildItem "$vcpkgBin\libssl*.dll" -ErrorAction SilentlyContinue | Copy-Item -Destination "bin\" -Force
    Get-ChildItem "$vcpkgBin\libcrypto*.dll" -ErrorAction SilentlyContinue | Copy-Item -Destination "bin\" -Force
    Get-ChildItem "$vcpkgBin\zlib*.dll" -ErrorAction SilentlyContinue | Copy-Item -Destination "bin\" -Force
    Write-Host "√ Dependencies copied" -ForegroundColor Green
}

Write-Host ""
Write-Host "======================================"
Write-Host "Build completed successfully!"
Write-Host "======================================"
Write-Host ""
Write-Host "To run the server:"
Write-Host "  .\bin\zixiao-git-server.exe -config .\configs\server.yaml" -ForegroundColor Cyan
Write-Host ""
Write-Host "To test the server:"
Write-Host "  .\scripts\api-test.ps1" -ForegroundColor Cyan
Write-Host ""
