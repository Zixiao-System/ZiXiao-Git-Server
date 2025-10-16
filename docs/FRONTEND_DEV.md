# ZiXiao Git Server - Frontend Development Guide

## Project Structure

```
frontend/
├── public/              # Static assets
├── src/
│   ├── assets/         # Images, styles, etc.
│   ├── components/     # Reusable Vue components
│   ├── views/          # Page components
│   │   ├── Home.vue           # Landing page
│   │   ├── Login.vue          # Login page
│   │   ├── Register.vue       # Registration page
│   │   ├── Dashboard.vue      # User dashboard
│   │   ├── Repositories.vue   # Repository list
│   │   ├── RepositoryDetail.vue  # Repository details
│   │   └── NotFound.vue       # 404 page
│   ├── router/         # Vue Router configuration
│   │   └── index.js
│   ├── stores/         # Pinia state management
│   │   ├── auth.js     # Authentication store
│   │   └── repository.js  # Repository store
│   ├── services/       # API service layer
│   │   └── index.js    # API functions
│   ├── utils/          # Utility functions
│   │   └── api.js      # Axios configuration
│   ├── App.vue         # Root component
│   └── main.js         # Application entry point
├── index.html          # HTML template
├── vite.config.js      # Vite configuration
├── package.json        # Dependencies
└── .env.*              # Environment variables
```

## Technology Stack

- **Vue 3**: Progressive JavaScript framework with Composition API
- **Vue Router**: Official router for Vue.js
- **Pinia**: State management library (Vuex replacement)
- **Axios**: HTTP client for API requests
- **Vite**: Next-generation frontend build tool
- **MDUI 2.x**: Material Design UI components

## Getting Started

### Installation

```bash
cd frontend
npm install
```

### Development Server

```bash
npm run dev
```

Runs on `http://localhost:3000` with hot module replacement.

### Build for Production

```bash
npm run build
```

Outputs to `../web/dist/`.

### Preview Production Build

```bash
npm run preview
```

## State Management (Pinia)

### Auth Store (`stores/auth.js`)

Manages user authentication state:

```javascript
const authStore = useAuthStore()

// Login
await authStore.login(username, password)

// Register
await authStore.register(username, password, email)

// Logout
authStore.logout()

// Check authentication
if (authStore.isAuthenticated) {
  // User is logged in
}

// Access user data
console.log(authStore.user)
```

### Repository Store (`stores/repository.js`)

Manages repository data:

```javascript
const repoStore = useRepositoryStore()

// Fetch all repositories
await repoStore.fetchRepositories()

// Fetch single repository
await repoStore.fetchRepository(id)

// Create repository
await repoStore.createRepository({ name, description, is_public })

// Update repository
await repoStore.updateRepository(id, { description, is_public })

// Delete repository
await repoStore.deleteRepository(id)

// Access data
console.log(repoStore.repositories)
console.log(repoStore.currentRepository)
```

## API Service Layer

### Services (`services/index.js`)

Centralized API functions:

```javascript
import { authService, repositoryService } from '@/services'

// Authentication
await authService.login(username, password)
await authService.register(username, password, email)
await authService.getCurrentUser()

// Repositories
await repositoryService.getRepositories()
await repositoryService.getRepository(id)
await repositoryService.createRepository(data)
await repositoryService.updateRepository(id, data)
await repositoryService.deleteRepository(id)
await repositoryService.getCollaborators(id)
await repositoryService.addCollaborator(id, data)
```

### Axios Configuration (`utils/api.js`)

Configured with:
- Base URL from environment variables
- JWT token auto-attachment
- Response interceptors for error handling
- 401 redirect to login

## Routing

### Routes

| Path | Component | Auth Required |
|------|-----------|---------------|
| `/` | Home | No |
| `/login` | Login | No |
| `/register` | Register | No |
| `/dashboard` | Dashboard | Yes |
| `/repositories` | Repositories | Yes |
| `/repositories/:id` | RepositoryDetail | Yes |
| `*` (404) | NotFound | No |

### Navigation Guards

Automatic redirect:
- Unauthenticated users accessing protected routes → Login
- Authenticated users accessing login/register → Dashboard

## UI Components (MDUI 2.x)

### Common Components

```vue
<!-- Buttons -->
<mdui-button variant="filled">Primary Button</mdui-button>
<mdui-button variant="outlined">Secondary Button</mdui-button>
<mdui-button variant="text">Text Button</mdui-button>

<!-- Text Fields -->
<mdui-text-field
  v-model="value"
  label="Label"
  placeholder="Placeholder"
  required
></mdui-text-field>

<!-- Cards -->
<mdui-card variant="elevated">
  Content
</mdui-card>

<!-- Dialogs -->
<mdui-dialog :open="showDialog" @close="showDialog = false">
  <mdui-dialog-headline>Title</mdui-dialog-headline>
  <mdui-dialog-body>Content</mdui-dialog-body>
  <mdui-dialog-actions>
    <mdui-button @click="showDialog = false">Close</mdui-button>
  </mdui-dialog-actions>
</mdui-dialog>

<!-- Snackbar (Toast) -->
<script>
import mdui from 'mdui'

mdui.snackbar({
  message: 'Success!',
  icon: 'done',
  placement: 'top'
})
</script>
```

### Theme System

```javascript
import mdui from 'mdui'

// Set theme
mdui.setTheme('light')  // 'light', 'dark', or 'auto'

// Get current theme
const theme = document.documentElement.className
```

## Environment Variables

### `.env` (Default)

```env
VITE_API_BASE_URL=/api/v1
```

### `.env.development`

```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

### `.env.production`

```env
VITE_API_BASE_URL=/api/v1
```

### Usage

```javascript
const apiUrl = import.meta.env.VITE_API_BASE_URL
```

## Code Style

### Vue 3 Composition API

```vue
<script setup>
import { ref, onMounted } from 'vue'

const count = ref(0)

function increment() {
  count.value++
}

onMounted(() => {
  console.log('Component mounted')
})
</script>
```

### Naming Conventions

- **Components**: PascalCase (e.g., `UserProfile.vue`)
- **Composables**: camelCase with `use` prefix (e.g., `useAuth`)
- **Files**: kebab-case for utilities (e.g., `api-client.js`)
- **Variables**: camelCase

## Best Practices

### 1. Error Handling

```javascript
try {
  await repoStore.createRepository(data)
  mdui.snackbar({ message: 'Success!', icon: 'done' })
} catch (error) {
  mdui.snackbar({ message: error.error || 'Failed' })
}
```

### 2. Loading States

```vue
<script setup>
const loading = ref(false)

async function loadData() {
  loading.value = true
  try {
    await fetchData()
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <mdui-button :loading="loading" @click="loadData">
    Load
  </mdui-button>
</template>
```

### 3. Form Validation

```vue
<script setup>
const form = ref({
  username: '',
  password: ''
})

function validate() {
  if (!form.value.username) {
    mdui.snackbar({ message: 'Username is required' })
    return false
  }
  if (form.value.password.length < 8) {
    mdui.snackbar({ message: 'Password must be at least 8 characters' })
    return false
  }
  return true
}

async function submit() {
  if (!validate()) return
  // Submit form
}
</script>
```

## Testing

### Unit Tests (TODO)

```bash
npm run test:unit
```

### E2E Tests (TODO)

```bash
npm run test:e2e
```

## Build Optimization

### Vite Configuration

```javascript
// vite.config.js
export default defineConfig({
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          'vendor': ['vue', 'vue-router', 'pinia'],
          'mdui': ['mdui']
        }
      }
    }
  }
})
```

### Code Splitting

Vue Router automatically code-splits routes:

```javascript
{
  path: '/dashboard',
  component: () => import('@/views/Dashboard.vue')
}
```

## Troubleshooting

### Issue: "Cannot find module '@/...'"

**Solution**: Check `vite.config.js` alias configuration:

```javascript
resolve: {
  alias: {
    '@': fileURLToPath(new URL('./src', import.meta.url))
  }
}
```

### Issue: API requests failing in development

**Solution**: Check Vite proxy configuration:

```javascript
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true
    }
  }
}
```

### Issue: MDUI components not rendering

**Solution**: Ensure MDUI is imported in `main.js`:

```javascript
import 'mdui/mdui.css'
import 'mdui'
```

And configure Vite to treat MDUI tags as custom elements:

```javascript
vue({
  template: {
    compilerOptions: {
      isCustomElement: (tag) => tag.startsWith('mdui-')
    }
  }
})
```

## Resources

- [Vue 3 Documentation](https://vuejs.org/)
- [Pinia Documentation](https://pinia.vuejs.org/)
- [Vue Router Documentation](https://router.vuejs.org/)
- [Vite Documentation](https://vitejs.dev/)
- [MDUI Documentation](https://www.mdui.org/)
