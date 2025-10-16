<template>
  <div class="auth-container">
    <mdui-card
      class="auth-card"
      variant="elevated"
    >
      <div class="auth-header">
        <h1>登录</h1>
        <p>欢迎回到 ZiXiao Git Server</p>
      </div>

      <form @submit.prevent="handleLogin">
        <mdui-text-field
          v-model="formData.username"
          label="用户名"
          placeholder="请输入用户名"
          required
          :disabled="loading"
        />

        <mdui-text-field
          v-model="formData.password"
          type="password"
          label="密码"
          placeholder="请输入密码"
          required
          :disabled="loading"
          toggle-password
        />

        <mdui-button
          type="submit"
          variant="filled"
          full-width
          :loading="loading"
          style="margin-top: 24px;"
        >
          登录
        </mdui-button>
      </form>

      <div class="auth-footer">
        <p>还没有账号？ <a @click="$router.push('/register')">立即注册</a></p>
        <p><a @click="$router.push('/')">返回首页</a></p>
      </div>
    </mdui-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { snackbar } from 'mdui'

const router = useRouter()
const authStore = useAuthStore()

const formData = ref({
  username: '',
  password: ''
})

const loading = ref(false)

async function handleLogin() {
  loading.value = true
  try {
    await authStore.login(formData.value.username, formData.value.password)
    snackbar({
      message: '登录成功！',
      icon: 'done',
      placement: 'top'
    })
    router.push('/dashboard')
  } catch (error) {
    snackbar({
      message: error.error || '登录失败，请检查用户名和密码',
      placement: 'top'
    })
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, var(--mdui-color-primary-container) 0%, var(--mdui-color-secondary-container) 100%);
  padding: 24px;
}

.auth-card {
  width: 100%;
  max-width: 400px;
  padding: 32px;
}

.auth-header {
  text-align: center;
  margin-bottom: 32px;
}

.auth-header h1 {
  margin: 0 0 8px 0;
  font-size: 2rem;
}

.auth-header p {
  margin: 0;
  opacity: 0.7;
}

mdui-text-field {
  margin-bottom: 16px;
  width: 100%;
}

.auth-footer {
  margin-top: 24px;
  text-align: center;
}

.auth-footer p {
  margin: 8px 0;
  font-size: 0.9rem;
}

.auth-footer a {
  color: var(--mdui-color-primary);
  cursor: pointer;
  text-decoration: none;
}

.auth-footer a:hover {
  text-decoration: underline;
}
</style>
