<template>
  <div class="auth-container">
    <mdui-card class="auth-card" variant="elevated">
      <div class="auth-header">
        <h1>注册</h1>
        <p>创建你的 ZiXiao Git Server 账号</p>
      </div>

      <form @submit.prevent="handleRegister">
        <mdui-text-field
          v-model="formData.username"
          label="用户名"
          placeholder="请输入用户名"
          required
          :disabled="loading"
        ></mdui-text-field>

        <mdui-text-field
          v-model="formData.email"
          type="email"
          label="邮箱"
          placeholder="请输入邮箱"
          required
          :disabled="loading"
        ></mdui-text-field>

        <mdui-text-field
          v-model="formData.password"
          type="password"
          label="密码"
          placeholder="请输入密码（至少8位）"
          required
          :disabled="loading"
          toggle-password
        ></mdui-text-field>

        <mdui-text-field
          v-model="formData.confirmPassword"
          type="password"
          label="确认密码"
          placeholder="请再次输入密码"
          required
          :disabled="loading"
          toggle-password
        ></mdui-text-field>

        <mdui-button
          type="submit"
          variant="filled"
          full-width
          :loading="loading"
          style="margin-top: 24px;"
        >
          注册
        </mdui-button>
      </form>

      <div class="auth-footer">
        <p>已有账号？ <a @click="$router.push('/login')">立即登录</a></p>
        <p><a @click="$router.push('/')">返回首页</a></p>
      </div>
    </mdui-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import mdui from 'mdui'

const router = useRouter()
const authStore = useAuthStore()

const formData = ref({
  username: '',
  email: '',
  password: '',
  confirmPassword: ''
})

const loading = ref(false)

async function handleRegister() {
  if (formData.value.password.length < 8) {
    mdui.snackbar({
      message: '密码长度至少为8位',
      placement: 'top'
    })
    return
  }

  if (formData.value.password !== formData.value.confirmPassword) {
    mdui.snackbar({
      message: '两次输入的密码不一致',
      placement: 'top'
    })
    return
  }

  loading.value = true
  try {
    await authStore.register(
      formData.value.username,
      formData.value.password,
      formData.value.email
    )
    mdui.snackbar({
      message: '注册成功！请登录',
      icon: 'done',
      placement: 'top'
    })
    router.push('/login')
  } catch (error) {
    mdui.snackbar({
      message: error.error || '注册失败，请重试',
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
