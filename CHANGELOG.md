# Changelog

All notable changes to ZiXiao Git Server will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2025-10-16

### Added

#### Frontend
- Vue 3 SPA with modern UI/UX
- MDUI 2.x Material Design components
- Pinia state management
- Vue Router with navigation guards
- Dark/light/auto theme support
- Responsive design for mobile and desktop
- Real-time JWT token management
- Axios interceptors for API requests

#### Backend
- Complete REST API for repository management
- JWT authentication with token expiration
- User registration and login
- Repository CRUD operations
- Collaborator management
- Public/private repository access control
- SQLite database with migration support

#### Git Core
- HTTP Git protocol implementation (smart HTTP)
- Git info/refs advertisement
- Git receive-pack (push)
- Git upload-pack (fetch/clone)
- C++ Git core library with CGo bindings
- Repository initialization and management
- Reference (branches/tags) management

#### Infrastructure
- GitHub Actions CI/CD workflows
  - Backend CI (Go + C++ testing)
  - Frontend CI (Vue linting and building)
  - Multi-platform release builds
  - Docker image building and publishing
- Docker support
  - Multi-stage Dockerfile
  - Docker Compose configuration
  - PostgreSQL and Redis profiles
- Nginx reverse proxy configuration
  - Local, development, and production configs
  - SPA routing support
  - Git protocol proxying

#### Documentation
- Comprehensive README with examples
- CLAUDE.md for architecture deep dive
- Frontend development guide
- Frontend deployment guide
- GitHub Actions documentation
- API documentation
- Installation scripts for all platforms

### Changed
- Migrated from static HTML to Vue 3 SPA
- Updated project structure for full-stack architecture
- Improved error handling and validation
- Enhanced security with proper CORS and JWT configuration

### Fixed
- CGo linking issues on macOS (OpenSSL paths)
- MDUI 2.x import errors (named imports)
- Go module dependencies (go.sum)
- Nginx configuration for SPA routing

## [0.1.0] - 2024-10-16

### Added
- Initial project setup
- Basic Go HTTP server with Gin framework
- C++ Git core library
- CGo bridge layer
- SQLite database support
- Basic authentication
- Repository creation and management
- Static HTML frontend

[Unreleased]: https://github.com/Zixiao-System/ZiXiao-Git-Server/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/Zixiao-System/ZiXiao-Git-Server/compare/v0.1.0...v1.0.0
[0.1.0]: https://github.com/Zixiao-System/ZiXiao-Git-Server/releases/tag/v0.1.0
