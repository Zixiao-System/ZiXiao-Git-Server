<template>
  <div class="dashboard-layout">
    <mdui-top-app-bar>
      <mdui-button-icon
        icon="arrow_back"
        @click="$router.back()"
      />
      <mdui-top-app-bar-title>仓库详情</mdui-top-app-bar-title>
      <div style="flex-grow: 1" />
      <mdui-button-icon
        icon="settings"
        @click="showEditDialog = true"
      />
    </mdui-top-app-bar>

    <div class="content-area">
      <div
        v-if="loading"
        style="text-align: center; padding: 48px;"
      >
        <mdui-circular-progress />
      </div>

      <div v-else-if="repository">
        <div class="repo-header">
          <div>
            <h1><mdui-icon name="folder" /> {{ repository.name }}</h1>
            <p class="repo-description">
              {{ repository.description || '无描述' }}
            </p>
            <div class="repo-badges">
              <mdui-chip v-if="repository.is_public">
                公开
              </mdui-chip>
              <mdui-chip v-else>
                私有
              </mdui-chip>
            </div>
          </div>
        </div>

        <mdui-card style="padding: 24px; margin-top: 24px;">
          <h2>克隆仓库</h2>
          <div class="code-block">
            <code>git clone http://localhost:8080/{{ authStore.user?.username }}/{{ repository.name }}.git</code>
            <mdui-button-icon
              icon="content_copy"
              @click="copyCloneCommand"
            />
          </div>
        </mdui-card>

        <mdui-card style="padding: 24px; margin-top: 24px;">
          <h2>仓库信息</h2>
          <mdui-list>
            <mdui-list-item>
              <strong>创建时间:</strong> {{ formatDate(repository.created_at) }}
            </mdui-list-item>
            <mdui-list-item>
              <strong>更新时间:</strong> {{ formatDate(repository.updated_at) }}
            </mdui-list-item>
            <mdui-list-item>
              <strong>所有者:</strong> {{ authStore.user?.username }}
            </mdui-list-item>
          </mdui-list>
        </mdui-card>

        <mdui-card style="padding: 24px; margin-top: 24px;">
          <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px;">
            <h2>协作者</h2>
            <mdui-button
              variant="outlined"
              icon="person_add"
            >
              添加协作者
            </mdui-button>
          </div>
          <p style="opacity: 0.7;">
            暂无协作者
          </p>
        </mdui-card>
      </div>

      <mdui-card
        v-else
        style="padding: 48px; text-align: center;"
      >
        <p>仓库不存在或已被删除</p>
        <mdui-button @click="$router.push('/repositories')">
          返回仓库列表
        </mdui-button>
      </mdui-card>
    </div>

    <!-- Edit Repository Dialog -->
    <mdui-dialog
      :open="showEditDialog"
      @close="showEditDialog = false"
    >
      <mdui-dialog-headline>编辑仓库</mdui-dialog-headline>
      <mdui-dialog-body>
        <form @submit.prevent="handleUpdateRepo">
          <mdui-text-field
            v-model="editForm.name"
            label="仓库名称"
            required
            disabled
          />
          <mdui-text-field
            v-model="editForm.description"
            label="描述"
            style="margin-top: 16px;"
          />
          <mdui-checkbox
            v-model="editForm.is_public"
            style="margin-top: 16px;"
          >
            公开仓库
          </mdui-checkbox>
        </form>
      </mdui-dialog-body>
      <mdui-dialog-actions>
        <mdui-button @click="showEditDialog = false">
          取消
        </mdui-button>
        <mdui-button
          variant="filled"
          :loading="updating"
          @click="handleUpdateRepo"
        >
          保存
        </mdui-button>
      </mdui-dialog-actions>
    </mdui-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useRepositoryStore } from '@/stores/repository'
import { snackbar } from 'mdui'

const route = useRoute()
const authStore = useAuthStore()
const repoStore = useRepositoryStore()

const loading = ref(false)
const updating = ref(false)
const showEditDialog = ref(false)
const repository = ref(null)

const editForm = ref({
  name: '',
  description: '',
  is_public: false
})

async function loadRepository() {
  loading.value = true
  try {
    repository.value = await repoStore.fetchRepository(route.params.id)
    editForm.value = {
      name: repository.value.name,
      description: repository.value.description,
      is_public: repository.value.is_public
    }
  } catch (error) {
    snackbar({
      message: '加载仓库失败',
      placement: 'top'
    })
  } finally {
    loading.value = false
  }
}

async function handleUpdateRepo() {
  updating.value = true
  try {
    await repoStore.updateRepository(route.params.id, editForm.value)
    repository.value = repoStore.currentRepository
    snackbar({
      message: '仓库更新成功！',
      icon: 'done',
      placement: 'top'
    })
    showEditDialog.value = false
  } catch (error) {
    snackbar({
      message: error.error || '更新仓库失败',
      placement: 'top'
    })
  } finally {
    updating.value = false
  }
}

function copyCloneCommand() {
  const command = `git clone http://localhost:8080/${authStore.user?.username}/${repository.value.name}.git`
  navigator.clipboard.writeText(command).then(() => {
    snackbar({
      message: '已复制到剪贴板',
      icon: 'done',
      placement: 'bottom'
    })
  })
}

function formatDate(dateString) {
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

onMounted(() => {
  loadRepository()
})
</script>

<style scoped>
.dashboard-layout {
  min-height: 100vh;
  background: var(--mdui-color-surface);
}

.content-area {
  padding: 24px;
  max-width: 900px;
  margin: 0 auto;
  margin-top: 64px;
}

.repo-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.repo-header h1 {
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 0 0 8px 0;
}

.repo-description {
  opacity: 0.7;
  margin: 8px 0;
}

.repo-badges {
  display: flex;
  gap: 8px;
  margin-top: 16px;
}

.code-block {
  background: var(--mdui-color-surface-variant);
  padding: 16px;
  border-radius: 8px;
  font-family: 'Courier New', monospace;
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 16px;
}

.code-block code {
  flex: 1;
  overflow-x: auto;
}
</style>
