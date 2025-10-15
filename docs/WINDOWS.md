# Windows Development Guide

## Prerequisites

### Option 1: Native Windows (Recommended)

#### Install Tools
1. **Visual Studio 2022** (Community Edition or higher)
   - Download from: https://visualstudio.microsoft.com/
   - Select "Desktop development with C++" workload
   - Include: MSVC v143, Windows SDK, CMake tools

2. **Go 1.21+**
   - Download from: https://golang.org/dl/
   - Add to PATH during installation

3. **Git for Windows**
   - Download from: https://git-scm.com/download/win

4. **vcpkg** (for OpenSSL and zlib)
```powershell
# Install vcpkg
cd C:\
git clone https://github.com/Microsoft/vcpkg.git
cd vcpkg
.\bootstrap-vcpkg.bat
.\vcpkg integrate install

# Install dependencies
.\vcpkg install openssl:x64-windows
.\vcpkg install zlib:x64-windows
```

### Option 2: WSL2 (Windows Subsystem for Linux)

```powershell
# Enable WSL2
wsl --install

# Install Ubuntu
wsl --install -d Ubuntu

# Inside WSL2
sudo apt-get update
sudo apt-get install golang g++ libssl-dev zlib1g-dev make
```

## Building on Windows (Native)

### Using Visual Studio 2022

1. **Open Project**
   - File → Open → CMake
   - Select `git-core/CMakeLists.txt`

2. **Configure CMake**
   - CMake Settings → CMake toolchain file:
     ```
     C:/vcpkg/scripts/buildsystems/vcpkg.cmake
     ```

3. **Build**
   - Build → Build All
   - Or press `Ctrl+Shift+B`

### Using Command Line

1. **Build C++ Library**
```powershell
cd git-core
mkdir build
cd build

cmake -G "Visual Studio 17 2022" -A x64 ^
  -DCMAKE_TOOLCHAIN_FILE=C:/vcpkg/scripts/buildsystems/vcpkg.cmake ^
  ..

cmake --build . --config Release
```

2. **Build Go Server**
```powershell
cd ../..
$env:CGO_ENABLED=1
$env:CGO_CFLAGS="-I./git-core/include"
$env:CGO_LDFLAGS="-L./git-core/build/Release -lgitcore"

go build -o bin/zixiao-git-server.exe ./cmd/server
```

### Using Makefile (with MinGW)

1. **Install MinGW-w64**
   - Download from: https://www.mingw-w64.org/
   - Or use Chocolatey: `choco install mingw`

2. **Build**
```bash
make build
```

## Running on Windows

### Native
```powershell
.\bin\zixiao-git-server.exe -config .\configs\server.yaml
```

### WSL2
```bash
./bin/zixiao-git-server -config ./configs/server.yaml
```

## PowerShell Build Script

Create `scripts/build.ps1`:
```powershell
# Check if running as Administrator
$isAdmin = ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)

Write-Host "======================================"
Write-Host "ZiXiao Git Server - Windows Build"
Write-Host "======================================"
Write-Host ""

# Create directories
New-Item -ItemType Directory -Force -Path "git-core\lib" | Out-Null
New-Item -ItemType Directory -Force -Path "bin" | Out-Null
New-Item -ItemType Directory -Force -Path "data\repositories" | Out-Null
New-Item -ItemType Directory -Force -Path "logs" | Out-Null

# Build C++ library
Write-Host "Building C++ Git core library..."
cd git-core

if (Test-Path "build") {
    Remove-Item -Recurse -Force build
}

New-Item -ItemType Directory -Force -Path "build" | Out-Null
cd build

cmake -G "Visual Studio 17 2022" -A x64 `
  -DCMAKE_TOOLCHAIN_FILE="C:/vcpkg/scripts/buildsystems/vcpkg.cmake" `
  ..

cmake --build . --config Release

# Copy DLL
Copy-Item "Release\gitcore.dll" "..\lib\" -Force
cd ..\..

Write-Host "C++ library built successfully"

# Build Go server
Write-Host ""
Write-Host "Building Go server..."

$env:CGO_ENABLED=1
$env:CGO_CFLAGS="-I$PWD\git-core\include"
$env:CGO_LDFLAGS="-L$PWD\git-core\lib -lgitcore"

go build -o bin\zixiao-git-server.exe .\cmd\server

Write-Host ""
Write-Host "======================================"
Write-Host "Build completed successfully!"
Write-Host "======================================"
Write-Host ""
Write-Host "To run the server:"
Write-Host "  .\bin\zixiao-git-server.exe -config .\configs\server.yaml"
Write-Host ""
```

## Windows Service Installation

### Create Service
```powershell
# Using NSSM (Non-Sucking Service Manager)
choco install nssm

# Install service
nssm install ZiXiaoGitServer "C:\path\to\bin\zixiao-git-server.exe"
nssm set ZiXiaoGitServer AppParameters "-config C:\path\to\configs\server.yaml"
nssm set ZiXiaoGitServer AppDirectory "C:\path\to\ZiXiao-Git-Server"

# Start service
nssm start ZiXiaoGitServer
```

### Using sc.exe
```powershell
sc.exe create ZiXiaoGitServer binPath= "C:\path\to\bin\zixiao-git-server.exe -config C:\path\to\configs\server.yaml"
sc.exe start ZiXiaoGitServer
```

## Troubleshooting

### OpenSSL Not Found
```powershell
# Reinstall via vcpkg
cd C:\vcpkg
.\vcpkg install openssl:x64-windows --force
```

### CGo Compilation Errors
```powershell
# Set environment variables
$env:CGO_ENABLED=1
$env:CC="gcc"  # Or "cl.exe" for MSVC
```

### DLL Not Found
```powershell
# Copy dependencies to bin/
copy git-core\lib\gitcore.dll bin\
copy C:\vcpkg\installed\x64-windows\bin\libssl-*.dll bin\
copy C:\vcpkg\installed\x64-windows\bin\libcrypto-*.dll bin\
copy C:\vcpkg\installed\x64-windows\bin\zlib1.dll bin\
```

## Performance Tips

1. **Use Release Build**: Significant performance improvement
2. **Disable Windows Defender** for project directory (optional)
3. **Use SSD** for repository storage
4. **Increase file handle limits** if needed

## Firewall Configuration

```powershell
# Allow inbound connections
New-NetFirewallRule -DisplayName "ZiXiao Git Server" `
  -Direction Inbound `
  -LocalPort 8080 `
  -Protocol TCP `
  -Action Allow
```

## Common Issues

### Git Bash vs PowerShell
- Use PowerShell for native Windows development
- Git Bash may have path issues with vcpkg

### Path Separators
- Use `\` in PowerShell
- Use `/` in WSL2 or Git Bash

### Line Endings
- Configure Git: `git config --global core.autocrlf true`
