# Visual Studio Code 开发指南

## 安装扩展

打开 VS Code，会自动提示安装推荐扩展，或手动安装：

### 必需扩展
- **Go** (golang.go) - Go 语言支持
- **C/C++** (ms-vscode.cpptools) - C++ IntelliSense
- **clangd** (llvm-vs-code-extensions.vscode-clangd) - C++ 语言服务器

### 可选扩展
- **GitLens** - Git 增强
- **Code Spell Checker** - 拼写检查
- **EditorConfig** - 编辑器配置

## 配置说明

项目已包含以下配置文件：
- `.vscode/settings.json` - 工作区设置
- `.vscode/tasks.json` - 构建任务
- `.vscode/launch.json` - 调试配置
- `.vscode/extensions.json` - 推荐扩展

## 构建项目

### 方法一：使用任务 (推荐)

1. 按 `Cmd+Shift+B` (Mac) 或 `Ctrl+Shift+B` (Linux)
2. 选择任务:
   - **Build All** - 构建完整项目 (默认)
   - **Build C++ Only** - 只构建 C++ 库
   - **Build Go Only** - 只构建 Go 程序

### 方法二：使用终端

打开集成终端 (`` Ctrl+` ``):
```bash
make build
```

## 运行项目

### 方法一：使用任务

1. 按 `Cmd+Shift+P` (Mac) 或 `Ctrl+Shift+P` (Linux)
2. 输入 `Tasks: Run Task`
3. 选择 `Run Server`

### 方法二：使用调试器 (推荐)

1. 打开 `cmd/server/main.go`
2. 按 `F5` 开始调试
3. 服务器会启动，断点会生效

### 方法三：使用终端

```bash
make run
```

## 调试

### 调试 Go 代码

1. 在 Go 代码中设置断点 (点击行号左侧)
2. 按 `F5` 启动调试
3. 使用调试控制台查看变量

### 调试 C++ 代码

C++ 代码通过 CGo 调用，调试较为复杂：

1. 在 C++ 代码中添加日志:
```cpp
#include <iostream>
std::cout << "Debug info" << std::endl;
```

2. 或使用 LLDB 附加到进程:
```bash
lldb ./bin/zixiao-git-server
```

## 测试

### 运行项目测试

任务面板选择:
- **Run Tests** - 运行项目验证
- **Run API Tests** - 运行 API 测试

或在终端:
```bash
./scripts/test.sh
./scripts/api-test.sh
```

## Go 语言特性

### 自动格式化

保存文件时自动格式化 (已配置)

### 自动导入

保存时自动添加/删除导入

### 代码补全

输入时自动提示函数、类型等

### 跳转定义

- `Cmd+Click` 或 `F12` - 跳转到定义
- `Cmd+Shift+O` - 查看文件符号
- `Cmd+T` - 查看工作区符号

## C++ 语言特性

### IntelliSense

已配置 Include 路径，支持:
- 代码补全
- 函数签名提示
- 错误提示

### 格式化

保存时自动格式化 (使用 clang-format)

## 常见问题

### Go 扩展无法工作

1. 确保安装了 Go 1.21+
2. 运行命令: `Go: Install/Update Tools`
3. 选择所有工具安装

### C++ IntelliSense 错误

1. 检查 `settings.json` 中的 `includePath`
2. 确保 OpenSSL 已安装
3. 重启 VS Code

### CGo 编译错误

检查环境变量:
```bash
export CGO_ENABLED=1
export CGO_CFLAGS="-I./git-core/include"
export CGO_LDFLAGS="-L./git-core/lib -lgitcore"
```

## 快捷键

- `Cmd+Shift+B` - 构建
- `F5` - 开始调试
- `Shift+F5` - 停止调试
- `F9` - 切换断点
- `F10` - 单步跳过
- `F11` - 单步进入
- `` Ctrl+` `` - 切换终端
- `Cmd+P` - 快速打开文件
- `Cmd+Shift+P` - 命令面板

## 更多资源

- [VS Code Go 文档](https://code.visualstudio.com/docs/languages/go)
- [VS Code C++ 文档](https://code.visualstudio.com/docs/languages/cpp)
