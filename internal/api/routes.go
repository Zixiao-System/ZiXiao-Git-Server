package api

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(r *gin.Engine) {
	// CORS middleware
	r.Use(CORSMiddleware())

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Public routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", Register)
			auth.POST("/login", Login)
		}

		// User routes
		users := v1.Group("/users")
		{
			users.GET("/:username", GetUser)
			users.GET("/:owner/repos", OptionalAuthMiddleware(), ListRepositories)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(AuthMiddleware())
		{
			// Current user
			protected.GET("/user", GetCurrentUser)

			// Repositories
			repos := protected.Group("/repos")
			{
				repos.POST("", CreateRepository)
				repos.GET("/:owner/:repo", OptionalAuthMiddleware(), GetRepository)
				repos.DELETE("/:owner/:repo", DeleteRepository)

				// Collaborators
				repos.POST("/:owner/:repo/collaborators", AddCollaborator)
				repos.DELETE("/:owner/:repo/collaborators/:username", RemoveCollaborator)
			}
		}
	}

	// Git HTTP protocol routes
	git := r.Group("/:owner/:repo")
	git.Use(OptionalAuthMiddleware())
	{
		git.GET("/info/refs", GitInfoRefs)
		git.POST("/git-receive-pack", GitReceivePack)
		git.POST("/git-upload-pack", GitUploadPack)
	}
}
