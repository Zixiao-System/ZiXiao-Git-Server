package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Git      GitConfig      `yaml:"git"`
	Security SecurityConfig `yaml:"security"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"` // debug, release
}

type DatabaseConfig struct {
	Type     string `yaml:"type"`     // sqlite3, postgres, sqlserver
	Path     string `yaml:"path"`     // For SQLite
	Host     string `yaml:"host"`     // For PostgreSQL/SQL Server
	Port     int    `yaml:"port"`     // For PostgreSQL/SQL Server
	Name     string `yaml:"name"`     // Database name
	User     string `yaml:"user"`     // Username
	Password string `yaml:"password"` // Password
	SSLMode  string `yaml:"sslmode"`  // For PostgreSQL (disable, require, verify-ca, verify-full)
}

type GitConfig struct {
	RepoPath     string `yaml:"repo_path"`
	MaxRepoSize  int64  `yaml:"max_repo_size"`  // in MB
	MaxFileSize  int64  `yaml:"max_file_size"`  // in MB
	AllowedTypes []string `yaml:"allowed_types"` // file extensions
}

type SecurityConfig struct {
	JWTSecret     string `yaml:"jwt_secret"`
	JWTExpiration int    `yaml:"jwt_expiration"` // in hours
	PasswordMin   int    `yaml:"password_min"`
	EnableSSH     bool   `yaml:"enable_ssh"`
	SSHPort       int    `yaml:"ssh_port"`
}

var GlobalConfig *Config

// Load loads configuration from file
func Load(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set defaults
	if cfg.Server.Host == "" {
		cfg.Server.Host = "0.0.0.0"
	}
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if cfg.Server.Mode == "" {
		cfg.Server.Mode = "release"
	}
	if cfg.Database.Type == "" {
		cfg.Database.Type = "sqlite"
	}
	if cfg.Database.Path == "" {
		cfg.Database.Path = "./data/gitserver.db"
	}
	if cfg.Git.RepoPath == "" {
		cfg.Git.RepoPath = "./data/repositories"
	}
	if cfg.Git.MaxRepoSize == 0 {
		cfg.Git.MaxRepoSize = 1024 // 1GB default
	}
	if cfg.Git.MaxFileSize == 0 {
		cfg.Git.MaxFileSize = 100 // 100MB default
	}
	if cfg.Security.JWTExpiration == 0 {
		cfg.Security.JWTExpiration = 24 // 24 hours
	}
	if cfg.Security.PasswordMin == 0 {
		cfg.Security.PasswordMin = 8
	}
	if cfg.Security.SSHPort == 0 {
		cfg.Security.SSHPort = 2222
	}

	GlobalConfig = &cfg
	return &cfg, nil
}

// Save saves configuration to file
func (c *Config) Save(configPath string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetRepoPath returns the full path for a repository
func (c *Config) GetRepoPath(owner, repoName string) string {
	return fmt.Sprintf("%s/%s/%s.git", c.Git.RepoPath, owner, repoName)
}
