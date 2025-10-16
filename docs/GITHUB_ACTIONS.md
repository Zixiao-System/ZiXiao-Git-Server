# GitHub Actions CI/CD 文档

本文档说明 ZiXiao Git Server 的 GitHub Actions 工作流配置和使用方法。

## 目录

- [工作流概览](#工作流概览)
- [后端 CI](#后端-ci)
- [前端 CI](#前端-ci)
- [发布流程](#发布流程)
- [前端部署](#前端部署)
- [本地测试](#本地测试)
- [配置 Secrets](#配置-secrets)

## 工作流概览

项目包含以下 GitHub Actions 工作流：

1. **Backend CI** (`.github/workflows/ci-backend.yml`)
   - 后端代码检查、测试和构建
   - 支持 Linux 和 macOS 多平台构建
   - 自动运行代码覆盖率测试

2. **Frontend CI** (`.github/workflows/ci-frontend.yml`)
   - 前端代码检查和构建
   - ESLint 代码质量检查
   - 构建产物分析

3. **Release** (`.github/workflows/release.yml`)
   - 自动化发布流程
   - 多平台二进制文件构建
   - Docker 镜像构建和推送
   - GitHub Release 创建

4. **Deploy Frontend** (`.github/workflows/deploy-frontend.yml`)
   - 前端自动部署到 GitHub Pages
   - 支持手动触发部署

## 后端 CI

### 触发条件

- 推送到 `main` 或 `develop` 分支
- 针对 `main` 或 `develop` 的 Pull Request
- 修改以下路径时触发：
  - `cmd/**`
  - `internal/**`
  - `pkg/**`
  - `git-core/**`
  - `go.mod` / `go.sum`
  - `Makefile`

### 工作流步骤

#### 1. Go 代码检查 (lint-go)

```yaml
- go fmt 检查代码格式
- go vet 静态分析
- golint 代码风格检查
- staticcheck 高级静态分析
```

#### 2. Go 测试 (test-go)

```yaml
- 安装系统依赖 (OpenSSL, zlib)
- 构建 C++ 库
- 运行单元测试和集成测试
- 生成代码覆盖率报告
- 上传到 Codecov
```

#### 3. C++ 库构建 (build-cpp)

```yaml
- 多平台构建 (Ubuntu, macOS)
- 验证库文件生成
```

#### 4. Go 服务器构建 (build-backend)

```yaml
- 完整构建后端服务器
- 测试服务器启动
- 上传构建产物
```

### 本地运行测试

模拟 CI 环境在本地运行测试：

```bash
# Go 代码格式检查
gofmt -l .

# Go 代码静态分析
go vet ./...

# 安装 golint 和 staticcheck
go install golang.org/x/lint/golint@latest
go install honnef.co/go/tools/cmd/staticcheck@latest

# 运行 golint
golint ./...

# 运行 staticcheck
staticcheck ./...

# 运行测试并生成覆盖率
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

# 查看覆盖率报告
go tool cover -html=coverage.out
```

## 前端 CI

### 触发条件

- 推送到 `main` 或 `develop` 分支
- 针对 `main` 或 `develop` 的 Pull Request
- 修改 `frontend/**` 目录时触发

### 工作流步骤

#### 1. 代码检查 (lint)

```yaml
- ESLint 检查 JavaScript/Vue 代码
- 自动修复可修复的问题
```

#### 2. 构建 (build)

```yaml
- npm ci 安装依赖
- npm run build 生产构建
- 检查构建产物大小
- 上传构建产物
```

#### 3. 测试 (test)

```yaml
- 运行单元测试（如果存在）
```

#### 4. Bundle 分析 (analyze)

```yaml
- 分析构建产物大小
- 生成大小报告
```

### 本地运行测试

```bash
cd frontend

# 安装依赖
npm ci

# 运行 ESLint
npm run lint

# 修复 ESLint 错误
npm run lint -- --fix

# 构建生产版本
npm run build

# 分析构建产物
ls -lh ../web/dist/
du -sh ../web/dist/
```

## 发布流程

### 触发方式

#### 1. 标签推送（推荐）

```bash
# 创建版本标签
git tag -a v1.0.0 -m "Release v1.0.0"

# 推送标签到远程
git push origin v1.0.0
```

#### 2. 手动触发

在 GitHub Actions 页面手动运行 Release 工作流，输入版本号。

### 发布流程步骤

#### 1. 构建前端

```yaml
- 构建 Vue 3 前端应用
- 打包为 tar.gz 归档文件
```

#### 2. 构建后端

多平台构建：

| 平台 | 架构 | 产物名称 |
|------|------|---------|
| Linux | AMD64 | `zixiao-git-server-linux-amd64.tar.gz` |
| Linux | ARM64 | `zixiao-git-server-linux-arm64.tar.gz` |
| macOS | Intel | `zixiao-git-server-darwin-amd64.tar.gz` |
| macOS | Apple Silicon | `zixiao-git-server-darwin-arm64.tar.gz` |

每个包包含：
- 编译好的二进制文件
- 配置文件模板
- 启动脚本
- README 和 LICENSE

#### 3. 创建 GitHub Release

```yaml
- 生成版本说明
- 附加所有构建产物
- 创建 GitHub Release
```

#### 4. 构建和推送 Docker 镜像

```yaml
- 多架构构建 (AMD64, ARM64)
- 推送到 Docker Hub
- 标签策略:
  - v1.0.0 (完整版本)
  - v1.0 (主版本.次版本)
  - v1 (主版本)
  - latest (最新版本)
```

### 发布检查清单

发布新版本前，请确保：

- [ ] 所有 CI 测试通过
- [ ] 更新 `CHANGELOG.md`
- [ ] 更新 `cmd/server/main.go` 中的版本号
- [ ] 创建并推送版本标签
- [ ] 验证 Release 创建成功
- [ ] 验证 Docker 镜像推送成功
- [ ] 测试下载的二进制文件

## 前端部署

### 部署到 GitHub Pages

#### 自动部署

推送到 `main` 分支时自动触发部署。

#### 手动部署

1. 进入 Actions 页面
2. 选择 "Deploy Frontend" 工作流
3. 点击 "Run workflow"
4. 选择分支并运行

### 配置自定义域名

在 `.github/workflows/deploy-frontend.yml` 中修改 `cname` 字段：

```yaml
- name: Deploy to GitHub Pages
  uses: peaceiris/actions-gh-pages@v3
  with:
    cname: your-domain.com
```

然后在域名 DNS 设置中添加 CNAME 记录指向 GitHub Pages。

### 配置 API 地址

在工作流中设置前端 API 地址：

```yaml
- name: Build for production
  run: |
    cd frontend
    npm run build
  env:
    VITE_API_BASE_URL: https://api.your-domain.com/api/v1
```

## 配置 Secrets

### 必需的 Secrets

在 GitHub 仓库设置中添加以下 Secrets：

#### Docker Hub 认证

```
DOCKER_USERNAME: Docker Hub 用户名
DOCKER_PASSWORD: Docker Hub 密码或访问令牌
```

设置步骤：
1. 访问 https://hub.docker.com/settings/security
2. 创建访问令牌
3. 在 GitHub 仓库 Settings → Secrets and variables → Actions 中添加

#### Codecov (可选)

```
CODECOV_TOKEN: Codecov 上传令牌
```

获取步骤：
1. 在 https://codecov.io 注册并关联仓库
2. 获取上传令牌
3. 添加到 GitHub Secrets

### 可选的 Secrets

```
SLACK_WEBHOOK: Slack 通知 Webhook URL
DISCORD_WEBHOOK: Discord 通知 Webhook URL
```

## 本地测试

### 使用 act 在本地运行 GitHub Actions

安装 [act](https://github.com/nektos/act)：

```bash
# macOS
brew install act

# Linux
curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash
```

运行工作流：

```bash
# 列出所有工作流
act -l

# 运行后端 CI
act -W .github/workflows/ci-backend.yml

# 运行前端 CI
act -W .github/workflows/ci-frontend.yml

# 使用特定事件触发
act push -W .github/workflows/ci-backend.yml
```

注意：
- 某些工作流可能需要配置 Secrets
- 多平台构建在本地可能无法完全模拟
- 建议使用真实的 GitHub Actions 进行最终测试

## 故障排查

### 常见问题

#### 1. C++ 构建失败

**问题**: 找不到 OpenSSL 或 zlib 库

**解决方案**:
- 检查工作流中的依赖安装步骤
- 验证 `pkg/gitcore/gitcore.go` 中的 CGo 路径配置

#### 2. Go 测试失败

**问题**: 测试超时或失败

**解决方案**:
- 检查 `go.mod` 依赖是否最新
- 确保 C++ 库已正确构建
- 查看测试日志定位具体问题

#### 3. Docker 构建失败

**问题**: 多阶段构建失败

**解决方案**:
- 检查 Dockerfile 中的路径是否正确
- 验证所有 COPY 指令的源文件存在
- 确保各阶段依赖正确传递

#### 4. 前端构建失败

**问题**: npm install 或 build 失败

**解决方案**:
- 检查 `package-lock.json` 是否提交
- 验证 Node.js 版本兼容性
- 检查 MDUI 导入语法（必须使用命名导入）

#### 5. Release 创建失败

**问题**: GitHub Release 创建失败

**解决方案**:
- 检查标签格式是否为 `vX.Y.Z`
- 确保仓库有写入权限
- 验证 `GITHUB_TOKEN` 权限

### 查看工作流日志

1. 进入 GitHub 仓库
2. 点击 "Actions" 标签
3. 选择失败的工作流运行
4. 展开失败的步骤查看详细日志
5. 下载日志文件进行本地分析

## 工作流优化

### 加速构建

1. **使用缓存**

```yaml
- name: Setup Go
  uses: actions/setup-go@v5
  with:
    go-version: '1.21'
    cache: true  # 启用 Go 模块缓存

- name: Setup Node.js
  uses: actions/setup-node@v4
  with:
    node-version: '20'
    cache: 'npm'  # 启用 npm 缓存
```

2. **并行执行**

```yaml
strategy:
  matrix:
    os: [ubuntu-latest, macos-latest]
  max-parallel: 2  # 最多并行 2 个任务
```

3. **条件执行**

```yaml
- name: Run expensive tests
  if: github.event_name == 'push' && github.ref == 'refs/heads/main'
  run: make test-integration
```

### 减少运行次数

使用 `paths` 过滤器只在相关文件变更时运行：

```yaml
on:
  push:
    paths:
      - 'internal/**'
      - 'pkg/**'
      - '!**/*.md'  # 排除 Markdown 文件
```

## 最佳实践

1. **保持工作流简洁**
   - 每个工作流专注于单一职责
   - 避免过度复杂的条件逻辑

2. **使用版本化的 Actions**
   - 使用 `@v4` 而不是 `@main`
   - 定期更新 Actions 版本

3. **保护敏感信息**
   - 使用 Secrets 存储敏感数据
   - 不要在日志中输出 Secrets

4. **及时通知**
   - 配置失败通知（Slack/Discord）
   - 定期检查工作流运行状态

5. **文档化**
   - 为复杂步骤添加注释
   - 保持 README 和文档同步更新

## 相关资源

- [GitHub Actions 官方文档](https://docs.github.com/en/actions)
- [Docker Build Push Action](https://github.com/docker/build-push-action)
- [Codecov Action](https://github.com/codecov/codecov-action)
- [act - 本地运行 GitHub Actions](https://github.com/nektos/act)

---

**最后更新**: 2025-10-16
**维护者**: ZiXiao Team
