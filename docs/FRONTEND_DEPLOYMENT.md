# ZiXiao Git Server - Frontend Deployment Guide

This guide covers deploying the Vue.js frontend with Nginx as the web server.

## Table of Contents

- [Architecture Overview](#architecture-overview)
- [Development Setup](#development-setup)
- [Production Deployment](#production-deployment)
- [Nginx Configuration](#nginx-configuration)
- [Troubleshooting](#troubleshooting)

## Architecture Overview

The ZiXiao Git Server uses a decoupled architecture:

```
┌─────────────────────────────────────┐
│   Nginx (Port 80/443)               │  ← Frontend Web Server
│   - Serves static Vue.js files      │
│   - Proxies /api/* to backend       │
│   - Handles Git protocol endpoints  │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│   Go Backend (Port 8080)            │  ← API Server
│   - REST API                        │
│   - Git HTTP protocol               │
│   - Authentication & Business Logic │
└─────────────────────────────────────┘
```

## Development Setup

### Prerequisites

- Node.js 18+ and npm
- Nginx
- Go backend server running on port 8080

### Quick Start

1. **Install dependencies**:
   ```bash
   cd frontend
   npm install
   ```

2. **Development mode** (with hot reload):
   ```bash
   cd frontend
   npm run dev
   ```
   Frontend will be available at `http://localhost:3000`
   API requests are proxied to `http://localhost:8080`

3. **Build for development**:
   ```bash
   ./scripts/build-frontend.sh  # Linux/macOS
   # or
   .\scripts\build-frontend.ps1  # Windows
   ```

4. **Run with Nginx**:
   ```bash
   # macOS
   brew services start nginx

   # Linux
   sudo systemctl start nginx

   # Windows
   cd C:\nginx
   start nginx
   ```

Visit `http://localhost` to access the application.

## Production Deployment

### Step 1: Build Frontend

```bash
cd frontend
npm run build
```

This creates optimized production files in `web/dist/`.

### Step 2: Deploy to Server

#### Option A: Docker Deployment (Recommended)

Create a `Dockerfile`:

```dockerfile
# Build stage
FROM node:18-alpine AS build
WORKDIR /app
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Production stage
FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
COPY configs/nginx-prod.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

Build and run:

```bash
docker build -t zixiao-git-server-frontend .
docker run -d -p 80:80 --name zixiao-frontend zixiao-git-server-frontend
```

#### Option B: Traditional Server Deployment

1. **Copy built files to server**:
   ```bash
   scp -r web/dist/* user@your-server:/usr/share/nginx/html/
   ```

2. **Copy Nginx configuration**:
   ```bash
   scp configs/nginx-prod.conf user@your-server:/etc/nginx/sites-available/zixiao-git-server.conf
   ```

3. **Enable site** (Linux):
   ```bash
   sudo ln -s /etc/nginx/sites-available/zixiao-git-server.conf /etc/nginx/sites-enabled/
   sudo nginx -t
   sudo systemctl reload nginx
   ```

### Step 3: Configure SSL (Production)

1. **Obtain SSL certificate** (using Let's Encrypt):
   ```bash
   sudo apt-get install certbot python3-certbot-nginx
   sudo certbot --nginx -d git.yourdomain.com
   ```

2. **Update Nginx configuration**:
   Edit `configs/nginx-prod.conf` and uncomment SSL lines:
   ```nginx
   ssl_certificate /etc/letsencrypt/live/git.yourdomain.com/fullchain.pem;
   ssl_certificate_key /etc/letsencrypt/live/git.yourdomain.com/privkey.pem;
   ```

3. **Set up auto-renewal**:
   ```bash
   sudo certbot renew --dry-run
   ```

## Nginx Configuration

### Development Configuration (`nginx-local.conf`)

- Serves files from local project directory
- Suitable for local development
- No SSL

### Production Configuration (`nginx-prod.conf`)

- Serves from `/usr/share/nginx/html`
- HTTPS with SSL support
- Security headers
- Gzip compression
- Rate limiting (optional)

### Key Nginx Directives

```nginx
# Frontend routing - SPA support
location / {
    root /usr/share/nginx/html;
    try_files $uri $uri/ /index.html;  # Fallback to index.html for SPA routing
    index index.html;
}

# Backend API proxy
location /api/ {
    proxy_pass http://localhost:8080;
    # Headers for proper proxying
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}

# Git HTTP protocol
location ~ ^/[^/]+/[^/]+\.git/ {
    proxy_pass http://localhost:8080;
    client_max_body_size 0;  # Allow large git pushes
    proxy_buffering off;
}
```

## Environment Variables

Create `.env.production` in `frontend/`:

```env
VITE_API_BASE_URL=/api/v1
```

## Troubleshooting

### Issue: 404 on page refresh

**Cause**: Nginx not configured for SPA routing

**Solution**: Ensure `try_files $uri $uri/ /index.html;` is in your Nginx config

### Issue: API requests failing

**Cause**: Backend not running or CORS issues

**Solution**:
1. Verify backend is running: `curl http://localhost:8080/api/v1/health`
2. Check Nginx error logs: `tail -f /var/log/nginx/error.log`
3. Verify proxy_pass configuration

### Issue: Static assets not loading

**Cause**: Incorrect base URL or paths

**Solution**:
1. Check `vite.config.js` base setting
2. Verify Nginx root directive points to correct directory
3. Check file permissions: `ls -la /usr/share/nginx/html`

### Issue: Git clone/push fails

**Cause**: Nginx not proxying Git protocol endpoints

**Solution**:
1. Verify Git location block in Nginx config
2. Check `client_max_body_size 0;` for large repositories
3. Ensure `proxy_buffering off;` is set

## Performance Optimization

### 1. Enable Gzip Compression

```nginx
gzip on;
gzip_vary on;
gzip_types text/css application/javascript image/svg+xml;
gzip_comp_level 6;
```

### 2. Browser Caching

```nginx
location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
    expires 1y;
    add_header Cache-Control "public, immutable";
}
```

### 3. HTTP/2

```nginx
listen 443 ssl http2;
```

## Monitoring

### Check Nginx Status

```bash
# Linux
sudo systemctl status nginx

# macOS
brew services list | grep nginx

# Windows
tasklist | findstr nginx
```

### View Logs

```bash
# Access logs
tail -f /var/log/nginx/zixiao-git-server-access.log

# Error logs
tail -f /var/log/nginx/zixiao-git-server-error.log
```

### Test Configuration

```bash
# Linux/macOS
sudo nginx -t

# Windows
C:\nginx\nginx.exe -t
```

## Backup and Rollback

### Backup

```bash
# Backup current frontend
tar -czf frontend-backup-$(date +%Y%m%d).tar.gz /usr/share/nginx/html/

# Backup Nginx config
cp /etc/nginx/sites-available/zixiao-git-server.conf{,.backup}
```

### Rollback

```bash
# Restore frontend
tar -xzf frontend-backup-20240101.tar.gz -C /

# Restore Nginx config
cp /etc/nginx/sites-available/zixiao-git-server.conf{.backup,}
sudo systemctl reload nginx
```

## Security Checklist

- [ ] SSL/TLS configured with valid certificate
- [ ] HTTP redirects to HTTPS
- [ ] Security headers configured (X-Frame-Options, CSP, etc.)
- [ ] Rate limiting enabled for API endpoints
- [ ] File upload size limits configured
- [ ] Nginx running as non-root user
- [ ] Regular security updates applied
- [ ] Access logs monitored for suspicious activity

## Additional Resources

- [Vue.js Deployment Guide](https://vuejs.org/guide/best-practices/production-deployment.html)
- [Nginx Documentation](https://nginx.org/en/docs/)
- [Vite Production Build](https://vitejs.dev/guide/build.html)

## Support

For issues and questions:
- GitHub Issues: https://github.com/Zixiao-System/ZiXiao-Git-Server/issues
- Documentation: See `docs/` directory
