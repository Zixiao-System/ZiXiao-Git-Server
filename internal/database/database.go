package database

import (
	"database/sql"
	"fmt"
)

// Config holds database configuration
type Config struct {
	Type     string
	Path     string // For SQLite
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string // For PostgreSQL
}

// DB is the global database connection
var DB *sql.DB

// Init initializes the database connection
func Init(cfg Config) error {
	dsn, err := buildDSN(cfg)
	if err != nil {
		return fmt.Errorf("failed to build DSN: %w", err)
	}

	DB, err = sql.Open(cfg.Type, dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)

	// Create tables
	if err := createTables(cfg.Type); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

// buildDSN builds the data source name for the given database type
func buildDSN(cfg Config) (string, error) {
	switch cfg.Type {
	case "sqlite3":
		return cfg.Path, nil

	case "postgres":
		sslMode := cfg.SSLMode
		if sslMode == "" {
			sslMode = "disable"
		}
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, sslMode), nil

	case "sqlserver", "mssql":
		// SQL Server connection string
		return fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name), nil

	default:
		return "", fmt.Errorf("unsupported database type: %s", cfg.Type)
	}
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// createTables creates all necessary database tables
func createTables(dbType string) error {
	schema := getSchema(dbType)
	_, err := DB.Exec(schema)
	return err
}

// getSchema returns the appropriate schema for the database type
func getSchema(dbType string) string {
	switch dbType {
	case "postgres":
		return getPostgreSQLSchema()
	case "sqlserver", "mssql":
		return getSQLServerSchema()
	default: // sqlite3
		return getSQLiteSchema()
	}
}

// getSQLiteSchema returns SQLite schema
func getSQLiteSchema() string {
	return `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		full_name TEXT,
		is_admin BOOLEAN DEFAULT 0,
		is_active BOOLEAN DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS repositories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT,
		owner_id INTEGER NOT NULL,
		is_private BOOLEAN DEFAULT 0,
		default_branch TEXT DEFAULT 'main',
		size INTEGER DEFAULT 0,
		stars INTEGER DEFAULT 0,
		forks INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE(owner_id, name)
	);

	CREATE TABLE IF NOT EXISTS ssh_keys (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		key TEXT NOT NULL,
		fingerprint TEXT NOT NULL UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS collaborations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		repository_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		permission TEXT NOT NULL DEFAULT 'read',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (repository_id) REFERENCES repositories(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE(repository_id, user_id)
	);

	CREATE TABLE IF NOT EXISTS access_tokens (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		token TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL,
		expires_at DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS activities (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		repository_id INTEGER,
		action TEXT NOT NULL,
		ref_name TEXT,
		content TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (repository_id) REFERENCES repositories(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_repositories_owner ON repositories(owner_id);
	CREATE INDEX IF NOT EXISTS idx_ssh_keys_user ON ssh_keys(user_id);
	CREATE INDEX IF NOT EXISTS idx_collaborations_repo ON collaborations(repository_id);
	CREATE INDEX IF NOT EXISTS idx_collaborations_user ON collaborations(user_id);
	CREATE INDEX IF NOT EXISTS idx_activities_user ON activities(user_id);
	CREATE INDEX IF NOT EXISTS idx_activities_repo ON activities(repository_id);
	`
}

// getPostgreSQLSchema returns PostgreSQL schema
func getPostgreSQLSchema() string {
	return `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) NOT NULL UNIQUE,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		full_name VARCHAR(255),
		is_admin BOOLEAN DEFAULT FALSE,
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS repositories (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		owner_id INTEGER NOT NULL,
		is_private BOOLEAN DEFAULT FALSE,
		default_branch VARCHAR(255) DEFAULT 'main',
		size BIGINT DEFAULT 0,
		stars INTEGER DEFAULT 0,
		forks INTEGER DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE(owner_id, name)
	);

	CREATE TABLE IF NOT EXISTS ssh_keys (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL,
		title VARCHAR(255) NOT NULL,
		key TEXT NOT NULL,
		fingerprint VARCHAR(255) NOT NULL UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS collaborations (
		id SERIAL PRIMARY KEY,
		repository_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		permission VARCHAR(50) NOT NULL DEFAULT 'read',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (repository_id) REFERENCES repositories(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE(repository_id, user_id)
	);

	CREATE TABLE IF NOT EXISTS access_tokens (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL,
		token VARCHAR(255) NOT NULL UNIQUE,
		name VARCHAR(255) NOT NULL,
		expires_at TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS activities (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL,
		repository_id INTEGER,
		action VARCHAR(255) NOT NULL,
		ref_name VARCHAR(255),
		content TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (repository_id) REFERENCES repositories(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_repositories_owner ON repositories(owner_id);
	CREATE INDEX IF NOT EXISTS idx_ssh_keys_user ON ssh_keys(user_id);
	CREATE INDEX IF NOT EXISTS idx_collaborations_repo ON collaborations(repository_id);
	CREATE INDEX IF NOT EXISTS idx_collaborations_user ON collaborations(user_id);
	CREATE INDEX IF NOT EXISTS idx_activities_user ON activities(user_id);
	CREATE INDEX IF NOT EXISTS idx_activities_repo ON activities(repository_id);
	`
}

// getSQLServerSchema returns SQL Server schema
func getSQLServerSchema() string {
	return `
	IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'users')
	CREATE TABLE users (
		id INT IDENTITY(1,1) PRIMARY KEY,
		username NVARCHAR(255) NOT NULL UNIQUE,
		email NVARCHAR(255) NOT NULL UNIQUE,
		password NVARCHAR(255) NOT NULL,
		full_name NVARCHAR(255),
		is_admin BIT DEFAULT 0,
		is_active BIT DEFAULT 1,
		created_at DATETIME DEFAULT GETDATE(),
		updated_at DATETIME DEFAULT GETDATE()
	);

	IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'repositories')
	CREATE TABLE repositories (
		id INT IDENTITY(1,1) PRIMARY KEY,
		name NVARCHAR(255) NOT NULL,
		description NVARCHAR(MAX),
		owner_id INT NOT NULL,
		is_private BIT DEFAULT 0,
		default_branch NVARCHAR(255) DEFAULT 'main',
		size BIGINT DEFAULT 0,
		stars INT DEFAULT 0,
		forks INT DEFAULT 0,
		created_at DATETIME DEFAULT GETDATE(),
		updated_at DATETIME DEFAULT GETDATE(),
		FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE(owner_id, name)
	);

	IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'ssh_keys')
	CREATE TABLE ssh_keys (
		id INT IDENTITY(1,1) PRIMARY KEY,
		user_id INT NOT NULL,
		title NVARCHAR(255) NOT NULL,
		[key] NVARCHAR(MAX) NOT NULL,
		fingerprint NVARCHAR(255) NOT NULL UNIQUE,
		created_at DATETIME DEFAULT GETDATE(),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'collaborations')
	CREATE TABLE collaborations (
		id INT IDENTITY(1,1) PRIMARY KEY,
		repository_id INT NOT NULL,
		user_id INT NOT NULL,
		permission NVARCHAR(50) NOT NULL DEFAULT 'read',
		created_at DATETIME DEFAULT GETDATE(),
		FOREIGN KEY (repository_id) REFERENCES repositories(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE(repository_id, user_id)
	);

	IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'access_tokens')
	CREATE TABLE access_tokens (
		id INT IDENTITY(1,1) PRIMARY KEY,
		user_id INT NOT NULL,
		token NVARCHAR(255) NOT NULL UNIQUE,
		name NVARCHAR(255) NOT NULL,
		expires_at DATETIME,
		created_at DATETIME DEFAULT GETDATE(),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'activities')
	CREATE TABLE activities (
		id INT IDENTITY(1,1) PRIMARY KEY,
		user_id INT NOT NULL,
		repository_id INT,
		action NVARCHAR(255) NOT NULL,
		ref_name NVARCHAR(255),
		content NVARCHAR(MAX),
		created_at DATETIME DEFAULT GETDATE(),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (repository_id) REFERENCES repositories(id) ON DELETE CASCADE
	);

	IF NOT EXISTS (SELECT * FROM sys.indexes WHERE name = 'idx_repositories_owner')
	CREATE INDEX idx_repositories_owner ON repositories(owner_id);

	IF NOT EXISTS (SELECT * FROM sys.indexes WHERE name = 'idx_ssh_keys_user')
	CREATE INDEX idx_ssh_keys_user ON ssh_keys(user_id);

	IF NOT EXISTS (SELECT * FROM sys.indexes WHERE name = 'idx_collaborations_repo')
	CREATE INDEX idx_collaborations_repo ON collaborations(repository_id);

	IF NOT EXISTS (SELECT * FROM sys.indexes WHERE name = 'idx_collaborations_user')
	CREATE INDEX idx_collaborations_user ON collaborations(user_id);

	IF NOT EXISTS (SELECT * FROM sys.indexes WHERE name = 'idx_activities_user')
	CREATE INDEX idx_activities_user ON activities(user_id);

	IF NOT EXISTS (SELECT * FROM sys.indexes WHERE name = 'idx_activities_repo')
	CREATE INDEX idx_activities_repo ON activities(repository_id);
	`
}
