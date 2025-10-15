package repository

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/zixiao/git-server/internal/config"
	"github.com/zixiao/git-server/internal/database"
	"github.com/zixiao/git-server/internal/models"
	"github.com/zixiao/git-server/pkg/gitcore"
)

var (
	ErrRepoExists    = fmt.Errorf("repository already exists")
	ErrRepoNotFound  = fmt.Errorf("repository not found")
	ErrAccessDenied  = fmt.Errorf("access denied")
	ErrInvalidName   = fmt.Errorf("invalid repository name")
)

// Create creates a new repository
func Create(ownerID int64, name, description string, isPrivate bool) (*models.Repository, error) {
	if name == "" || len(name) > 100 {
		return nil, ErrInvalidName
	}

	// Insert into database
	result, err := database.DB.Exec(`
		INSERT INTO repositories (name, description, owner_id, is_private, default_branch)
		VALUES (?, ?, ?, ?, ?)
	`, name, description, ownerID, isPrivate, "main")

	if err != nil {
		if err.Error() == "UNIQUE constraint failed: repositories.owner_id, repositories.name" {
			return nil, ErrRepoExists
		}
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	repoID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get repository ID: %w", err)
	}

	// Get owner username
	var ownerName string
	err = database.DB.QueryRow("SELECT username FROM users WHERE id = ?", ownerID).Scan(&ownerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get owner username: %w", err)
	}

	// Create repository on disk
	repoPath := config.GlobalConfig.GetRepoPath(ownerName, name)
	repo := gitcore.NewRepository(repoPath)
	defer repo.Free()

	if err := repo.Init(true); err != nil {
		// Rollback database insert
		database.DB.Exec("DELETE FROM repositories WHERE id = ?", repoID)
		return nil, fmt.Errorf("failed to initialize git repository: %w", err)
	}

	// Log activity
	database.DB.Exec(`
		INSERT INTO activities (user_id, repository_id, action, content)
		VALUES (?, ?, ?, ?)
	`, ownerID, repoID, "create", fmt.Sprintf("Created repository %s", name))

	return &models.Repository{
		ID:            repoID,
		Name:          name,
		Description:   description,
		OwnerID:       ownerID,
		OwnerName:     ownerName,
		IsPrivate:     isPrivate,
		DefaultBranch: "main",
		Size:          0,
		Stars:         0,
		Forks:         0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil
}

// Get retrieves a repository by owner and name
func Get(ownerName, repoName string) (*models.Repository, error) {
	var repo models.Repository
	err := database.DB.QueryRow(`
		SELECT r.id, r.name, r.description, r.owner_id, u.username, r.is_private,
		       r.default_branch, r.size, r.stars, r.forks, r.created_at, r.updated_at
		FROM repositories r
		JOIN users u ON r.owner_id = u.id
		WHERE u.username = ? AND r.name = ?
	`, ownerName, repoName).Scan(&repo.ID, &repo.Name, &repo.Description, &repo.OwnerID,
		&repo.OwnerName, &repo.IsPrivate, &repo.DefaultBranch, &repo.Size, &repo.Stars,
		&repo.Forks, &repo.CreatedAt, &repo.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, ErrRepoNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query repository: %w", err)
	}

	return &repo, nil
}

// GetByID retrieves a repository by ID
func GetByID(repoID int64) (*models.Repository, error) {
	var repo models.Repository
	err := database.DB.QueryRow(`
		SELECT r.id, r.name, r.description, r.owner_id, u.username, r.is_private,
		       r.default_branch, r.size, r.stars, r.forks, r.created_at, r.updated_at
		FROM repositories r
		JOIN users u ON r.owner_id = u.id
		WHERE r.id = ?
	`, repoID).Scan(&repo.ID, &repo.Name, &repo.Description, &repo.OwnerID,
		&repo.OwnerName, &repo.IsPrivate, &repo.DefaultBranch, &repo.Size, &repo.Stars,
		&repo.Forks, &repo.CreatedAt, &repo.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, ErrRepoNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query repository: %w", err)
	}

	return &repo, nil
}

// List lists repositories for a user
func List(ownerID int64) ([]*models.Repository, error) {
	rows, err := database.DB.Query(`
		SELECT r.id, r.name, r.description, r.owner_id, u.username, r.is_private,
		       r.default_branch, r.size, r.stars, r.forks, r.created_at, r.updated_at
		FROM repositories r
		JOIN users u ON r.owner_id = u.id
		WHERE r.owner_id = ?
		ORDER BY r.updated_at DESC
	`, ownerID)

	if err != nil {
		return nil, fmt.Errorf("failed to query repositories: %w", err)
	}
	defer rows.Close()

	repos := []*models.Repository{}
	for rows.Next() {
		var repo models.Repository
		err := rows.Scan(&repo.ID, &repo.Name, &repo.Description, &repo.OwnerID,
			&repo.OwnerName, &repo.IsPrivate, &repo.DefaultBranch, &repo.Size,
			&repo.Stars, &repo.Forks, &repo.CreatedAt, &repo.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan repository: %w", err)
		}
		repos = append(repos, &repo)
	}

	return repos, nil
}

// Delete deletes a repository
func Delete(repoID, userID int64) error {
	// Get repository
	repo, err := GetByID(repoID)
	if err != nil {
		return err
	}

	// Check ownership
	if repo.OwnerID != userID {
		return ErrAccessDenied
	}

	// Delete from database
	_, err = database.DB.Exec("DELETE FROM repositories WHERE id = ?", repoID)
	if err != nil {
		return fmt.Errorf("failed to delete repository: %w", err)
	}

	// Delete from disk
	repoPath := config.GlobalConfig.GetRepoPath(repo.OwnerName, repo.Name)
	if err := os.RemoveAll(repoPath); err != nil {
		return fmt.Errorf("failed to delete repository files: %w", err)
	}

	// Log activity
	database.DB.Exec(`
		INSERT INTO activities (user_id, action, content)
		VALUES (?, ?, ?)
	`, userID, "delete", fmt.Sprintf("Deleted repository %s", repo.Name))

	return nil
}

// CheckAccess checks if a user has access to a repository
func CheckAccess(repoID, userID int64, permission string) (bool, error) {
	repo, err := GetByID(repoID)
	if err != nil {
		return false, err
	}

	// Owner has full access
	if repo.OwnerID == userID {
		return true, nil
	}

	// Public repositories allow read access
	if !repo.IsPrivate && permission == "read" {
		return true, nil
	}

	// Check collaborations
	var count int
	err = database.DB.QueryRow(`
		SELECT COUNT(*) FROM collaborations
		WHERE repository_id = ? AND user_id = ? AND permission = ?
	`, repoID, userID, permission).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("failed to check access: %w", err)
	}

	return count > 0, nil
}

// AddCollaborator adds a collaborator to a repository
func AddCollaborator(repoID, userID int64, permission string) error {
	_, err := database.DB.Exec(`
		INSERT INTO collaborations (repository_id, user_id, permission)
		VALUES (?, ?, ?)
	`, repoID, userID, permission)

	if err != nil {
		return fmt.Errorf("failed to add collaborator: %w", err)
	}

	return nil
}

// RemoveCollaborator removes a collaborator from a repository
func RemoveCollaborator(repoID, userID int64) error {
	_, err := database.DB.Exec(`
		DELETE FROM collaborations WHERE repository_id = ? AND user_id = ?
	`, repoID, userID)

	if err != nil {
		return fmt.Errorf("failed to remove collaborator: %w", err)
	}

	return nil
}
