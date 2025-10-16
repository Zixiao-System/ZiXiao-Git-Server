# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ZiXiao Git Server is a lightweight, high-performance Git server with a hybrid architecture:
- **Vue 3 Frontend**: Modern SPA with Vue Router, Pinia state management, MDUI 2.x components
- **Nginx**: Reverse proxy serving frontend and proxying API/Git requests
- **Go Backend**: HTTP API, business logic, authentication, database (port 8080)
- **CGo Bridge**: C API wrapper enabling Go to call C++ functions
- **C++ Core**: Low-level Git operations (repository, objects, protocol, pack files)

**Tech Stack**: Vue 3, Vite, Pinia, Vue Router, Nginx, Go 1.21+, C++17, Gin framework, JWT auth, SQLite3/PostgreSQL/SQL Server, MDUI 2.x

## Common Commands

```bash
# Backend
make build          # Build C++ library + Go server
make build-cpp      # Build C++ library only
make build-go       # Build Go server only
make run            # Build and run Go server (port 8080)
make clean          # Clean build artifacts
make test           # Run Go tests

# Frontend
cd frontend
npm install         # Install dependencies
npm run dev         # Dev server (port 3000)
npm run build       # Production build to web/dist/
npm run preview     # Preview production build
npm run lint        # Lint Vue files

# Build frontend (from project root)
./scripts/build-frontend.sh        # Linux/macOS
./scripts/build-frontend.ps1       # Windows PowerShell

# Full stack setup
./scripts/install.sh               # Install all dependencies (macOS/Linux)
./scripts/install.ps1              # Install all dependencies (Windows)
./scripts/api-test.sh              # Test API endpoints
```

## Architecture Deep Dive

### Four-Layer Full Stack Architecture

```
┌─────────────────────────────────────┐
│   Nginx (Port 80)                   │  ← Serves Vue SPA, proxies API/Git
│   configs/nginx-*.conf              │
└─────────┬───────────────────────────┘
          │
          ├──────────────────┐
          │                  │
┌─────────▼───────────────┐  │
│   Vue 3 SPA             │  │         ← Vue Router, Pinia stores, MDUI components
│   frontend/src/         │  │           Built to web/dist/ by Vite
│   - views/              │  │
│   - stores/             │  │
│   - services/           │  │
└─────────────────────────┘  │
                             │
                  ┌──────────▼─────────────────┐
                  │   Go HTTP API (Gin:8080)   │  ← REST API, JWT Auth, CORS
                  │   internal/api/            │
                  └──────────┬─────────────────┘
                             │
                  ┌──────────▼─────────────────┐
                  │   Go Business Logic        │  ← Repository CRUD, User Mgmt
                  │   internal/{repo, auth}    │     Database Ops (SQLite)
                  └──────────┬─────────────────┘
                             │
                  ┌──────────▼─────────────────┐
                  │   CGo Bridge Layer         │  ← pkg/gitcore/gitcore.go
                  │   #cgo directives          │     Type conversions (Go ↔ C)
                  └──────────┬─────────────────┘
                             │
                  ┌──────────▼─────────────────┐
                  │   C API Wrapper (C ABI)    │  ← git-core/include/git_c_api.h
                  │   extern "C" functions     │     C-compatible signatures
                  └──────────┬─────────────────┘
                             │
                  ┌──────────▼─────────────────┐
                  │   C++ Git Core Library     │  ← git-core/src/*.cpp
                  │   GitRepository, etc       │     Git operations (libgit2-style)
                  └────────────────────────────┘
```

### Request Flow Examples

**Frontend Page Load**:
1. Browser → `http://localhost/` → Nginx
2. Nginx serves `web/dist/index.html`
3. Vue app loads, initializes router and Pinia stores

**API Request** (e.g., create repository):
1. Vue component → `axios.post('/api/v1/repos')` → Nginx
2. Nginx proxies → `http://localhost:8080/api/v1/repos` → Go Gin
3. Gin handler validates JWT → calls business logic → returns JSON
4. Response → Nginx → Vue (Pinia store updates state)

**Git Operation** (e.g., git push):
1. Git client → `http://localhost/owner/repo.git/git-receive-pack` → Nginx
2. Nginx proxies → Go Gin → CGo bridge → C++ Git core
3. C++ processes pack file, updates refs, returns status

### CGo Bridge Pattern

**Location**: `pkg/gitcore/gitcore.go`

The CGo bridge is critical to understanding how Go and C++ communicate:

```go
/*
#cgo CXXFLAGS: -std=c++17 -I${SRCDIR}/../git-core/include
#cgo LDFLAGS: -L${SRCDIR}/../git-core/lib -lgitcore -lstdc++ -lz -lcrypto
#include "git_c_api.h"
*/
import "C"

type Repository struct {
    ptr unsafe.Pointer  // Holds C++ GitRepository* pointer
}

func Init(path string) (*Repository, error) {
    cPath := C.CString(path)
    defer C.free(unsafe.Pointer(cPath))

    ptr := C.git_repository_init(cPath)
    if ptr == nil {
        return nil, errors.New("failed to init repository")
    }

    return &Repository{ptr: ptr}, nil
}
```

**Key Pattern**:
1. Go calls C function via `C.git_repository_init()`
2. C API wrapper (`git_c_api.cpp`) converts to C++ call
3. C++ core (`git_repository.cpp`) performs actual Git operation
4. Result travels back up the chain

### Repository Initialization Flow

When a repository is created via `POST /api/v1/repos`:

1. **API Handler** (`internal/api/repo_handlers.go:CreateRepository`)
   - Validates JWT token via `AuthMiddleware`
   - Extracts user ID from context
   - Validates repository name and description

2. **Business Logic** (`internal/repository/repository.go:CreateRepository`)
   - Checks if repository already exists
   - Inserts record into SQLite database
   - Constructs filesystem path: `./data/repositories/{owner}/{repo}.git`

3. **CGo Bridge** (`pkg/gitcore/gitcore.go:Init`)
   - Converts Go string to C string
   - Calls C API function
   - Returns Go-wrapped repository handle

4. **C++ Core** (`git-core/src/git_repository.cpp:Init`)
   - Creates bare repository directory structure
   - Initializes `.git` subdirectories (objects, refs, hooks, etc.)
   - Writes initial config and HEAD files

### Authentication Flow

**JWT Token Generation** (`internal/auth/auth.go`):
```go
func GenerateToken(userID uint, username string) (string, error) {
    claims := jwt.MapClaims{
        "user_id":  userID,
        "username": username,
        "exp":      time.Now().Add(time.Hour * 24).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(getJWTSecret()))
}
```

**Middleware Chain** (`internal/api/middleware.go:AuthMiddleware`):
```
HTTP Request
    │
    ▼
AuthMiddleware checks Authorization header
    │
    ├─ Missing/Invalid → 401 Unauthorized
    │
    ▼
Parse and validate JWT token
    │
    ├─ Invalid/Expired → 401 Unauthorized
    │
    ▼
Extract user_id from claims
    │
    ▼
Set gin.Context["user_id"] = userID
    │
    ▼
Call next handler
```

### Git Protocol Handlers

**HTTP Git Protocol** (`internal/api/git_handlers.go`):

Three endpoints implement the Git smart HTTP protocol:

1. **GET /:owner/:repo/info/refs?service=git-upload-pack**
   - Advertises available refs (branches, tags)
   - Returns pkt-line formatted response
   - Used by `git clone` and `git fetch`

2. **POST /:owner/:repo/git-receive-pack**
   - Receives pushed objects and refs
   - Updates server-side refs
   - Used by `git push`

3. **POST /:owner/:repo/git-upload-pack**
   - Sends requested objects to client
   - Implements pack negotiation
   - Used by `git fetch` and `git clone`

**Pattern**: Each handler calls corresponding C++ function via CGo bridge:
```go
func GitInfoRefs(c *gin.Context) {
    owner := c.Param("owner")
    repo := c.Param("repo")
    service := c.Query("service")

    // Call C++ via CGo
    refs := gitcore.GetRefs(repoPath)

    // Format as Git pkt-line protocol
    c.Header("Content-Type", fmt.Sprintf("application/x-%s-advertisement", service))
    c.String(200, formatPktLine(refs))
}
```

## Configuration Structure

**File**: `configs/server.yaml`

```yaml
server:
  host: 0.0.0.0        # Bind address
  port: 8080           # HTTP port
  mode: release        # Gin mode: debug/release/test

database:
  type: sqlite3       # Database type: sqlite3, postgres, sqlserver
  path: ./data/gitserver.db  # For SQLite

  # PostgreSQL (uncomment to use)
  # type: postgres
  # host: localhost
  # port: 5432
  # name: zixiao_git
  # user: postgres
  # password: postgres
  # sslmode: disable  # Options: disable, require, verify-ca, verify-full

  # SQL Server (uncomment to use)
  # type: sqlserver
  # host: localhost
  # port: 1433
  # name: zixiao_git
  # user: sa
  # password: YourPassword123

git:
  repo_path: ./data/repositories  # Where bare repos are stored
  max_repo_size: 1024             # MB, per-repository limit

security:
  jwt_secret: CHANGE_ME_IN_PRODUCTION_USE_RANDOM_STRING
  jwt_expiration: 24              # Hours
  password_min: 8                 # Minimum password length
```

**Environment Override**: JWT secret can be set via `JWT_SECRET` environment variable (takes precedence over config file).

## Database Schema

**File**: `internal/database/database.go`

ZiXiao Git Server supports multiple database backends:
- **SQLite3**: Default, file-based, suitable for development and small deployments
- **PostgreSQL**: Recommended for production, better concurrency and performance
- **SQL Server**: Enterprise-grade, suitable for Windows environments

**Database Configuration**: See `database.Config` struct with fields:
- `Type`: Database type (sqlite3, postgres, sqlserver)
- `Path`: SQLite database file path
- `Host`, `Port`, `Name`, `User`, `Password`: Connection parameters for PostgreSQL/SQL Server
- `SSLMode`: SSL mode for PostgreSQL (disable, require, verify-ca, verify-full)

**Schema Management**: Each database type has its own schema function:
- `getSQLiteSchema()`: Uses AUTOINCREMENT, INTEGER, TEXT, DATETIME
- `getPostgreSQLSchema()`: Uses SERIAL, VARCHAR, TIMESTAMP, BOOLEAN
- `getSQLServerSchema()`: Uses IDENTITY, NVARCHAR, BIT, DATETIME

Key tables:
- `users`: id, username, password_hash (bcrypt), email, created_at, updated_at
- `repositories`: id, name, description, owner_id, is_public, created_at, updated_at
- `collaborations`: id, repository_id, user_id, permission (read/write/admin)
- `ssh_keys`: id, user_id, title, key, fingerprint, created_at
- `access_tokens`: id, user_id, token, name, last_used_at, created_at, expires_at
- `activities`: id, user_id, repository_id, action, created_at

**Important**:
- Foreign keys are enforced with `ON DELETE CASCADE` to maintain referential integrity
- Connection pooling is configured with `SetMaxOpenConns(25)` and `SetMaxIdleConns(5)`
- For detailed database configuration, migration, and optimization, see `docs/DATABASE.md`
- For Docker deployment with different databases, see `docs/DOCKER_DEPLOYMENT.md`

## Platform-Specific Build Notes

### macOS
**Backend Dependencies**:
- Dependencies via Homebrew: `brew install cmake openssl zlib go`
- Uses Clang by default
- IDE: Xcode with CMake generation (`cmake -G Xcode`)

**Frontend Dependencies**:
- Node.js: `brew install node`
- Nginx: `brew install nginx`
- Nginx config location: `/opt/homebrew/etc/nginx/servers/`

**OpenSSL Paths** (for CGo linking):
- Homebrew (Apple Silicon): `/opt/homebrew/opt/openssl/lib`
- Homebrew (Intel): `/usr/local/opt/openssl/lib`

### Linux
**Backend Dependencies**:
- Ubuntu/Debian: `apt-get install build-essential cmake libssl-dev zlib1g-dev golang`
- CentOS/RHEL: `yum install gcc gcc-c++ cmake openssl-devel zlib-devel golang`
- Uses GCC by default

**Frontend Dependencies**:
- Node.js: `apt-get install nodejs npm` (Ubuntu) or `yum install nodejs npm` (CentOS)
- Nginx: `apt-get install nginx` (Ubuntu) or `yum install nginx` (CentOS)
- Nginx config location: `/etc/nginx/sites-available/` and `/etc/nginx/sites-enabled/`

### Windows
**Backend Dependencies**:
- Requires Visual Studio 2022 with C++ workload
- Dependencies via vcpkg: `vcpkg install openssl zlib`
- Must set `VCPKG_ROOT` environment variable
- Go: Download from https://golang.org/dl/
- PowerShell scripts: `./scripts/build.ps1`, `./scripts/install.ps1`
- Can also use WSL2 for Linux-style build

**Frontend Dependencies**:
- Node.js: Download from https://nodejs.org/
- Nginx: Download from http://nginx.org/en/download.html
  - Extract to `C:\nginx\`
  - Configure with PowerShell scripts

## API Structure

**Base URL**: `http://localhost:8080/api/v1`

**Authentication Endpoints**:
- `POST /auth/register` - Register new user
- `POST /auth/login` - Login and get JWT token

**Repository Endpoints** (all require JWT auth):
- `GET /repos` - List user's repositories
- `POST /repos` - Create new repository
- `GET /repos/:id` - Get repository details
- `PUT /repos/:id` - Update repository
- `DELETE /repos/:id` - Delete repository
- `GET /repos/:id/collaborators` - List collaborators
- `POST /repos/:id/collaborators` - Add collaborator

**Git Protocol Endpoints** (no auth for public repos):
- `GET /:owner/:repo/info/refs` - Ref advertisement
- `POST /:owner/:repo/git-receive-pack` - Receive pack (push)
- `POST /:owner/:repo/git-upload-pack` - Upload pack (fetch/clone)

**Full API docs**: See `docs/API.md`

## Frontend Architecture

**Location**: `frontend/src/`

Built with **Vue 3 + Vite + MDUI 2.x**:

### Project Structure
```
frontend/
├── src/
│   ├── App.vue              # Root component with navigation
│   ├── main.js              # Entry point, Vue app initialization
│   ├── router/
│   │   └── index.js         # Vue Router config with auth guards
│   ├── stores/
│   │   ├── auth.js          # Pinia store for authentication state
│   │   └── repository.js    # Pinia store for repository data
│   ├── services/
│   │   └── index.js         # API service layer (authService, repoService)
│   ├── utils/
│   │   └── api.js           # Axios instance with JWT interceptors
│   └── views/
│       ├── Home.vue         # Landing page
│       ├── Login.vue        # Login form
│       ├── Register.vue     # Registration form
│       ├── Dashboard.vue    # User dashboard
│       ├── Repositories.vue # Repository list
│       ├── RepositoryDetail.vue  # Repository details
│       └── NotFound.vue     # 404 page
├── vite.config.js           # Vite config (proxy, build options)
└── package.json             # Dependencies
```

### Key Patterns

**MDUI 2.x Named Imports** (CRITICAL):
```javascript
// CORRECT - Use named imports
import { setTheme, snackbar } from 'mdui'

// Import components in main.js
import 'mdui/components/button.js'
import 'mdui/components/card.js'

// WRONG - MDUI 2.x doesn't export default
import mdui from 'mdui'  // Will cause build errors
```

**Pinia Store Pattern** (`frontend/src/stores/auth.js`):
```javascript
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))
  const isAuthenticated = computed(() => !!token.value)

  async function login(username, password) {
    const response = await authService.login(username, password)
    token.value = response.token
    user.value = response.user
    localStorage.setItem('token', response.token)
    localStorage.setItem('user', JSON.stringify(response.user))
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  return { token, user, isAuthenticated, login, logout }
})
```

**API Service Layer** (`frontend/src/services/index.js`):
```javascript
import apiClient from '@/utils/api'

export const authService = {
  async login(username, password) {
    return await apiClient.post('/auth/login', { username, password })
  },
  async register(username, password, email) {
    return await apiClient.post('/auth/register', { username, password, email })
  }
}

export const repositoryService = {
  async getRepositories() {
    return await apiClient.get('/repos')
  },
  async createRepository(data) {
    return await apiClient.post('/repos', data)
  }
  // ... more methods
}
```

**Axios Interceptors** (`frontend/src/utils/api.js`):
```javascript
// Request interceptor - attach JWT token
apiClient.interceptors.request.use((config) => {
  const authStore = useAuthStore()
  if (authStore.token) {
    config.headers.Authorization = `Bearer ${authStore.token}`
  }
  return config
})

// Response interceptor - handle 401
apiClient.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      const authStore = useAuthStore()
      authStore.logout()
      window.location.href = '/login'
    }
    return Promise.reject(error.response.data)
  }
)
```

**Vue Router Navigation Guards** (`frontend/src/router/index.js`):
```javascript
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next({ name: 'login', query: { redirect: to.fullPath } })
  } else {
    next()
  }
})
```

### Development Workflow

**Local Development** (frontend only):
```bash
cd frontend
npm run dev    # Vite dev server on port 3000
               # API requests proxy to localhost:8080
```

**Full Stack Development**:
```bash
# Terminal 1: Backend
make run       # Go server on port 8080

# Terminal 2: Frontend
cd frontend
npm run dev    # Vue app on port 3000 with proxy

# Terminal 3: Nginx (optional, for production-like setup)
nginx -c $(pwd)/configs/nginx-local-generated.conf
```

**Production Build**:
```bash
cd frontend
npm run build              # Builds to web/dist/
cd ..
nginx -s reload            # Reload Nginx to serve new build
```

### Nginx Configuration

**Development** (`configs/nginx-local.conf`):
- Serves Vue SPA from `web/dist/`
- Proxies `/api/*` to Go backend (port 8080)
- Proxies Git protocol endpoints (`/:owner/:repo.git/*`)
- Uses project-relative paths (replaced by install script)

**Production** (`configs/nginx-prod.conf`):
- SSL/TLS support (configure SSL certificates)
- Security headers (HSTS, X-Frame-Options)
- Gzip compression
- Static asset caching

**Key Nginx Patterns**:
```nginx
# SPA routing - all requests fall back to index.html
location / {
    root /path/to/web/dist;
    try_files $uri $uri/ /index.html;
}

# API proxy
location /api/ {
    proxy_pass http://localhost:8080;
    proxy_set_header Authorization $http_authorization;
}

# Git HTTP protocol
location ~ ^/[^/]+/[^/]+\.git/ {
    proxy_pass http://localhost:8080;
    client_max_body_size 0;        # Allow large git pushes
    proxy_buffering off;
}
```

## Testing

**Location**: `scripts/api-test.sh` (Bash) and `scripts/api-test.ps1` (PowerShell)

Both scripts test the complete API flow:
1. Register user
2. Login and get JWT token
3. Create repository
4. List repositories
5. Get repository details
6. Add collaborator
7. Delete repository

**Usage**: Start server first (`make run`), then run test script.

## Important Patterns and Conventions

### Error Handling in API Handlers

```go
func SomeHandler(c *gin.Context) {
    // Get authenticated user ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(401, gin.H{"error": "unauthorized"})
        return
    }

    // Validate input
    var req SomeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // Business logic
    result, err := someBusinessLogic(userID.(uint), req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, result)
}
```

### CGo Memory Management

**Critical**: Always free C-allocated memory in Go:
```go
cPath := C.CString(path)
defer C.free(unsafe.Pointer(cPath))  // Essential to prevent memory leaks
```

### Repository Path Construction

**Convention**: `{repo_path}/{owner}/{repo}.git`

Example: `./data/repositories/alice/myproject.git`

This matches GitHub/GitLab structure and allows easy file-based repository storage.

## Key Files to Understand First

When working on this codebase, start by reading these files in order:

**Full Stack Overview**:
1. **README.md** - Overall project structure and features
2. **configs/server.yaml** - Configuration options
3. **CLAUDE.md** (this file) - Architecture deep dive

**Backend Architecture**:
4. **internal/database/database.go** - Database schema
5. **pkg/gitcore/gitcore.go** - CGo bridge interface
6. **internal/api/routes.go** - API endpoint structure
7. **internal/auth/auth.go** - Authentication patterns
8. **git-core/include/git_c_api.h** - C API surface

**Frontend Architecture**:
9. **frontend/src/main.js** - Vue app initialization
10. **frontend/src/router/index.js** - Routing and navigation guards
11. **frontend/src/stores/auth.js** - Authentication state management
12. **frontend/src/utils/api.js** - Axios setup and interceptors
13. **frontend/vite.config.js** - Build configuration

**Deployment**:
14. **configs/nginx-local.conf** - Nginx reverse proxy setup

## Common Development Tasks

### Adding a New API Endpoint

1. Define route in `internal/api/routes.go`
2. Add handler function in appropriate `*_handlers.go` file
3. Update authentication middleware if needed
4. Update `docs/API.md` with endpoint documentation
5. Add test case to `scripts/api-test.sh`
6. Create corresponding frontend service method in `frontend/src/services/index.js`

### Adding a New Vue Page

1. Create component in `frontend/src/views/YourPage.vue`:
   ```vue
   <script setup>
   import { ref, onMounted } from 'vue'
   import { snackbar } from 'mdui'
   import { useAuthStore } from '@/stores/auth'

   const authStore = useAuthStore()

   onMounted(() => {
     // Fetch data
   })
   </script>

   <template>
     <div class="container">
       <!-- MDUI components -->
     </div>
   </template>
   ```

2. Add route in `frontend/src/router/index.js`:
   ```javascript
   {
     path: '/your-path',
     name: 'yourPage',
     component: () => import('@/views/YourPage.vue'),
     meta: { requiresAuth: true }  // If authentication needed
   }
   ```

3. Add navigation link in `frontend/src/App.vue` if needed

### Adding a New Pinia Store

1. Create store file `frontend/src/stores/yourStore.js`:
   ```javascript
   import { defineStore } from 'pinia'
   import { ref } from 'vue'

   export const useYourStore = defineStore('yourStore', () => {
     const data = ref([])

     async function fetchData() {
       // API call
     }

     return { data, fetchData }
   })
   ```

2. Use in component:
   ```javascript
   import { useYourStore } from '@/stores/yourStore'
   const yourStore = useYourStore()
   ```

### Adding a New Git Operation

1. Add C++ method to appropriate header in `git-core/include/`
2. Implement in corresponding `.cpp` file in `git-core/src/`
3. Add C wrapper function in `git-core/include/git_c_api.h`
4. Implement C wrapper in `git-core/src/git_c_api.cpp`
5. Add Go wrapper in `pkg/gitcore/gitcore.go`
6. Rebuild C++ library: `make build-cpp`

### Modifying Database Schema

1. Update schema SQL in appropriate function in `internal/database/database.go`:
   - `getSQLiteSchema()` for SQLite3
   - `getPostgreSQLSchema()` for PostgreSQL
   - `getSQLServerSchema()` for SQL Server
2. Update corresponding model struct in `internal/models/models.go`
3. Consider migration strategy for existing databases (currently requires manual handling)
4. Test with fresh database:
   - SQLite: `rm data/gitserver.db && make run`
   - PostgreSQL: Drop and recreate database
   - SQL Server: Drop and recreate database

### Updating Nginx Configuration

1. Edit config files: `configs/nginx-{local,dev,prod}.conf`
2. For local development:
   ```bash
   ./scripts/install.sh    # Regenerates nginx-local-generated.conf
   nginx -s reload         # Reload Nginx
   ```
3. Test configuration: `nginx -t -c /path/to/config`

## Debugging Tips

### Frontend Debugging

**Vue DevTools**:
- Install Vue DevTools browser extension
- Inspect component hierarchy, Pinia stores, router
- Monitor events and performance

**Vite Dev Server**:
```bash
cd frontend
npm run dev    # Hot module replacement enabled
               # View at http://localhost:3000
```

**Common Issues**:
- **MDUI import errors**: Always use named imports, never default import
  ```javascript
  // ✓ CORRECT
  import { snackbar, setTheme } from 'mdui'

  // ✗ WRONG - will cause "default is not exported" error
  import mdui from 'mdui'
  ```
- **401 errors**: Check JWT token in localStorage, verify Authorization header
- **CORS errors**: Ensure Vite proxy is configured or Nginx is running
- **Build errors**: Run `npm run lint` to catch Vue/JS errors

**Browser DevTools**:
- Network tab: Inspect API requests, check headers
- Console: View errors and logs
- Application tab: Check localStorage for token and user data

### CGo Debugging

- Enable verbose CGo: `CGO_ENABLED=1 CGO_LDFLAGS="-Wl,-v" go build -x ./cmd/server`
- Check C++ library is built: `ls -l git-core/lib/libgitcore.*`
- Verify library linking: `otool -L bin/zixiao-git-server` (macOS) or `ldd bin/zixiao-git-server` (Linux)

### API Debugging

- Run in debug mode: Edit `configs/server.yaml` to set `server.mode: debug`
- Enable Gin logging: Already enabled by `gin.Default()`
- Check JWT token: Use https://jwt.io to decode token payload
- Test with curl: `curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/repos`

### Git Protocol Debugging

- Enable Git client tracing: `GIT_TRACE=1 GIT_CURL_VERBOSE=1 git clone http://...`
- Check server logs for protocol-level errors
- Verify repository path: `ls -la data/repositories/{owner}/{repo}.git`

## Security Considerations

- **JWT Secret**: Must be changed in production (use random 64-character string)
- **Password Hashing**: Uses bcrypt with cost 10 (defined in `internal/auth/auth.go`)
- **CORS**: Configured in `internal/api/middleware.go:CORSMiddleware` - adjust for production
- **SQL Injection**: Protected by parameterized queries via database/sql
- **Repository Access**: Always check ownership or collaboration before operations

## Known Limitations and TODOs

- **No database migrations system** (manual schema updates required) - important for production
- **SSH protocol** not yet implemented (only HTTP Git protocol)
- **No repository forking** feature
- **No pull request / merge request** system
- **No CI/CD** integration
- **No webhooks** support
- **No Git LFS** support
- **No repository size enforcement** (max_repo_size in config not enforced)

## Additional Resources

- **Database Configuration**: `docs/DATABASE.md` - Detailed guide for SQLite3, PostgreSQL, SQL Server
- **Docker Deployment**: `docs/DOCKER_DEPLOYMENT.md` - Complete Docker deployment guide
- **API Documentation**: `docs/API.md`
- **Frontend Development Guide**: `docs/FRONTEND_DEV.md`
- **Frontend Deployment Guide**: `docs/FRONTEND_DEPLOYMENT.md`
- **Frontend Migration Summary**: `FRONTEND_MIGRATION.md`
- **Windows Setup**: `docs/WINDOWS.md`
- **VS Code Setup**: `docs/VSCODE.md`
- **Xcode Setup**: `docs/XCODE.md`
- **Project Summary**: `PROJECT_SUMMARY.md`
- **Quick Start**: `QUICKSTART.md`

---

**Last Updated**: 2025-10-16
**Project Version**: 1.0.0 (Vue 3 Frontend)
**License**: MIT
