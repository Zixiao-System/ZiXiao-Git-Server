<template>
  <div>
    <mdui-top-app-bar>
      <mdui-top-app-bar-title>ZiXiao Git Server</mdui-top-app-bar-title>
      <div style="flex-grow: 1" />
      <mdui-button-icon
        :icon="themeIcon"
        @click="toggleTheme"
      />
      <mdui-button-icon
        icon="open_in_new"
        href="https://github.com/Zixiao-System/ZiXiao-Git-Server"
        target="_blank"
      />
    </mdui-top-app-bar>

    <div class="hero-section">
      <div class="hero-title">
        🚀 ZiXiao Git Server
      </div>
      <div class="hero-subtitle">
        轻量级、高性能的 Git 服务器
      </div>
    </div>

    <div class="content-section">
      <div class="stats-grid">
        <mdui-card
          class="stat-card"
          variant="elevated"
        >
          <mdui-icon
            name="code"
            class="feature-icon"
          />
          <div class="stat-value">
            Go + C++
          </div>
          <div class="stat-label">
            混合架构
          </div>
        </mdui-card>

        <mdui-card
          class="stat-card"
          variant="elevated"
        >
          <mdui-icon
            name="api"
            class="feature-icon"
          />
          <div class="stat-value">
            REST API
          </div>
          <div class="stat-label">
            完整接口
          </div>
        </mdui-card>

        <mdui-card
          class="stat-card"
          variant="elevated"
        >
          <mdui-icon
            name="security"
            class="feature-icon"
          />
          <div class="stat-value">
            JWT
          </div>
          <div class="stat-label">
            安全认证
          </div>
        </mdui-card>
      </div>

      <h2 style="text-align: center; margin-bottom: 32px;">
        核心特性
      </h2>
      <div class="features-grid">
        <mdui-card
          class="feature-card"
          variant="outlined"
        >
          <mdui-icon
            name="lock"
            class="feature-icon"
          />
          <div class="feature-title">
            🔐 用户认证
          </div>
          <div class="feature-description">
            JWT token 认证，bcrypt 密码加密，确保用户数据安全
          </div>
        </mdui-card>

        <mdui-card
          class="feature-card"
          variant="outlined"
        >
          <mdui-icon
            name="folder"
            class="feature-icon"
          />
          <div class="feature-title">
            📦 仓库管理
          </div>
          <div class="feature-description">
            完整的仓库 CRUD 操作，支持公开/私有仓库，权限细粒度控制
          </div>
        </mdui-card>

        <mdui-card
          class="feature-card"
          variant="outlined"
        >
          <mdui-icon
            name="cloud_sync"
            class="feature-icon"
          />
          <div class="feature-title">
            🌐 Git 协议
          </div>
          <div class="feature-description">
            完整的 HTTP Git 协议支持，支持 push/pull/clone 操作
          </div>
        </mdui-card>

        <mdui-card
          class="feature-card"
          variant="outlined"
        >
          <mdui-icon
            name="speed"
            class="feature-icon"
          />
          <div class="feature-title">
            ⚡ 高性能
          </div>
          <div class="feature-description">
            C++ 实现 Git 核心操作，Go 处理业务逻辑，性能卓越
          </div>
        </mdui-card>
      </div>

      <h2 style="text-align: center; margin-bottom: 32px;">
        快速开始
      </h2>
      <mdui-card
        class="quickstart-card"
        variant="filled"
      >
        <h3 style="margin-top: 0;">
          克隆仓库
        </h3>
        <div class="code-block">
          <mdui-button-icon
            class="copy-button"
            icon="content_copy"
            @click="copyCommand"
          />
          <code>git clone http://localhost:8080/username/repo.git</code>
        </div>
      </mdui-card>

      <mdui-card
        style="padding: 24px; margin-bottom: 32px;"
        variant="outlined"
      >
        <h3 style="margin-top: 0;">
          多平台支持
        </h3>
        <mdui-list>
          <mdui-list-item icon="computer">
            <strong>macOS</strong> - 原生支持，Homebrew 依赖管理
          </mdui-list-item>
          <mdui-list-item icon="devices">
            <strong>Linux</strong> - Ubuntu/Debian, CentOS/RHEL
          </mdui-list-item>
          <mdui-list-item icon="desktop_windows">
            <strong>Windows</strong> - Visual Studio 2022, vcpkg 支持
          </mdui-list-item>
          <mdui-list-item icon="code">
            <strong>IDE 支持</strong> - VS Code, Xcode, Visual Studio, CLion
          </mdui-list-item>
        </mdui-list>
      </mdui-card>

      <div class="action-buttons">
        <mdui-button
          variant="filled"
          icon="play_arrow"
          @click="$router.push('/register')"
        >
          开始使用
        </mdui-button>
        <mdui-button
          variant="outlined"
          icon="book"
          href="https://github.com/Zixiao-System/ZiXiao-Git-Server"
          target="_blank"
        >
          查看文档
        </mdui-button>
        <mdui-button
          variant="text"
          icon="login"
          @click="$router.push('/login')"
        >
          登录
        </mdui-button>
      </div>
    </div>

    <div style="text-align: center; padding: 48px 24px; opacity: 0.7;">
      <p>© 2024 ZiXiao System. MIT License.</p>
      <p>灵感来自 GitLab, Gitea 和 Gogs</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { setTheme, snackbar } from 'mdui'

const themeIcon = ref('brightness_auto')
let currentTheme = 'auto'

function updateThemeIcon() {
  const html = document.documentElement
  const currentClass = html.className

  if (currentClass.includes('mdui-theme-dark')) {
    themeIcon.value = 'dark_mode'
  } else if (currentClass.includes('mdui-theme-light')) {
    themeIcon.value = 'light_mode'
  } else {
    themeIcon.value = 'brightness_auto'
  }
}

function toggleTheme() {
  const themes = ['auto', 'light', 'dark']
  const currentIndex = themes.indexOf(currentTheme)
  currentTheme = themes[(currentIndex + 1) % themes.length]

  setTheme(currentTheme)
  localStorage.setItem('theme', currentTheme)
  updateThemeIcon()

  snackbar({
    message: `已切换到${currentTheme === 'auto' ? '自动' : currentTheme === 'light' ? '亮色' : '暗色'}主题`,
    placement: 'bottom'
  })
}

function copyCommand() {
  const command = 'git clone http://localhost:8080/username/repo.git'
  navigator.clipboard.writeText(command).then(() => {
    snackbar({
      message: '已复制到剪贴板',
      icon: 'done',
      placement: 'bottom'
    })
  })
}

onMounted(() => {
  currentTheme = localStorage.getItem('theme') || 'auto'
  setTheme(currentTheme)
  updateThemeIcon()
})
</script>

<style scoped>
.hero-section {
  padding: 48px 24px;
  text-align: center;
  background: linear-gradient(135deg, var(--mdui-color-primary) 0%, var(--mdui-color-secondary) 100%);
  color: white;
}

.hero-title {
  font-size: 3rem;
  font-weight: 700;
  margin-bottom: 16px;
}

.hero-subtitle {
  font-size: 1.5rem;
  opacity: 0.9;
  margin-bottom: 32px;
}

.content-section {
  max-width: 1200px;
  margin: 0 auto;
  padding: 48px 24px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 24px;
  margin-bottom: 48px;
}

.stat-card {
  text-align: center;
  padding: 32px 24px;
}

.stat-value {
  font-size: 2.5rem;
  font-weight: 700;
  color: var(--mdui-color-primary);
  margin-bottom: 8px;
}

.stat-label {
  font-size: 1rem;
  opacity: 0.7;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 24px;
  margin-bottom: 48px;
}

.feature-card {
  padding: 24px;
}

.feature-icon {
  font-size: 48px;
  color: var(--mdui-color-primary);
  margin-bottom: 16px;
}

.feature-title {
  font-size: 1.25rem;
  font-weight: 600;
  margin-bottom: 8px;
}

.feature-description {
  opacity: 0.7;
  line-height: 1.6;
}

.quickstart-card {
  margin-bottom: 32px;
  padding: 24px;
}

.code-block {
  background: var(--mdui-color-surface-variant);
  padding: 16px;
  border-radius: 8px;
  font-family: 'Courier New', monospace;
  font-size: 0.95rem;
  overflow-x: auto;
  position: relative;
}

.copy-button {
  position: absolute;
  top: 8px;
  right: 8px;
}

.action-buttons {
  display: flex;
  gap: 16px;
  justify-content: center;
  flex-wrap: wrap;
}

@media (max-width: 768px) {
  .hero-title {
    font-size: 2rem;
  }

  .hero-subtitle {
    font-size: 1.25rem;
  }

  .stats-grid,
  .features-grid {
    grid-template-columns: 1fr;
  }
}
</style>
