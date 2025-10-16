# Docker Deployment Guide

This guide explains how to deploy ZiXiao Git Server using Docker with different database backends.

## Quick Start (SQLite)

The simplest deployment uses SQLite (default):

```bash
# Clone the repository
git clone https://github.com/Zixiao-System/ZiXiao-Git-Server.git
cd ZiXiao-Git-Server

# Create environment file
cp .env.example .env
# Edit .env and change JWT_SECRET

# Start services
docker-compose up -d

# Check logs
docker-compose logs -f git-server
```

Access the application at:
- Frontend: http://localhost
- API: http://localhost:8080
- Adminer (if enabled): http://localhost:8081

## PostgreSQL Deployment

For production deployments, PostgreSQL is recommended:

### 1. Configure Environment

```bash
cp .env.example .env
```

Edit `.env`:
```bash
JWT_SECRET=your-super-secret-key-change-this-in-production

# Database type
DB_TYPE=postgres
DB_HOST=postgres
DB_PORT=5432
DB_NAME=zixiao_git
DB_USER=zixiao
DB_PASSWORD=your-secure-password

# PostgreSQL container settings
POSTGRES_DB=zixiao_git
POSTGRES_USER=zixiao
POSTGRES_PASSWORD=your-secure-password
```

### 2. Update server.yaml

Edit `configs/server.yaml`:
```yaml
database:
  type: postgres
  host: postgres  # Docker service name
  port: 5432
  name: zixiao_git
  user: zixiao
  password: your-secure-password
  sslmode: disable  # Use 'require' for external PostgreSQL
```

### 3. Start Services

```bash
# Uncomment depends_on in docker-compose.yml
# Under git-server service:
depends_on:
  - postgres

# Start with PostgreSQL profile
docker-compose --profile postgres up -d
```

### 4. Verify Connection

```bash
# Check logs
docker-compose logs git-server

# Should see: "Database: postgres (zixiao_git)"

# Connect to PostgreSQL
docker exec -it zixiao-postgres psql -U zixiao -d zixiao_git

# List tables
\dt
```

## SQL Server Deployment

For Windows-based deployments or enterprise environments:

### 1. Configure Environment

Edit `.env`:
```bash
JWT_SECRET=your-super-secret-key

# Database type
DB_TYPE=sqlserver
DB_HOST=sqlserver
DB_PORT=1433
DB_NAME=zixiao_git
DB_USER=sa
DB_PASSWORD=YourStrong!Password123

# SQL Server container settings
SQL_SERVER_PASSWORD=YourStrong!Password123
```

### 2. Update server.yaml

```yaml
database:
  type: sqlserver
  host: sqlserver
  port: 1433
  name: zixiao_git
  user: sa
  password: YourStrong!Password123
```

### 3. Start Services

```bash
# Uncomment depends_on in docker-compose.yml
# Under git-server service:
depends_on:
  - sqlserver

# Start with SQL Server profile
docker-compose --profile sqlserver up -d
```

### 4. Create Database

```bash
# Connect to SQL Server
docker exec -it zixiao-sqlserver /opt/mssql-tools/bin/sqlcmd \
  -S localhost -U sa -P 'YourStrong!Password123'

# Create database
CREATE DATABASE zixiao_git;
GO

# List databases
SELECT name FROM sys.databases;
GO
```

## External Database

To use an external database (not in Docker):

### PostgreSQL External

Edit `configs/server.yaml`:
```yaml
database:
  type: postgres
  host: your-postgres-host.example.com
  port: 5432
  name: zixiao_git
  user: zixiao
  password: your-password
  sslmode: require  # Enable SSL for external connections
```

Start without database profile:
```bash
docker-compose up -d
```

### SQL Server External

Edit `configs/server.yaml`:
```yaml
database:
  type: sqlserver
  host: your-sqlserver-host.example.com
  port: 1433
  name: zixiao_git
  user: your-user
  password: your-password
```

## Using Adminer (Database UI)

Adminer provides a web interface to manage your database:

```bash
# Start with adminer profile
docker-compose --profile adminer up -d

# Access at http://localhost:8081
```

Connection details:
- **SQLite**: Not supported by Adminer
- **PostgreSQL**:
  - System: PostgreSQL
  - Server: postgres
  - Username: zixiao
  - Password: (from .env)
  - Database: zixiao_git
- **SQL Server**:
  - System: MS SQL
  - Server: sqlserver
  - Username: sa
  - Password: (from .env)
  - Database: zixiao_git

## Multiple Profiles

You can combine multiple profiles:

```bash
# PostgreSQL + Redis + Adminer
docker-compose --profile postgres --profile redis --profile adminer up -d

# SQL Server + Adminer
docker-compose --profile sqlserver --profile adminer up -d
```

## Production Deployment

### 1. Use Docker Swarm or Kubernetes

For high availability, use orchestration:

```bash
# Docker Swarm
docker stack deploy -c docker-compose.yml zixiao

# Kubernetes
kubectl apply -f kubernetes/
```

### 2. External Database

Use managed database services:
- AWS RDS (PostgreSQL)
- Azure SQL Database
- Google Cloud SQL

### 3. Persistent Volumes

Ensure data persistence:

```yaml
volumes:
  postgres-data:
    driver: local
    driver_opts:
      type: none
      device: /mnt/data/postgres
      o: bind
```

### 4. Backup Strategy

#### PostgreSQL Backup

```bash
# Backup
docker exec zixiao-postgres pg_dump -U zixiao zixiao_git > backup.sql

# Restore
docker exec -i zixiao-postgres psql -U zixiao zixiao_git < backup.sql

# Automated backup
docker exec zixiao-postgres sh -c 'pg_dump -U zixiao zixiao_git | gzip' > backup-$(date +%Y%m%d).sql.gz
```

#### SQL Server Backup

```bash
# Backup
docker exec zixiao-sqlserver /opt/mssql-tools/bin/sqlcmd \
  -S localhost -U sa -P 'YourStrong!Password123' \
  -Q "BACKUP DATABASE zixiao_git TO DISK='/var/opt/mssql/backup/zixiao_git.bak'"

# Copy backup out
docker cp zixiao-sqlserver:/var/opt/mssql/backup/zixiao_git.bak ./backup/

# Restore
docker exec zixiao-sqlserver /opt/mssql-tools/bin/sqlcmd \
  -S localhost -U sa -P 'YourStrong!Password123' \
  -Q "RESTORE DATABASE zixiao_git FROM DISK='/var/opt/mssql/backup/zixiao_git.bak' WITH REPLACE"
```

### 5. Monitoring

```bash
# View logs
docker-compose logs -f

# Container stats
docker stats

# Database stats
# PostgreSQL
docker exec zixiao-postgres psql -U zixiao -d zixiao_git -c "SELECT count(*) FROM users;"

# SQL Server
docker exec zixiao-sqlserver /opt/mssql-tools/bin/sqlcmd \
  -S localhost -U sa -P 'Password' -Q "SELECT count(*) FROM users;"
```

## Troubleshooting

### Git Server Can't Connect to Database

```bash
# Check network
docker network ls
docker network inspect zixiao_zixiao-network

# Check if database is running
docker-compose ps

# Check database logs
docker-compose logs postgres
docker-compose logs sqlserver

# Verify connection from git-server
docker exec -it zixiao-git-server sh
ping postgres
telnet postgres 5432
```

### PostgreSQL Connection Refused

```bash
# Check PostgreSQL logs
docker-compose logs postgres

# Ensure pg_hba.conf allows connections
docker exec -it zixiao-postgres cat /var/lib/postgresql/data/pg_hba.conf

# Should have:
# host all all all md5
```

### SQL Server Login Failed

```bash
# Check password complexity
# SQL Server requires:
# - At least 8 characters
# - Uppercase, lowercase, numbers, and symbols

# Reset SA password
docker exec -it zixiao-sqlserver /opt/mssql-tools/bin/sqlcmd \
  -S localhost -U sa -P 'OldPassword' \
  -Q "ALTER LOGIN sa WITH PASSWORD='NewPassword123!'"
```

### Container Out of Memory

```bash
# Check memory usage
docker stats

# Limit container memory
# In docker-compose.yml:
services:
  sqlserver:
    mem_limit: 2g
    memswap_limit: 2g
```

## Performance Tuning

### PostgreSQL

Add to `docker-compose.yml`:
```yaml
postgres:
  command:
    - postgres
    - -c
    - shared_buffers=256MB
    - -c
    - max_connections=200
    - -c
    - work_mem=4MB
```

### SQL Server

```yaml
sqlserver:
  environment:
    - MSSQL_MEMORY_LIMIT_MB=2048
```

## Security Best Practices

1. **Change Default Passwords**
   ```bash
   # Generate strong password
   openssl rand -base64 32
   ```

2. **Use Docker Secrets**
   ```yaml
   services:
     postgres:
       secrets:
         - postgres_password
   secrets:
     postgres_password:
       external: true
   ```

3. **Network Isolation**
   ```yaml
   networks:
     frontend:
     backend:
       internal: true
   ```

4. **Read-Only Containers**
   ```yaml
   git-server:
     read_only: true
     tmpfs:
       - /tmp
   ```

5. **Run as Non-Root**
   ```yaml
   git-server:
     user: "1000:1000"
   ```

## Scaling

### Horizontal Scaling

```yaml
git-server:
  deploy:
    replicas: 3
    update_config:
      parallelism: 1
      delay: 10s
```

### Load Balancer

Use Nginx or HAProxy:
```yaml
nginx:
  image: nginx:alpine
  volumes:
    - ./nginx.conf:/etc/nginx/nginx.conf
  ports:
    - "80:80"
  depends_on:
    - git-server
```

## References

- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [PostgreSQL Docker Image](https://hub.docker.com/_/postgres)
- [SQL Server Docker Image](https://hub.docker.com/_/microsoft-mssql-server)
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)
