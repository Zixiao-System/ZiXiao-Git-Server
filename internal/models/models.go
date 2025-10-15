package models

import (
	"time"
)

// User represents a user account
type User struct {
	ID        int64     `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"` // Hashed password
	FullName  string    `json:"full_name" db:"full_name"`
	IsAdmin   bool      `json:"is_admin" db:"is_admin"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Repository represents a git repository
type Repository struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	OwnerID     int64     `json:"owner_id" db:"owner_id"`
	OwnerName   string    `json:"owner_name" db:"-"` // Joined field
	IsPrivate   bool      `json:"is_private" db:"is_private"`
	DefaultBranch string  `json:"default_branch" db:"default_branch"`
	Size        int64     `json:"size" db:"size"` // in bytes
	Stars       int       `json:"stars" db:"stars"`
	Forks       int       `json:"forks" db:"forks"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// SSHKey represents a user's SSH public key
type SSHKey struct {
	ID          int64     `json:"id" db:"id"`
	UserID      int64     `json:"user_id" db:"user_id"`
	Title       string    `json:"title" db:"title"`
	Key         string    `json:"key" db:"key"` // Public key content
	Fingerprint string    `json:"fingerprint" db:"fingerprint"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// Collaboration represents repository access permissions
type Collaboration struct {
	ID           int64     `json:"id" db:"id"`
	RepositoryID int64     `json:"repository_id" db:"repository_id"`
	UserID       int64     `json:"user_id" db:"user_id"`
	Permission   string    `json:"permission" db:"permission"` // read, write, admin
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// AccessToken represents an API access token
type AccessToken struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	Name      string    `json:"name" db:"name"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Activity represents user or repository activity
type Activity struct {
	ID           int64     `json:"id" db:"id"`
	UserID       int64     `json:"user_id" db:"user_id"`
	RepositoryID int64     `json:"repository_id" db:"repository_id"`
	Action       string    `json:"action" db:"action"` // create, push, fork, star, etc.
	RefName      string    `json:"ref_name" db:"ref_name"`
	Content      string    `json:"content" db:"content"` // JSON encoded details
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
