package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zixiao/git-server/internal/config"
	"github.com/zixiao/git-server/internal/database"
	"github.com/zixiao/git-server/internal/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrInvalidCredentials is returned when username or password is incorrect
	ErrInvalidCredentials = errors.New("invalid username or password")
	// ErrUserExists is returned when trying to create a user that already exists
	ErrUserExists = errors.New("user already exists")
	// ErrUserNotFound is returned when a user cannot be found
	ErrUserNotFound = errors.New("user not found")
	// ErrInvalidToken is returned when a token is invalid or expired
	ErrInvalidToken = errors.New("invalid token")
)

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// VerifyPassword verifies a password against a hash
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken generates a JWT token for a user
func GenerateToken(user *models.User) (string, error) {
	cfg := config.GlobalConfig
	expirationTime := time.Now().Add(time.Duration(cfg.Security.JWTExpiration) * time.Hour)

	claims := &JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "zixiao-git-server",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Security.JWTSecret))
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string) (*JWTClaims, error) {
	cfg := config.GlobalConfig

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.Security.JWTSecret), nil
		})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// Register creates a new user account
func Register(username, email, password, fullName string) (*models.User, error) {
	// Validate password length
	if len(password) < config.GlobalConfig.Security.PasswordMin {
		return nil, fmt.Errorf("password must be at least %d characters",
			config.GlobalConfig.Security.PasswordMin)
	}

	// Hash password
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Insert user into database
	result, err := database.DB.Exec(`
		INSERT INTO users (username, email, password, full_name, is_admin, is_active)
		VALUES (?, ?, ?, ?, ?, ?)
	`, username, email, hashedPassword, fullName, false, true)

	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.username" ||
			err.Error() == "UNIQUE constraint failed: users.email" {
			return nil, ErrUserExists
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	return &models.User{
		ID:        userID,
		Username:  username,
		Email:     email,
		FullName:  fullName,
		IsAdmin:   false,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// Login authenticates a user and returns a token
func Login(username, password string) (string, *models.User, error) {
	var user models.User
	err := database.DB.QueryRow(`
		SELECT id, username, email, password, full_name, is_admin, is_active, created_at, updated_at
		FROM users WHERE username = ? AND is_active = 1
	`, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password,
		&user.FullName, &user.IsAdmin, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return "", nil, ErrInvalidCredentials
	}
	if err != nil {
		return "", nil, fmt.Errorf("failed to query user: %w", err)
	}

	// Verify password
	if !VerifyPassword(password, user.Password) {
		return "", nil, ErrInvalidCredentials
	}

	// Generate token
	token, err := GenerateToken(&user)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Don't return password hash
	user.Password = ""

	return token, &user, nil
}

// GetUserByID retrieves a user by ID
func GetUserByID(userID int64) (*models.User, error) {
	var user models.User
	err := database.DB.QueryRow(`
		SELECT id, username, email, full_name, is_admin, is_active, created_at, updated_at
		FROM users WHERE id = ?
	`, userID).Scan(&user.ID, &user.Username, &user.Email, &user.FullName,
		&user.IsAdmin, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	return &user, nil
}

// GetUserByUsername retrieves a user by username
func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := database.DB.QueryRow(`
		SELECT id, username, email, full_name, is_admin, is_active, created_at, updated_at
		FROM users WHERE username = ?
	`, username).Scan(&user.ID, &user.Username, &user.Email, &user.FullName,
		&user.IsAdmin, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	return &user, nil
}

// GenerateAccessToken generates a random access token
func GenerateAccessToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// CreateAccessToken creates a new access token for a user
func CreateAccessToken(userID int64, name string, expiresAt *time.Time) (*models.AccessToken, error) {
	token, err := GenerateAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	var expiresAtSQL interface{}
	if expiresAt != nil {
		expiresAtSQL = expiresAt
	}

	result, err := database.DB.Exec(`
		INSERT INTO access_tokens (user_id, token, name, expires_at)
		VALUES (?, ?, ?, ?)
	`, userID, token, name, expiresAtSQL)

	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	tokenID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get token ID: %w", err)
	}

	accessToken := &models.AccessToken{
		ID:        tokenID,
		UserID:    userID,
		Token:     token,
		Name:      name,
		CreatedAt: time.Now(),
	}

	if expiresAt != nil {
		accessToken.ExpiresAt = *expiresAt
	}

	return accessToken, nil
}

// ValidateAccessToken validates an access token and returns the user
func ValidateAccessToken(token string) (*models.User, error) {
	var accessToken models.AccessToken
	err := database.DB.QueryRow(`
		SELECT id, user_id, token, name, expires_at, created_at
		FROM access_tokens WHERE token = ?
	`, token).Scan(&accessToken.ID, &accessToken.UserID, &accessToken.Token,
		&accessToken.Name, &accessToken.ExpiresAt, &accessToken.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, ErrInvalidToken
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query access token: %w", err)
	}

	// Check if token is expired
	if !accessToken.ExpiresAt.IsZero() && accessToken.ExpiresAt.Before(time.Now()) {
		return nil, ErrInvalidToken
	}

	// Get user
	return GetUserByID(accessToken.UserID)
}
