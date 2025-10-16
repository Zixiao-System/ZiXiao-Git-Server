<template>
  <div class="not-found">
    <mdui-card class="not-found-card">
      <mdui-icon
        name="error_outline"
        style="font-size: 96px; opacity: 0.5;"
      />
      <h1>404</h1>
      <p>页面未找到</p>
      <mdui-button
        variant="filled"
        icon="home"
        @click="$router.push('/')"
      >
        返回首页
      </mdui-button>
    </mdui-card>
  </div>
</template>

<script setup>
</script>

<style scoped>
.not-found {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--mdui-color-surface);
}

.not-found-card {
  padding: 48px;
  text-align: center;
  max-width: 400px;
}

.not-found-card h1 {
  font-size: 4rem;
  margin: 16px 0;
  color: var(--mdui-color-primary);
}

.not-found-card p {
  opacity: 0.7;
  margin-bottom: 24px;
}
</style>
