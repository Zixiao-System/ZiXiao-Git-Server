# Xcode ZiXiao Git Server 开发指南

## 打开项目

### 方法一：使用 CMake 生成 Xcode 项目

1. 安装 CMake (如果未安装):
```bash
brew install cmake
```

2. 生成 Xcode 项目:
```bash
cd git-core
mkdir build-xcode
cd build-xcode
cmake -G Xcode ..
```

3. 打开生成的项目:
```bash
open ZiXiaoGitCore.xcodeproj
```

### 方法二：手动创建 Xcode 项目

1. 打开 Xcode
2. File → New → Project
3. 选择 macOS → Library
4. 填写项目信息:
   - Product Name: GitCore
   - Language: C++
   - Type: Dynamic Library

5. 添加源文件:
   - 将 `git-core/src/*.cpp` 添加到项目
   - 将 `git-core/include/*.h` 添加到项目

6. 配置 Build Settings:
   - **Header Search Paths**:
     - `$(PROJECT_DIR)/include`
     - `/opt/homebrew/opt/openssl/include`

   - **Library Search Paths**:
     - `/opt/homebrew/opt/openssl/lib`

   - **Other Linker Flags**:
     - `-lssl`
     - `-lcrypto`
     - `-lz`

   - **C++ Language Dialect**: `C++17`
   - **Enable Modules**: `NO`

## 构建

1. 选择 Scheme: `GitCore`
2. Product → Build (⌘B)
3. 编译后的库在: `DerivedData/.../Products/Debug/libgitcore.dylib`

## 调试 C++ 代码

1. 设置断点在 C++ 代码中
2. 从 Go 程序调用 C++ 函数时会触发断点
3. 使用 LLDB 调试器检查变量

## Go 开发

对于 Go 代码开发，推荐使用:
- **Goland**: JetBrains 的专业 Go IDE
- **VS Code**: 配合 Go 扩展

在 Xcode 中只需要维护 C++ 部分的代码。

## 常见问题

### OpenSSL 未找到
```bash
brew install openssl
```
然后在 Build Settings 中更新 Header Search Paths 和 Library Search Paths。

### C++17 特性不可用
在 Build Settings → C++ Language Dialect 中选择 `C++17` 或 `GNU++17`。

### 链接错误
确保在 Other Linker Flags 中添加了:
- `-lssl`
- `-lcrypto`
- `-lz`
