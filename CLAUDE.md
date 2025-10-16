# CLAUDE.md - ZiXiao Git Server

This file provides context for Claude Code instances working on the ZiXiao Git Server codebase.

## Project Overview

ZiXiao Git Server is a lightweight, high-performance Git server with a hybrid architecture:
- **C++ Core**: Low-level Git operations (repository, objects, protocol, pack files)
- **Go Backend**: HTTP API, business logic, authentication, database
- **CGo Bridge**: C API wrapper enabling Go to call C++ functions

**Tech Stack**: Go 1.21+, C++17, Gin web framework, JWT auth, SQLite, MDUI 2.x frontend

## Common Commands

```bash
# Build everything (C++ library + Go server)
make build

# Build C++ library only
make build-cpp

# Build Go server only
make build-go

# Run the server (builds first if needed)
make run

# Clean build artifacts
make clean

# Install dependencies (macOS/Linux)
./scripts/install.sh

# Install dependencies (Windows PowerShell)
./scripts/install.ps1

# API testing
./scripts/api-test.sh
```

## Architecture Deep Dive

### Three-Layer Hybrid Architecture

```
┌─────────────────────────────────────┐
│   Go HTTP API Layer (Gin)           │  ← REST API, JWT Auth, CORS
│   internal/api/                      │
└─────────────────┬───────────────────┘
                  │
┌─────────────────▼───────────────────┐
│   Go Business Logic Layer           │  ← Repository CRUD, User Management
│   internal/{repository, auth, etc}  │     Database Operations (SQLite)
└─────────────────┬───────────────────┘
                  │
┌─────────────────▼───────────────────┐
│   CGo Bridge Layer                  │  ← pkg/gitcore/gitcore.go
│   #cgo directives, C.* calls        │     Type conversions (Go ↔ C)
└─────────────────┬───────────────────┘
                  │
┌─────────────────▼───────────────────┐
│   C API Wrapper (C ABI)             │  ← git-core/include/git_c_api.h
│   extern "C" functions              │     C-compatible function signatures
└─────────────────┬───────────────────┘
                  │
┌─────────────────▼───────────────────┐
│   C++ Git Core Library              │  ← git-core/src/*.cpp
│   GitRepository, GitObject, etc     │     Actual Git operations (libgit2-style)
└─────────────────────────────────────┘
```

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
  type: sqlite         # Database type (currently only SQLite)
  path: ./data/gitserver.db

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

Key tables:
- `users`: id, username, password_hash (bcrypt), email, created_at, updated_at
- `repositories`: id, name, description, owner_id, is_public, created_at, updated_at
- `collaborations`: id, repository_id, user_id, permission (read/write/admin)
- `ssh_keys`: id, user_id, title, key, fingerprint, created_at
- `access_tokens`: id, user_id, token, name, last_used_at, created_at, expires_at
- `activities`: id, user_id, repository_id, action, created_at

**Important**: Foreign keys are enforced with `ON DELETE CASCADE` to maintain referential integrity.

## Platform-Specific Build Notes

### macOS
- Dependencies via Homebrew: `brew install cmake openssl zlib`
- Uses Clang by default
- IDE: Xcode with CMake generation (`cmake -G Xcode`)

### Linux
- Ubuntu/Debian: `apt-get install build-essential cmake libssl-dev zlib1g-dev`
- CentOS/RHEL: `yum install gcc gcc-c++ cmake openssl-devel zlib-devel`
- Uses GCC by default

### Windows
- Requires Visual Studio 2022 with C++ workload
- Dependencies via vcpkg: `vcpkg install openssl zlib`
- Must set `VCPKG_ROOT` environment variable
- PowerShell scripts: `./scripts/build.ps1`, `./scripts/install.ps1`
- Can also use WSL2 for Linux-style build

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

**File**: `web/index.html`

Built with **MDUI 2.x** (Material Design UI Web Components):
- Auto/light/dark theme support with `mdui.setTheme()`
- Material Design 3 components: `<mdui-card>`, `<mdui-button>`, `<mdui-top-app-bar>`
- Theme toggle cycles: auto → light → dark → auto
- Copy-to-clipboard with snackbar feedback
- Fully responsive grid layouts

**Pattern**: Uses Web Components (custom HTML elements) instead of framework like React/Vue.

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

1. **README.md** - Overall project structure and features
2. **configs/server.yaml** - Configuration options
3. **internal/database/database.go** - Database schema
4. **pkg/gitcore/gitcore.go** - CGo bridge interface
5. **internal/api/routes.go** - API endpoint structure
6. **internal/auth/auth.go** - Authentication patterns
7. **git-core/include/git_c_api.h** - C API surface

## Common Development Tasks

### Adding a New API Endpoint

1. Define route in `internal/api/routes.go`
2. Add handler function in appropriate `*_handlers.go` file
3. Update authentication middleware if needed
4. Update `docs/API.md` with endpoint documentation
5. Add test case to `scripts/api-test.sh`

### Adding a New Git Operation

1. Add C++ method to appropriate header in `git-core/include/`
2. Implement in corresponding `.cpp` file in `git-core/src/`
3. Add C wrapper function in `git-core/include/git_c_api.h`
4. Implement C wrapper in `git-core/src/git_c_api.cpp`
5. Add Go wrapper in `pkg/gitcore/gitcore.go`
6. Rebuild C++ library: `make build-cpp`

### Modifying Database Schema

1. Update schema SQL in `internal/database/database.go:initSchema()`
2. Update corresponding model struct in `internal/models/models.go`
3. Consider migration strategy for existing databases (currently requires manual handling)
4. Test with fresh database: `rm data/gitserver.db && make run`

## Debugging Tips

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

- No database migrations system (manual schema updates required)
- SSH protocol not yet implemented (only HTTP Git protocol)
- No repository forking feature
- No pull request / merge request system
- No CI/CD integration
- No webhooks support
- No Git LFS support
- No repository size enforcement (max_repo_size in config not enforced)

## Additional Resources

- **API Documentation**: `docs/API.md`
- **Windows Setup**: `docs/WINDOWS.md`
- **VS Code Setup**: `docs/VSCODE.md`
- **Xcode Setup**: `docs/XCODE.md`
- **Project Summary**: `PROJECT_SUMMARY.md`
- **Quick Start**: `QUICKSTART.md`

---

**Last Updated**: 2024-10-16
**Project Version**: 0.1.0
**License**: MIT
