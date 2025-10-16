# ZiXiao Git Server - Vue Frontend Integration Summary

## What's New

The frontend has been completely refactored from a static MDUI page to a dynamic Vue 3 Single Page Application (SPA) with Nginx as the default web server.

## Architecture

```
┌──────────────────────────────────────┐
│   Nginx (Port 80) - Web Server      │
│   - Serves Vue.js SPA                │
│   - Proxies /api/* to Go backend     │
│   - Handles Git protocol             │
└──────────────┬───────────────────────┘
               │
               ▼
┌──────────────────────────────────────┐
│   Go Backend (Port 8080) - API       │
│   - REST API endpoints               │
│   - Git HTTP protocol                │
│   - Authentication & Business logic  │
└──────────────────────────────────────┘
```

## Key Features

### Frontend (Vue 3)
- **Modern SPA**: Vue 3 with Composition API
- **Routing**: Vue Router with navigation guards
- **State Management**: Pinia for auth and repository data
- **API Layer**: Axios with JWT auto-attachment
- **UI Framework**: MDUI 2.x Material Design components
- **Build Tool**: Vite for fast development and optimized builds

### Pages Implemented
1. **Home** (`/`) - Landing page with features
2. **Login** (`/login`) - User authentication
3. **Register** (`/register`) - User registration
4. **Dashboard** (`/dashboard`) - User overview and stats
5. **Repositories** (`/repositories`) - Repository list with search
6. **Repository Detail** (`/repositories/:id`) - Single repository view
7. **404** - Not found page

### Nginx Integration
- **Development**: Local file serving with hot reload support
- **Production**: Optimized static file serving with caching
- **Auto-installation**: Scripts automatically install and configure Nginx
- **Multiple configs**: Dev, prod, and local configurations

## Quick Start

### 1. Install Dependencies

```bash
# Run the installation script (installs Node.js, Nginx, and all dependencies)
./scripts/install.sh  # Linux/macOS
# or
.\scripts\install.ps1  # Windows
```

### 2. Build Frontend

```bash
cd frontend
npm run build
```

Or use the build script:

```bash
./scripts/build-frontend.sh  # Linux/macOS
# or
.\scripts\build-frontend.ps1  # Windows
```

### 3. Start Backend

```bash
make build
make run
```

Backend runs on `http://localhost:8080`

### 4. Start Nginx

**macOS**:
```bash
brew services start nginx
# or restart if already running
brew services restart nginx
```

**Linux**:
```bash
sudo systemctl start nginx
# or restart
sudo systemctl restart nginx
```

**Windows**:
```bash
cd C:\nginx
start nginx
```

### 5. Access Application

Visit `http://localhost` in your browser.

## Development Workflow

### Option 1: Frontend Dev Server (Recommended for development)

```bash
# Terminal 1: Start backend
make run

# Terminal 2: Start frontend dev server
cd frontend
npm run dev
```

Access at `http://localhost:3000` (with hot reload)

### Option 2: Full Stack with Nginx

```bash
# Build frontend
cd frontend
npm run build

# Start backend
make run

# Start Nginx
brew services start nginx  # macOS
# or
sudo systemctl start nginx  # Linux
```

Access at `http://localhost`

## Project Structure

```
ZiXiao-Git-Server/
├── frontend/                  # Vue 3 frontend
│   ├── src/
│   │   ├── views/            # Page components
│   │   ├── stores/           # Pinia stores
│   │   ├── services/         # API services
│   │   ├── router/           # Vue Router
│   │   └── utils/            # Utilities
│   ├── package.json
│   └── vite.config.js
├── web/
│   └── dist/                 # Built frontend (generated)
├── configs/
│   ├── nginx-dev.conf        # Development Nginx config
│   ├── nginx-prod.conf       # Production Nginx config
│   ├── nginx-local.conf      # Local development template
│   └── nginx-local-generated.conf  # Generated local config
├── scripts/
│   ├── install.sh            # Updated with Nginx & frontend
│   ├── install.ps1           # Updated with Nginx & frontend
│   ├── build-frontend.sh     # Frontend build script
│   └── build-frontend.ps1    # Frontend build script (Windows)
└── docs/
    ├── FRONTEND_DEV.md       # Development guide
    └── FRONTEND_DEPLOYMENT.md # Deployment guide
```

## Configuration Files

### Nginx Configurations

1. **nginx-local.conf** - Template for local development
2. **nginx-dev.conf** - Development server configuration
3. **nginx-prod.conf** - Production with SSL support

### Environment Variables

**Frontend (`frontend/.env*`)**:
- `.env` - Default config
- `.env.development` - Dev mode (API: `http://localhost:8080`)
- `.env.production` - Production (API: `/api/v1`)

## Installation Script Updates

### Linux/macOS (`scripts/install.sh`)

**New Features**:
- Checks for and installs Node.js
- Checks for and installs Nginx
- Installs frontend dependencies (`npm install`)
- Configures Nginx automatically
- Creates Nginx configuration symlinks

### Windows (`scripts/install.ps1`)

**New Features**:
- Checks for and installs Node.js
- Downloads and installs Nginx for Windows
- Installs frontend dependencies
- Configures Nginx automatically
- Updates Nginx main config

## API Integration

### Axios Configuration

```javascript
// Automatic JWT token attachment
// Base URL from environment variables
// 401 auto-redirect to login
// Error handling with snackbar notifications
```

### Service Layer

```javascript
// Authentication
authService.login(username, password)
authService.register(username, password, email)

// Repositories
repositoryService.getRepositories()
repositoryService.createRepository(data)
repositoryService.updateRepository(id, data)
repositoryService.deleteRepository(id)
```

## State Management

### Auth Store
- User authentication state
- JWT token management
- Login/logout functionality

### Repository Store
- Repository CRUD operations
- Loading states
- Error handling

## Routing & Navigation

- **Public routes**: Home, Login, Register
- **Protected routes**: Dashboard, Repositories (require authentication)
- **Auto-redirect**: Unauthenticated → Login, Authenticated → Dashboard
- **SPA routing**: All handled by Vue Router with Nginx fallback

## Build & Deployment

### Development Build

```bash
cd frontend
npm run build
```

Output: `web/dist/`

### Production Deployment

See `docs/FRONTEND_DEPLOYMENT.md` for:
- Docker deployment
- Traditional server deployment
- SSL configuration
- Performance optimization
- Monitoring and troubleshooting

## Documentation

- **FRONTEND_DEV.md**: Complete development guide
- **FRONTEND_DEPLOYMENT.md**: Deployment and production guide
- **frontend/README.md**: Quick reference

## Testing

The application can be tested by:

1. **Register**: Create a new user account
2. **Login**: Authenticate with credentials
3. **Dashboard**: View user overview
4. **Create Repository**: Add new Git repository
5. **Repository List**: Browse and search repositories
6. **Repository Details**: View repository info and clone URL

## Nginx Endpoints

### Frontend
- `/` - Vue SPA (all routes handled by Vue Router)

### Backend Proxies
- `/api/*` - REST API endpoints
- `/:owner/:repo.git/*` - Git HTTP protocol

## Known Issues & Limitations

1. **Hot Reload**: Use `npm run dev` for hot reload during development
2. **Windows Nginx**: May require manual download if automatic fails
3. **SSL**: Production SSL needs manual certificate setup

## Next Steps

1. **Run installation script**: `./scripts/install.sh`
2. **Build frontend**: `cd frontend && npm run build`
3. **Start backend**: `make run`
4. **Start Nginx**: `brew services start nginx` (or equivalent)
5. **Visit**: `http://localhost`

## Support

- GitHub Issues: https://github.com/Zixiao-System/ZiXiao-Git-Server/issues
- Documentation: See `docs/` directory
- Frontend README: See `frontend/README.md`

---

**Note**: The old static `web/index.html` is replaced by the Vue SPA build output in `web/dist/`. Nginx now serves the Vue application with proper SPA routing support.
