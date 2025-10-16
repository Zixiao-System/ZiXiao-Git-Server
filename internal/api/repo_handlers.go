package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zixiao/git-server/internal/auth"
	"github.com/zixiao/git-server/internal/models"
	"github.com/zixiao/git-server/internal/repository"
)

// CreateRepositoryRequest represents a repository creation request
type CreateRepositoryRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description"`
	IsPrivate   bool   `json:"is_private"`
}

// CreateRepository handles repository creation
func CreateRepository(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	var req CreateRepositoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	repo, err := repository.Create(userID.(int64), req.Name, req.Description, req.IsPrivate)
	if err != nil {
		if err == repository.ErrRepoExists {
			c.JSON(http.StatusConflict, gin.H{"error": "repository already exists"})
			return
		}
		if err == repository.ErrInvalidName {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid repository name"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"repository": repo})
}

// GetRepository retrieves a repository
func GetRepository(c *gin.Context) {
	owner := c.Param("owner")
	repoName := c.Param("repo")

	repo, err := repository.Get(owner, repoName)
	if err != nil {
		if err == repository.ErrRepoNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "repository not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check access for private repositories
	if repo.IsPrivate {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}

		hasAccess, err := repository.CheckAccess(repo.ID, userID.(int64), "read")
		if err != nil || !hasAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"repository": repo})
}

// ListRepositories lists repositories for a user
func ListRepositories(c *gin.Context) {
	username := c.Param("username")

	// Get owner user
	user, err := auth.GetUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	repos, err := repository.List(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Filter private repositories if not owner
	userID, authenticated := c.Get("user_id")
	if !authenticated || userID.(int64) != user.ID {
		filtered := []*models.Repository{}
		for _, repo := range repos {
			if !repo.IsPrivate {
				filtered = append(filtered, repo)
			}
		}
		repos = filtered
	}

	c.JSON(http.StatusOK, gin.H{"repositories": repos})
}

// DeleteRepository deletes a repository
func DeleteRepository(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	owner := c.Param("owner")
	repoName := c.Param("repo")

	repo, err := repository.Get(owner, repoName)
	if err != nil {
		if err == repository.ErrRepoNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "repository not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = repository.Delete(repo.ID, userID.(int64))
	if err != nil {
		if err == repository.ErrAccessDenied {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "repository deleted"})
}

// AddCollaboratorRequest represents a request to add a collaborator
type AddCollaboratorRequest struct {
	Username   string `json:"username" binding:"required"`
	Permission string `json:"permission" binding:"required,oneof=read write admin"`
}

// AddCollaborator adds a collaborator to a repository
func AddCollaborator(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	owner := c.Param("owner")
	repoName := c.Param("repo")

	repo, err := repository.Get(owner, repoName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "repository not found"})
		return
	}

	// Check if user is owner
	if repo.OwnerID != userID.(int64) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	var req AddCollaboratorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get collaborator user
	collabUser, err := auth.GetUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	err = repository.AddCollaborator(repo.ID, collabUser.ID, req.Permission)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "collaborator added"})
}

// RemoveCollaborator removes a collaborator from a repository
func RemoveCollaborator(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	owner := c.Param("owner")
	repoName := c.Param("repo")
	username := c.Param("username")

	repo, err := repository.Get(owner, repoName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "repository not found"})
		return
	}

	// Check if user is owner
	if repo.OwnerID != userID.(int64) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// Get collaborator user
	collabUser, err := auth.GetUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	err = repository.RemoveCollaborator(repo.ID, collabUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "collaborator removed"})
}
