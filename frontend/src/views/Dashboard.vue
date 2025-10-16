<template>
  <div class="dashboard-layout">
    <mdui-top-app-bar>
      <mdui-button-icon
        icon="menu"
        @click="drawerOpen = !drawerOpen"
      />
      <mdui-top-app-bar-title>仪表盘</mdui-top-app-bar-title>
      <div style="flex-grow: 1" />
      <mdui-button-icon
        icon="logout"
        @click="handleLogout"
      />
    </mdui-top-app-bar>

    <mdui-navigation-drawer
      :open="drawerOpen"
      @close="drawerOpen = false"
    >
      <mdui-list>
        <mdui-list-item
          icon="dashboard"
          @click="$router.push('/dashboard')"
        >
          仪表盘
        </mdui-list-item>
        <mdui-list-item
          icon="folder"
          @click="$router.push('/repositories')"
        >
          我的仓库
        </mdui-list-item>
        <mdui-list-item icon="settings">
          设置
        </mdui-list-item>
      </mdui-list>
    </mdui-navigation-drawer>

    <div class="content-area">
      <div class="dashboard-header">
        <h1>欢迎, {{ authStore.user?.username }}!</h1>
        <mdui-button
          variant="filled"
          icon="add"
          @click="showCreateDialog = true"
        >
          创建仓库
        </mdui-button>
      </div>

      <div class="stats-grid">
        <mdui-card class="stat-card">
          <mdui-icon
            name="folder"
            style="font-size: 48px; color: var(--mdui-color-primary);"
          />
          <div class="stat-value">
            {{ repositories.length }}
          </div>
          <div class="stat-label">
            仓库数量
          </div>
        </mdui-card>

        <mdui-card class="stat-card">
          <mdui-icon
            name="people"
            style="font-size: 48px; color: var(--mdui-color-secondary);"
          />
          <div class="stat-value">
            0
          </div>
          <div class="stat-label">
            协作者
          </div>
        </mdui-card>

        <mdui-card class="stat-card">
          <mdui-icon
            name="trending_up"
            style="font-size: 48px; color: var(--mdui-color-tertiary);"
          />
          <div class="stat-value">
            0
          </div>
          <div class="stat-label">
            活动
          </div>
        </mdui-card>
      </div>

      <h2>最近的仓库</h2>
      <div
        v-if="repositories.length > 0"
        class="repo-list"
      >
        <mdui-card
          v-for="repo in repositories.slice(0, 5)"
          :key="repo.id"
          class="repo-card"
        >
          <div class="repo-info">
            <mdui-icon name="folder" />
            <div>
              <h3>{{ repo.name }}</h3>
              <p>{{ repo.description || '无描述' }}</p>
            </div>
          </div>
          <mdui-button
            variant="text"
            @click="$router.push(`/repositories/${repo.id}`)"
          >
            查看详情
          </mdui-button>
        </mdui-card>
      </div>
      <mdui-card
        v-else
        style="padding: 32px; text-align: center;"
      >
        <p>还没有仓库，创建第一个仓库吧！</p>
      </mdui-card>
    </div>

    <mdui-dialog
      :open="showCreateDialog"
      @close="showCreateDialog = false"
    >
      <mdui-dialog-headline>创建新仓库</mdui-dialog-headline>
      <mdui-dialog-body>
        <form @submit.prevent="handleCreateRepo">
          <mdui-text-field
            v-model="newRepo.name"
            label="仓库名称"
            required
          />
          <mdui-text-field
            v-model="newRepo.description"
            label="描述"
            style="margin-top: 16px;"
          />
          <mdui-checkbox
            v-model="newRepo.is_public"
            style="margin-top: 16px;"
          >
            公开仓库
          </mdui-checkbox>
        </form>
      </mdui-dialog-body>
      <mdui-dialog-actions>
        <mdui-button @click="showCreateDialog = false">
          取消
        </mdui-button>
        <mdui-button
          variant="filled"
          @click="handleCreateRepo"
        >
          创建
        </mdui-button>
      </mdui-dialog-actions>
    </mdui-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useRepositoryStore } from '@/stores/repository'
import { snackbar } from 'mdui'

const router = useRouter()
const authStore = useAuthStore()
const repoStore = useRepositoryStore()

const drawerOpen = ref(false)
const showCreateDialog = ref(false)
const repositories = ref([])

const newRepo = ref({
  name: '',
  description: '',
  is_public: false
})

async function loadRepositories() {
  try {
    repositories.value = await repoStore.fetchRepositories()
  } catch (error) {
    snackbar({
      message: '加载仓库失败',
      placement: 'top'
    })
  }
}

async function handleCreateRepo() {
  try {
    await repoStore.createRepository(newRepo.value)
    snackbar({
      message: '仓库创建成功！',
      icon: 'done',
      placement: 'top'
    })
    showCreateDialog.value = false
    newRepo.value = { name: '', description: '', is_public: false }
    await loadRepositories()
  } catch (error) {
    snackbar({
      message: error.error || '创建仓库失败',
      placement: 'top'
    })
  }
}

function handleLogout() {
  authStore.logout()
  snackbar({
    message: '已退出登录',
    placement: 'top'
  })
  router.push('/')
}

onMounted(() => {
  loadRepositories()
})
</script>

<style scoped>
.dashboard-layout {
  min-height: 100vh;
  background: var(--mdui-color-surface);
}

.content-area {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
  margin-top: 64px;
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
}

.dashboard-header h1 {
  margin: 0;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 24px;
  margin-bottom: 48px;
}

.stat-card {
  padding: 24px;
  text-align: center;
}

.stat-value {
  font-size: 2rem;
  font-weight: 700;
  color: var(--mdui-color-primary);
  margin: 8px 0;
}

.stat-label {
  opacity: 0.7;
}

.repo-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-top: 16px;
}

.repo-card {
  padding: 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.repo-info {
  display: flex;
  gap: 16px;
  align-items: center;
}

.repo-info h3 {
  margin: 0 0 4px 0;
}

.repo-info p {
  margin: 0;
  opacity: 0.7;
  font-size: 0.9rem;
}
</style>
