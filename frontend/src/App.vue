<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { setTheme } from 'mdui'

onMounted(() => {
  // Initialize MDUI theme
  const savedTheme = localStorage.getItem('theme') || 'auto'
  setTheme(savedTheme)
})
</script>

<style>
body {
  margin: 0;
  padding: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
}

* {
  box-sizing: border-box;
}
</style>
