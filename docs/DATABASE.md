# Database Configuration Examples

This file provides configuration examples for different database backends.

## SQLite (Default)

```yaml
database:
  type: sqlite3
  path: ./data/gitserver.db
```

**Pros:**
- No external dependencies
- Simple setup
- Good for development and small deployments

**Cons:**
- Single-file database
- Limited concurrent writes
- Not suitable for high-traffic deployments

## PostgreSQL (Recommended for Production)

```yaml
database:
  type: postgres
  host: localhost
  port: 5432
  name: zixiao_git
  user: zixiao
  password: your_secure_password
  sslmode: disable  # Use 'require' or 'verify-full' in production
```

### SSL Modes

- `disable`: No SSL (development only)
- `require`: SSL required, no verification
- `verify-ca`: SSL with CA verification
- `verify-full`: SSL with full certificate verification

### Setup PostgreSQL

#### Using Docker

```bash
docker run -d \
  --name zixiao-postgres \
  -e POSTGRES_USER=zixiao \
  -e POSTGRES_PASSWORD=your_secure_password \
  -e POSTGRES_DB=zixiao_git \
  -p 5432:5432 \
  postgres:15-alpine
```

#### Manual Setup

```bash
# Create user
sudo -u postgres createuser zixiao -P

# Create database
sudo -u postgres createdb -O zixiao zixiao_git

# Grant privileges
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE zixiao_git TO zixiao;"
```

### Connection Pooling

PostgreSQL supports connection pooling. Configure in code:

```go
DB.SetMaxOpenConns(25)
DB.SetMaxIdleConns(5)
DB.SetConnMaxLifetime(5 * time.Minute)
```

## SQL Server

```yaml
database:
  type: sqlserver
  host: localhost
  port: 1433
  name: zixiao_git
  user: sa
  password: YourStrong!Password
```

### Setup SQL Server

#### Using Docker

```bash
docker run -d \
  --name zixiao-sqlserver \
  -e 'ACCEPT_EULA=Y' \
  -e 'SA_PASSWORD=YourStrong!Password' \
  -p 1433:1433 \
  mcr.microsoft.com/mssql/server:2022-latest
```

#### Create Database

```sql
CREATE DATABASE zixiao_git;
GO

USE zixiao_git;
GO
```

### Connection String Format

The driver uses the following format:
```
sqlserver://username:password@host:port?database=dbname
```

## Environment Variables

You can override database configuration using environment variables:

```bash
# Database type
export DB_TYPE=postgres

# Connection details
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=zixiao_git
export DB_USER=zixiao
export DB_PASSWORD=your_password

# PostgreSQL specific
export DB_SSLMODE=disable
```

Update your config loading to check environment variables:

```go
if dbType := os.Getenv("DB_TYPE"); dbType != "" {
    cfg.Database.Type = dbType
}
if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
    cfg.Database.Host = dbHost
}
// ... and so on
```

## Performance Considerations

### SQLite
- Good for: < 100 concurrent users, < 10 repos/sec
- Max database size: Practically unlimited (tested up to 140TB)
- Backup: Simple file copy

### PostgreSQL
- Good for: > 100 concurrent users, > 100 repos/sec
- Max database size: Unlimited
- Backup: pg_dump, continuous archiving, PITR
- Requires: Separate server, maintenance

### SQL Server
- Good for: Enterprise deployments, Windows environments
- Max database size: 524 PB
- Backup: Built-in backup tools, Always On
- Requires: License, separate server

## Migration Between Databases

### From SQLite to PostgreSQL

1. Export SQLite data:
```bash
sqlite3 data/gitserver.db .dump > dump.sql
```

2. Convert to PostgreSQL format:
```bash
# Replace AUTOINCREMENT with SERIAL
sed -i 's/AUTOINCREMENT/SERIAL/g' dump.sql

# Replace INTEGER PRIMARY KEY with SERIAL PRIMARY KEY
sed -i 's/INTEGER PRIMARY KEY/SERIAL PRIMARY KEY/g' dump.sql
```

3. Import to PostgreSQL:
```bash
psql -U zixiao -d zixiao_git < dump.sql
```

### Using a Migration Tool

Consider using a migration tool like:
- [pgloader](https://github.com/dimitri/pgloader) - SQLite to PostgreSQL
- [go-migrate](https://github.com/golang-migrate/migrate) - Schema migrations

## Troubleshooting

### PostgreSQL Connection Issues

**Problem**: `connection refused`
```bash
# Check if PostgreSQL is running
sudo systemctl status postgresql

# Check if port is open
netstat -an | grep 5432

# Check pg_hba.conf for authentication
sudo cat /etc/postgresql/15/main/pg_hba.conf
```

**Problem**: `password authentication failed`
```bash
# Reset password
sudo -u postgres psql
ALTER USER zixiao WITH PASSWORD 'new_password';
```

### SQL Server Connection Issues

**Problem**: `login failed for user`
```bash
# Enable SQL Server authentication
# In SQL Server Management Studio:
# Server Properties → Security → SQL Server and Windows Authentication mode

# Or via command:
docker exec -it zixiao-sqlserver /opt/mssql-tools/bin/sqlcmd \
  -S localhost -U sa -P 'YourStrong!Password' \
  -Q "EXEC xp_instance_regwrite N'HKEY_LOCAL_MACHINE', N'Software\Microsoft\MSSQLServer\MSSQLServer', N'LoginMode', REG_DWORD, 2"
```

**Problem**: `cannot open database`
```sql
-- Check database exists
SELECT name FROM sys.databases;

-- Create if missing
CREATE DATABASE zixiao_git;
```

## Monitoring

### PostgreSQL

```sql
-- Active connections
SELECT count(*) FROM pg_stat_activity;

-- Database size
SELECT pg_size_pretty(pg_database_size('zixiao_git'));

-- Slow queries
SELECT query, mean_exec_time, calls
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;
```

### SQL Server

```sql
-- Active connections
SELECT COUNT(*) FROM sys.dm_exec_sessions WHERE is_user_process = 1;

-- Database size
EXEC sp_spaceused;

-- Slow queries
SELECT TOP 10
    total_elapsed_time / execution_count AS avg_time,
    execution_count,
    SUBSTRING(st.text, (qs.statement_start_offset/2)+1,
        ((CASE qs.statement_end_offset
            WHEN -1 THEN DATALENGTH(st.text)
            ELSE qs.statement_end_offset
        END - qs.statement_start_offset)/2) + 1) AS query_text
FROM sys.dm_exec_query_stats AS qs
CROSS APPLY sys.dm_exec_sql_text(qs.sql_handle) AS st
ORDER BY avg_time DESC;
```

## Best Practices

1. **Use connection pooling** in production
2. **Enable SSL/TLS** for remote connections
3. **Regular backups** - automated and tested
4. **Monitor performance** - slow queries, connections
5. **Set appropriate timeouts** in your application
6. **Use read replicas** for high-traffic deployments
7. **Implement retry logic** for transient failures

## References

- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [SQL Server Documentation](https://docs.microsoft.com/en-us/sql/)
- [database/sql Package](https://pkg.go.dev/database/sql)
- [lib/pq Driver](https://github.com/lib/pq)
- [go-mssqldb Driver](https://github.com/microsoft/go-mssqldb)
