# ZiXiao Git Server API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication

Most endpoints require authentication using JWT tokens. Include the token in the Authorization header:

```
Authorization: Bearer <token>
```

## Endpoints

### Authentication

#### Register a new user
```http
POST /auth/register
```

Request body:
```json
{
  "username": "alice",
  "email": "alice@example.com",
  "password": "password123",
  "full_name": "Alice Smith"
}
```

Response (201 Created):
```json
{
  "user": {
    "id": 1,
    "username": "alice",
    "email": "alice@example.com",
    "full_name": "Alice Smith",
    "is_admin": false,
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### Login
```http
POST /auth/login
```

Request body:
```json
{
  "username": "alice",
  "password": "password123"
}
```

Response (200 OK):
```json
{
  "user": {
    "id": 1,
    "username": "alice",
    "email": "alice@example.com",
    "full_name": "Alice Smith",
    "is_admin": false,
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Users

#### Get current user
```http
GET /user
Authorization: Bearer <token>
```

Response (200 OK):
```json
{
  "user": {
    "id": 1,
    "username": "alice",
    "email": "alice@example.com",
    "full_name": "Alice Smith",
    "is_admin": false,
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### Get user by username
```http
GET /users/:username
```

Response (200 OK):
```json
{
  "user": {
    "id": 1,
    "username": "alice",
    "email": "alice@example.com",
    "full_name": "Alice Smith",
    "is_admin": false,
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### Repositories

#### Create a repository
```http
POST /repos
Authorization: Bearer <token>
```

Request body:
```json
{
  "name": "my-project",
  "description": "My awesome project",
  "is_private": false
}
```

Response (201 Created):
```json
{
  "repository": {
    "id": 1,
    "name": "my-project",
    "description": "My awesome project",
    "owner_id": 1,
    "owner_name": "alice",
    "is_private": false,
    "default_branch": "main",
    "size": 0,
    "stars": 0,
    "forks": 0,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### Get repository
```http
GET /repos/:owner/:repo
```

Response (200 OK):
```json
{
  "repository": {
    "id": 1,
    "name": "my-project",
    "description": "My awesome project",
    "owner_id": 1,
    "owner_name": "alice",
    "is_private": false,
    "default_branch": "main",
    "size": 0,
    "stars": 0,
    "forks": 0,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### List user repositories
```http
GET /users/:owner/repos
```

Response (200 OK):
```json
{
  "repositories": [
    {
      "id": 1,
      "name": "my-project",
      "description": "My awesome project",
      "owner_id": 1,
      "owner_name": "alice",
      "is_private": false,
      "default_branch": "main",
      "size": 0,
      "stars": 0,
      "forks": 0,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

#### Delete repository
```http
DELETE /repos/:owner/:repo
Authorization: Bearer <token>
```

Response (200 OK):
```json
{
  "message": "repository deleted"
}
```

### Collaborators

#### Add collaborator
```http
POST /repos/:owner/:repo/collaborators
Authorization: Bearer <token>
```

Request body:
```json
{
  "username": "bob",
  "permission": "write"
}
```

Permissions: `read`, `write`, `admin`

Response (200 OK):
```json
{
  "message": "collaborator added"
}
```

#### Remove collaborator
```http
DELETE /repos/:owner/:repo/collaborators/:username
Authorization: Bearer <token>
```

Response (200 OK):
```json
{
  "message": "collaborator removed"
}
```

## Git HTTP Protocol

### Clone repository
```bash
git clone http://localhost:8080/alice/my-project.git
```

### Clone private repository
```bash
git clone http://alice:<token>@localhost:8080/alice/my-project.git
```

### Push to repository
```bash
git push http://alice:<token>@localhost:8080/alice/my-project.git main
```

## Error Responses

### 400 Bad Request
```json
{
  "error": "invalid request parameters"
}
```

### 401 Unauthorized
```json
{
  "error": "missing authorization token"
}
```

### 403 Forbidden
```json
{
  "error": "access denied"
}
```

### 404 Not Found
```json
{
  "error": "repository not found"
}
```

### 409 Conflict
```json
{
  "error": "repository already exists"
}
```

### 500 Internal Server Error
```json
{
  "error": "internal server error"
}
```
