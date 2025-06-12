package registryroutes

import (
	"ondc-registry/internals/config/container"

	"github.com/gin-gonic/gin"
)

// RegistryRoutes registers all registry-related routes
func RegistryRoutes(r *gin.RouterGroup, c *container.Container) {
	registries := r.Group("/registry")
	{
		registries.GET("/all", c.RegistryHdl.GetAllRegistries)
		registries.POST("/create", c.RegistryHdl.CreateRegistry)
		registries.PUT("/:id", c.RegistryHdl.UpdateRegistry)
		registries.DELETE("/:id", c.RegistryHdl.DeleteRegistry)
	}
}
