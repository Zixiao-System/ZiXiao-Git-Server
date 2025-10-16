import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'

// Import MDUI
import 'mdui/mdui.css'
import 'mdui'

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount('#app')
