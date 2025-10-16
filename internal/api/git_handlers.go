package api

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zixiao/git-server/internal/config"
	"github.com/zixiao/git-server/internal/repository"
	"github.com/zixiao/git-server/pkg/gitcore"
)

// GitInfoRefs handles git info/refs request
func GitInfoRefs(c *gin.Context) {
	owner := c.Param("owner")
	repoName := c.Param("repo")
	service := c.Query("service")

	// Get repository
	repo, err := repository.Get(owner, repoName)
	if err != nil {
		c.String(http.StatusNotFound, "Repository not found")
		return
	}

	// Check access for private repositories
	if repo.IsPrivate {
		userID, exists := c.Get("user_id")
		if !exists {
			c.Header("WWW-Authenticate", "Basic realm=\"Git\"")
			c.String(http.StatusUnauthorized, "Authentication required")
			return
		}

		hasAccess, err := repository.CheckAccess(repo.ID, userID.(int64), "read")
		if err != nil || !hasAccess {
			c.String(http.StatusForbidden, "Access denied")
			return
		}
	}

	// Get repository path
	repoPath := config.GlobalConfig.GetRepoPath(owner, repoName)
	gitRepo := gitcore.NewRepository(repoPath)
	defer gitRepo.Free()

	// Get all refs
	refs, err := gitRepo.ListRefs()
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to list refs")
		return
	}

	// Build refs map
	refsMap := make(map[string]string)
	for _, ref := range refs {
		sha, err := gitRepo.GetRef(ref)
		if err == nil && sha != "" {
			refsMap["refs/"+ref] = sha
		}
	}

	// Create advertisement
	adv, err := gitcore.CreateRefAdvertisement(refsMap, service)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to create advertisement")
		return
	}

	// Set headers
	c.Header("Content-Type", "application/x-"+service+"-advertisement")
	c.Header("Cache-Control", "no-cache")

	c.Data(http.StatusOK, "application/x-"+service+"-advertisement", adv)
}

// GitReceivePack handles git push (receive-pack)
func GitReceivePack(c *gin.Context) {
	owner := c.Param("owner")
	repoName := c.Param("repo")

	// Get repository
	repo, err := repository.Get(owner, repoName)
	if err != nil {
		c.String(http.StatusNotFound, "Repository not found")
		return
	}

	// Check write access
	userID, exists := c.Get("user_id")
	if !exists {
		c.Header("WWW-Authenticate", "Basic realm=\"Git\"")
		c.String(http.StatusUnauthorized, "Authentication required")
		return
	}

	hasAccess, err := repository.CheckAccess(repo.ID, userID.(int64), "write")
	if err != nil || !hasAccess {
		c.String(http.StatusForbidden, "Access denied")
		return
	}

	// Read pack data
	packData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to read pack data")
		return
	}

	// Get repository path
	repoPath := config.GlobalConfig.GetRepoPath(owner, repoName)
	gitRepo := gitcore.NewRepository(repoPath)
	defer gitRepo.Free()

	// Receive pack
	if err := gitRepo.ReceivePack(packData); err != nil {
		c.String(http.StatusInternalServerError, "Failed to receive pack")
		return
	}

	// Return success
	response := gitcore.PktLine("unpack ok\n") + gitcore.FlushPkt()
	c.Header("Content-Type", "application/x-git-receive-pack-result")
	c.String(http.StatusOK, response)
}

// GitUploadPack handles git pull/fetch (upload-pack)
func GitUploadPack(c *gin.Context) {
	owner := c.Param("owner")
	repoName := c.Param("repo")

	// Get repository
	repo, err := repository.Get(owner, repoName)
	if err != nil {
		c.String(http.StatusNotFound, "Repository not found")
		return
	}

	// Check read access for private repositories
	if repo.IsPrivate {
		userID, exists := c.Get("user_id")
		if !exists {
			c.Header("WWW-Authenticate", "Basic realm=\"Git\"")
			c.String(http.StatusUnauthorized, "Authentication required")
			return
		}

		hasAccess, err := repository.CheckAccess(repo.ID, userID.(int64), "read")
		if err != nil || !hasAccess {
			c.String(http.StatusForbidden, "Access denied")
			return
		}
	}

	// Read request body (wants and haves)
	_, err = io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to read request")
		return
	}

	// Parse wants and haves (simplified)
	wants := []string{}
	haves := []string{}

	// Get repository path
	repoPath := config.GlobalConfig.GetRepoPath(owner, repoName)
	gitRepo := gitcore.NewRepository(repoPath)
	defer gitRepo.Free()

	// Upload pack
	packData, err := gitRepo.UploadPack(wants, haves)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to upload pack")
		return
	}

	// Return pack data
	c.Header("Content-Type", "application/x-git-upload-pack-result")
	c.Data(http.StatusOK, "application/x-git-upload-pack-result", packData)
}
