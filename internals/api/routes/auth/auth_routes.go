package AuthRoutes

import (
	"ondc-registry/internals/api/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	auth.GET("/github", handlers.GitHubLogin)
	auth.GET("/callback", handlers.GitHubCallback)
}
