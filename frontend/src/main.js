import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'

// Import MDUI CSS and register all components
import 'mdui/mdui.css'
import 'mdui/components/button.js'
import 'mdui/components/button-icon.js'
import 'mdui/components/card.js'
import 'mdui/components/checkbox.js'
import 'mdui/components/chip.js'
import 'mdui/components/circular-progress.js'
import 'mdui/components/dialog.js'
import 'mdui/components/icon.js'
import 'mdui/components/list.js'
import 'mdui/components/navigation-drawer.js'
import 'mdui/components/text-field.js'
import 'mdui/components/top-app-bar.js'

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount('#app')
