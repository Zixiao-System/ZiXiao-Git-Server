<template>
  <div class="dashboard-layout">
    <mdui-top-app-bar>
      <mdui-button-icon
        icon="menu"
        @click="drawerOpen = !drawerOpen"
      />
      <mdui-top-app-bar-title>我的仓库</mdui-top-app-bar-title>
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
      <div class="page-header">
        <h1>我的仓库</h1>
        <mdui-button
          variant="filled"
          icon="add"
          @click="showCreateDialog = true"
        >
          创建仓库
        </mdui-button>
      </div>

      <mdui-text-field
        v-model="searchQuery"
        placeholder="搜索仓库..."
        icon="search"
        style="margin-bottom: 24px; width: 100%;"
      />

      <div
        v-if="loading"
        style="text-align: center; padding: 48px;"
      >
        <mdui-circular-progress />
      </div>

      <div
        v-else-if="filteredRepositories.length > 0"
        class="repo-grid"
      >
        <mdui-card
          v-for="repo in filteredRepositories"
          :key="repo.id"
          class="repo-card"
        >
          <div class="repo-header">
            <mdui-icon
              name="folder"
              class="repo-icon"
            />
            <div class="repo-title">
              <h3>{{ repo.name }}</h3>
              <mdui-chip v-if="repo.is_public">
                公开
              </mdui-chip>
              <mdui-chip v-else>
                私有
              </mdui-chip>
            </div>
          </div>

          <p class="repo-description">
            {{ repo.description || '无描述' }}
          </p>

          <div class="repo-meta">
            <span><mdui-icon name="schedule" /> {{ formatDate(repo.created_at) }}</span>
          </div>

          <div class="repo-actions">
            <mdui-button
              variant="text"
              icon="info"
              @click="$router.push(`/repositories/${repo.id}`)"
            >
              详情
            </mdui-button>
            <mdui-button
              variant="text"
              icon="delete"
              @click="confirmDelete(repo)"
            >
              删除
            </mdui-button>
          </div>
        </mdui-card>
      </div>

      <mdui-card
        v-else
        style="padding: 48px; text-align: center;"
      >
        <mdui-icon
          name="inbox"
          style="font-size: 64px; opacity: 0.3;"
        />
        <p>{{ searchQuery ? '没有找到匹配的仓库' : '还没有仓库，创建第一个仓库吧！' }}</p>
      </mdui-card>
    </div>

    <!-- Create Repository Dialog -->
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
            helper="只能包含字母、数字、连字符和下划线"
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
          :loading="creating"
          @click="handleCreateRepo"
        >
          创建
        </mdui-button>
      </mdui-dialog-actions>
    </mdui-dialog>

    <!-- Delete Confirmation Dialog -->
    <mdui-dialog
      :open="showDeleteDialog"
      @close="showDeleteDialog = false"
    >
      <mdui-dialog-headline>确认删除</mdui-dialog-headline>
      <mdui-dialog-body>
        确定要删除仓库 <strong>{{ repoToDelete?.name }}</strong> 吗？此操作不可撤销。
      </mdui-dialog-body>
      <mdui-dialog-actions>
        <mdui-button @click="showDeleteDialog = false">
          取消
        </mdui-button>
        <mdui-button
          variant="filled"
          :loading="deleting"
          @click="handleDeleteRepo"
        >
          删除
        </mdui-button>
      </mdui-dialog-actions>
    </mdui-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useRepositoryStore } from '@/stores/repository'
import { snackbar } from 'mdui'

const router = useRouter()
const authStore = useAuthStore()
const repoStore = useRepositoryStore()

const drawerOpen = ref(false)
const showCreateDialog = ref(false)
const showDeleteDialog = ref(false)
const searchQuery = ref('')
const loading = ref(false)
const creating = ref(false)
const deleting = ref(false)
const repoToDelete = ref(null)

const newRepo = ref({
  name: '',
  description: '',
  is_public: false
})

const filteredRepositories = computed(() => {
  if (!searchQuery.value) {
    return repoStore.repositories
  }
  const query = searchQuery.value.toLowerCase()
  return repoStore.repositories.filter(repo =>
    repo.name.toLowerCase().includes(query) ||
    (repo.description && repo.description.toLowerCase().includes(query))
  )
})

async function loadRepositories() {
  loading.value = true
  try {
    await repoStore.fetchRepositories()
  } catch (error) {
    snackbar({
      message: '加载仓库失败',
      placement: 'top'
    })
  } finally {
    loading.value = false
  }
}

async function handleCreateRepo() {
  creating.value = true
  try {
    await repoStore.createRepository(newRepo.value)
    snackbar({
      message: '仓库创建成功！',
      icon: 'done',
      placement: 'top'
    })
    showCreateDialog.value = false
    newRepo.value = { name: '', description: '', is_public: false }
  } catch (error) {
    snackbar({
      message: error.error || '创建仓库失败',
      placement: 'top'
    })
  } finally {
    creating.value = false
  }
}

function confirmDelete(repo) {
  repoToDelete.value = repo
  showDeleteDialog.value = true
}

async function handleDeleteRepo() {
  deleting.value = true
  try {
    await repoStore.deleteRepository(repoToDelete.value.id)
    snackbar({
      message: '仓库已删除',
      icon: 'done',
      placement: 'top'
    })
    showDeleteDialog.value = false
    repoToDelete.value = null
  } catch (error) {
    snackbar({
      message: error.error || '删除仓库失败',
      placement: 'top'
    })
  } finally {
    deleting.value = false
  }
}

function formatDate(dateString) {
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN')
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

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
}

.page-header h1 {
  margin: 0;
}

.repo-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 24px;
}

.repo-card {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.repo-header {
  display: flex;
  gap: 12px;
  align-items: flex-start;
}

.repo-icon {
  font-size: 32px;
  color: var(--mdui-color-primary);
}

.repo-title {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.repo-title h3 {
  margin: 0;
  font-size: 1.25rem;
}

.repo-description {
  margin: 0;
  opacity: 0.7;
  min-height: 40px;
}

.repo-meta {
  display: flex;
  gap: 16px;
  font-size: 0.9rem;
  opacity: 0.7;
  align-items: center;
}

.repo-meta span {
  display: flex;
  align-items: center;
  gap: 4px;
}

.repo-actions {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
  border-top: 1px solid var(--mdui-color-outline-variant);
  padding-top: 16px;
}
</style>
