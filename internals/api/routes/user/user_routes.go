package Userroutes

import (
	"ondc-registry/internals/api/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")
	users.GET("", handlers.GetUsers)

}
