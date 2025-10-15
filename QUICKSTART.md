# 快速开始指南

## 第一步：安装依赖

```bash
# macOS
brew install go openssl
xcode-select --install

# Ubuntu/Debian
sudo apt-get update
sudo apt-get install golang g++ libssl-dev zlib1g-dev
```

## 第二步：设置项目

```bash
cd ZiXiao-Git-Server
./scripts/install.sh
```

## 第三步：构建项目

```bash
make build
```

## 第四步：启动服务器

```bash
make run
```

服务器将在 `http://localhost:8080` 启动

## 第五步：测试 API

在另一个终端窗口运行：

```bash
./scripts/api-test.sh
```

## 第六步：使用 Git

### 注册用户
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

保存返回的 token。

### 创建仓库
```bash
curl -X POST http://localhost:8080/api/v1/repos \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "my-project",
    "description": "My first project",
    "is_private": false
  }'
```

### 使用 Git 克隆
```bash
git clone http://localhost:8080/alice/my-project.git
cd my-project
echo "# My Project" > README.md
git add README.md
git commit -m "Initial commit"
```

### 推送代码
对于公开仓库：
```bash
git push http://alice:YOUR_TOKEN@localhost:8080/alice/my-project.git main
```

## 常用命令

```bash
make build      # 构建项目
make run        # 运行服务器
make clean      # 清理构建产物
make test       # 运行测试
./scripts/test.sh    # 验证项目完整性
./scripts/api-test.sh # API 功能测试
```

## 配置

编辑 `configs/server.yaml`:

```yaml
server:
  host: 0.0.0.0
  port: 8080

security:
  jwt_secret: YOUR_RANDOM_SECRET  # 生产环境必须修改！
```

## 故障排除

### OpenSSL 未找到 (macOS)
```bash
brew install openssl
```

### C++ 编译错误
确保已安装 Xcode Command Line Tools:
```bash
xcode-select --install
```

### 权限错误
确保脚本可执行:
```bash
chmod +x scripts/*.sh
```

## 更多信息

- 完整文档: `README.md`
- API 文档: `docs/API.md`
- 项目总结: `PROJECT_SUMMARY.md`
