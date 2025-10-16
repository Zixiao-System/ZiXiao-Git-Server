# ZiXiao Git Server

[![Backend CI](https://github.com/Zixiao-System/ZiXiao-Git-Server/workflows/Backend%20CI/badge.svg)](https://github.com/Zixiao-System/ZiXiao-Git-Server/actions/workflows/ci-backend.yml)
[![Frontend CI](https://github.com/Zixiao-System/ZiXiao-Git-Server/workflows/Frontend%20CI/badge.svg)](https://github.com/Zixiao-System/ZiXiao-Git-Server/actions/workflows/ci-frontend.yml)
[![Release](https://github.com/Zixiao-System/ZiXiao-Git-Server/workflows/Release/badge.svg)](https://github.com/Zixiao-System/ZiXiao-Git-Server/actions/workflows/release.yml)
[![codecov](https://codecov.io/gh/Zixiao-System/ZiXiao-Git-Server/branch/main/graph/badge.svg)](https://codecov.io/gh/Zixiao-System/ZiXiao-Git-Server)
[![Go Report Card](https://goreportcard.com/badge/github.com/zixiao/git-server)](https://goreportcard.com/report/github.com/zixiao/git-server)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Docker](https://img.shields.io/docker/v/zixiao/git-server?label=docker)](https://hub.docker.com/r/zixiao/git-server)

一个使用 **Go** 和 **C++** 实现的轻量级 Git 服务器，类似于 GitLab，支持 HTTP Git 协议、用户认证、仓库管理等功能。

## 特性

- **混合架构**: Go 处理业务逻辑和 HTTP 服务，C++ 实现 Git 核心操作
- **Vue 3 前端**: 现代化单页应用，支持深色模式
- **用户认证**: JWT token 认证，密码加密存储
- **仓库管理**: 创建、删除、列表、权限控制
- **Git 协议**: 支持 HTTP Git 协议 (push/pull/clone)
- **权限系统**: 公开/私有仓库，协作者管理
- **REST API**: 完整的 RESTful API
- **多数据库支持**: SQLite3、PostgreSQL、SQL Server
- **Docker 支持**: 提供完整的 Docker 镜像和 Docker Compose 配置
- **CI/CD**: GitHub Actions 自动化构建和发布

## 技术栈

### 后端
- **Go 1.21+**: HTTP 服务器、业务逻辑
- **Gin**: Web 框架
- **JWT**: 用户认证
- **数据库**: SQLite3 / PostgreSQL / SQL Server

### Git 核心
- **C++ 17**: Git 对象模型、仓库操作
- **OpenSSL**: SHA-1 计算
- **zlib**: 数据压缩

## 项目结构

```
ZiXiao-Git-Server/
├── cmd/
│   └── server/          # 主程序入口
├── internal/
│   ├── api/             # HTTP API 处理器
│   ├── auth/            # 认证系统
│   ├── config/          # 配置管理
│   ├── database/        # 数据库操作
│   ├── models/          # 数据模型
│   └── repository/      # 仓库管理
├── pkg/
│   └── gitcore/         # C++ Git 核心库的 Go 接口
├── git-core/
│   ├── include/         # C++ 头文件
│   ├── src/             # C++ 源文件
│   └── lib/             # 编译后的动态库
├── configs/             # 配置文件
├── scripts/             # 构建脚本
├── data/                # 数据目录
│   └── repositories/    # Git 仓库存储
└── logs/                # 日志文件
```

## 快速开始

### 依赖要求

- Go 1.21+
- C++ 编译器 (g++ 支持 C++17)
- OpenSSL 开发库
- zlib 开发库

#### macOS
```bash
brew install go openssl
xcode-select --install
```

#### Ubuntu/Debian
```bash
sudo apt-get update
sudo apt-get install golang g++ libssl-dev zlib1g-dev
```

#### CentOS/RHEL
```bash
sudo yum install golang gcc-c++ openssl-devel zlib-devel
```

#### Windows
```powershell
# Using vcpkg
cd C:\
git clone https://github.com/Microsoft/vcpkg.git
cd vcpkg
.\bootstrap-vcpkg.bat
.\vcpkg integrate install
.\vcpkg install openssl:x64-windows zlib:x64-windows

# Install Go from https://golang.org/dl/
# Install Visual Studio 2022 with "Desktop development with C++"
```

### 安装

1. 克隆项目
```bash
git clone https://github.com/Zixiao-System/ZiXiao-Git-Server.git
cd ZiXiao-Git-Server
```

2. 运行安装脚本
```bash
./scripts/install.sh
```

3. 配置服务器
编辑 `configs/server.yaml`，修改 `jwt_secret` 等配置

4. 构建项目
```bash
make build
```

5. 运行服务器
```bash
make run
```

服务器将在 `http://localhost:8080` 启动

## Docker 部署

### 使用 Docker Compose (推荐)

1. 克隆项目
```bash
git clone https://github.com/Zixiao-System/ZiXiao-Git-Server.git
cd ZiXiao-Git-Server
```

2. 创建环境变量文件
```bash
cp .env.example .env
# 编辑 .env 文件，修改 JWT_SECRET 等配置
```

3. 启动服务
```bash
# 基础部署 (SQLite)
docker-compose up -d

# 使用 PostgreSQL
docker-compose --profile postgres up -d

# 使用 SQL Server
docker-compose --profile sqlserver up -d

# 使用 Redis 缓存
docker-compose --profile redis up -d

# 完整部署 (PostgreSQL + Redis + Adminer)
docker-compose --profile postgres --profile redis --profile adminer up -d
```

4. 查看日志
```bash
docker-compose logs -f
```

5. 停止服务
```bash
docker-compose down
```

### 使用 Docker

1. 拉取镜像
```bash
docker pull zixiao/git-server:latest
```

2. 运行容器
```bash
docker run -d \
  --name zixiao-git-server \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/logs:/app/logs \
  -e JWT_SECRET=your-secret-key \
  zixiao/git-server:latest
```

3. 使用自定义配置
```bash
docker run -d \
  --name zixiao-git-server \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/logs:/app/logs \
  -v $(pwd)/configs/server.yaml:/app/configs/server.yaml:ro \
  zixiao/git-server:latest
```

### 构建自定义镜像

```bash
# 构建镜像
docker build -t zixiao-git-server:custom .

# 多架构构建
docker buildx build --platform linux/amd64,linux/arm64 -t zixiao-git-server:custom .
```

## 使用方法

### 用户注册
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "email": "alice@example.com",
    "password": "password123",
    "full_name": "Alice Smith"
  }'
```

### 用户登录
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "password": "password123"
  }'
```

返回的 `token` 用于后续 API 调用。

### 创建仓库
```bash
curl -X POST http://localhost:8080/api/v1/repos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "my-project",
    "description": "My awesome project",
    "is_private": false
  }'
```

### Git 操作

#### 克隆仓库
```bash
git clone http://localhost:8080/alice/my-project.git
```

#### 推送代码
```bash
cd my-project
git add .
git commit -m "Initial commit"
git push origin main
```

#### 私有仓库认证
对于私有仓库，使用以下格式：
```bash
git clone http://alice:YOUR_TOKEN@localhost:8080/alice/my-project.git
```

## API 文档

### 认证 API

- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录

### 用户 API

- `GET /api/v1/user` - 获取当前用户信息 (需认证)
- `GET /api/v1/users/:username` - 获取用户信息

### 仓库 API

- `POST /api/v1/repos` - 创建仓库 (需认证)
- `GET /api/v1/repos/:owner/:repo` - 获取仓库信息
- `DELETE /api/v1/repos/:owner/:repo` - 删除仓库 (需认证)
- `GET /api/v1/users/:owner/repos` - 列出用户的仓库

### 协作者 API

- `POST /api/v1/repos/:owner/:repo/collaborators` - 添加协作者 (需认证)
- `DELETE /api/v1/repos/:owner/:repo/collaborators/:username` - 移除协作者 (需认证)

### Git HTTP 协议

- `GET /:owner/:repo/info/refs?service=git-upload-pack` - Git fetch/pull
- `POST /:owner/:repo/git-receive-pack` - Git push
- `POST /:owner/:repo/git-upload-pack` - Git fetch/pull

## 配置说明

### 数据库配置

ZiXiao Git Server 支持多种数据库后端：

- **SQLite3** (默认): 适合开发和小型部署，无需额外配置
- **PostgreSQL**: 推荐用于生产环境，支持高并发
- **SQL Server**: 适合企业环境和 Windows 部署

详细的数据库配置指南请参见：
- [数据库配置文档](docs/DATABASE.md) - 各数据库的详细配置、迁移和优化
- [Docker 部署文档](docs/DOCKER_DEPLOYMENT.md) - Docker 环境下的数据库部署

### 服务器配置

`configs/server.yaml` 配置选项：

```yaml
server:
  host: 0.0.0.0        # 监听地址
  port: 8080           # 监听端口
  mode: release        # debug 或 release

database:
  type: sqlite3        # 数据库类型: sqlite3, postgres, sqlserver
  path: ./data/gitserver.db  # SQLite 数据库文件路径

  # PostgreSQL 配置 (取消注释以使用)
  # type: postgres
  # host: localhost
  # port: 5432
  # name: zixiao_git
  # user: postgres
  # password: postgres
  # sslmode: disable   # disable, require, verify-ca, verify-full

  # SQL Server 配置 (取消注释以使用)
  # type: sqlserver
  # host: localhost
  # port: 1433
  # name: zixiao_git
  # user: sa
  # password: YourPassword123

git:
  repo_path: ./data/repositories  # 仓库存储路径
  max_repo_size: 1024  # 仓库最大大小 (MB)
  max_file_size: 100   # 文件最大大小 (MB)

security:
  jwt_secret: CHANGE_ME  # JWT 密钥 (生产环境必须修改)
  jwt_expiration: 24     # Token 有效期 (小时)
  password_min: 8        # 最小密码长度
  enable_ssh: false      # SSH 支持 (未实现)
  ssh_port: 2222
```

## IDE 支持

### Visual Studio Code

推荐安装扩展：
- **Go**: Go 语言支持
- **C/C++**: C++ IntelliSense
- **clangd**: C++ 代码补全和分析

`.vscode/settings.json` 配置：
```json
{
  "go.buildOnSave": "workspace",
  "go.lintOnSave": "workspace",
  "C_Cpp.default.includePath": [
    "${workspaceFolder}/git-core/include"
  ],
  "C_Cpp.default.compilerPath": "/usr/bin/g++",
  "files.exclude": {
    "**/*.o": true,
    "**/bin": true,
    "**/git-core/lib": true
  }
}
```

`.vscode/tasks.json` 配置：
```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Build All",
      "type": "shell",
      "command": "make build",
      "group": {
        "kind": "build",
        "isDefault": true
      }
    },
    {
      "label": "Run Server",
      "type": "shell",
      "command": "make run"
    }
  ]
}
```

### Xcode

1. 为 C++ 项目创建 Xcode 工程：
```bash
cd git-core
mkdir xcode-project && cd xcode-project
cmake -G Xcode ..
```

或手动创建项目：
- File → New → Project → macOS → Library
- 添加 `git-core/src/*.cpp` 文件
- 设置 Header Search Paths: `git-core/include`
- 设置 Library Search Paths: `/opt/homebrew/opt/openssl/lib`
- 链接库: `libssl.dylib`, `libcrypto.dylib`, `libz.dylib`

2. Go 开发推荐使用 Goland 或 VS Code

### Visual Studio (Windows)

暂不支持 Windows 原生编译，建议使用 WSL2：
```bash
# 在 WSL2 中安装依赖
sudo apt-get update
sudo apt-get install golang g++ libssl-dev zlib1g-dev

# 然后正常构建
make build
```

### CLion

1. 打开项目目录
2. CLion 会自动识别 CMakeLists.txt（如果存在）
3. 配置 C++ 标准为 C++17
4. 设置 Include 路径: `git-core/include`
5. Go 开发使用 Goland 插件

## 开发

### 构建命令

```bash
make build      # 构建项目
make build-cpp  # 只构建 C++ 库
make build-go   # 只构建 Go 程序
make clean      # 清理构建产物
make run        # 构建并运行
make test       # 运行测试
make init       # 初始化项目目录
```

### 项目架构

1. **C++ Git 核心层**
   - `GitRepository`: 仓库操作 (初始化、引用管理)
   - `GitObject`: Git 对象模型 (blob, tree, commit)
   - `GitProtocol`: Git 协议处理 (pkt-line, ref advertisement)
   - `GitPack`: Pack 文件处理

2. **Go 业务层**
   - `config`: 配置管理
   - `database`: 数据库初始化和 schema
   - `models`: 数据模型 (User, Repository, etc.)
   - `auth`: JWT 认证、密码加密
   - `repository`: 仓库 CRUD、权限检查
   - `api`: HTTP 路由和处理器

3. **CGo 桥接层**
   - `git_c_api.h/cpp`: C API 接口
   - `pkg/gitcore/gitcore.go`: Go CGo 绑定

## 性能优化

- C++ 编译优化: `-O2`
- Git 对象压缩: zlib
- 数据库索引: 用户、仓库、协作者
- 连接池: SQLite

## 安全性

- 密码使用 bcrypt 加密
- JWT token 认证
- 私有仓库访问控制
- SQL 注入防护 (参数化查询)

## 路线图

- [x] 基础 HTTP Git 协议
- [x] 用户认证和授权
- [x] 仓库 CRUD
- [x] 权限管理
- [x] Vue 3 前端 Web UI
- [x] Nginx 反向代理
- [x] Docker 和 Docker Compose 支持
- [x] GitHub Actions CI/CD
- [x] 多平台构建和发布
- [x] PostgreSQL 数据库支持
- [x] SQL Server 数据库支持
- [ ] 数据库迁移系统
- [ ] SSH 协议支持
- [ ] Webhook
- [ ] CI/CD 集成
- [ ] 代码审查
- [ ] Issue 跟踪
- [ ] Git LFS 支持

## 许可证

MIT License

## 贡献

欢迎提交 Pull Request 和 Issue！

开发指南：
- [数据库配置](docs/DATABASE.md)
- [Docker 部署](docs/DOCKER_DEPLOYMENT.md)
- [VS Code 开发指南](docs/VSCODE.md)
- [Xcode 开发指南](docs/XCODE.md)
- [Windows 开发指南](docs/WINDOWS.md)
- [API 文档](docs/API.md)

## 致谢

灵感来自 GitLab, Gitea 和 Gogs。
