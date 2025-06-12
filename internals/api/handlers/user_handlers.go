package handlers
import (
    "github.com/gin-gonic/gin"
    "net/http"
    "ondc-registry/internals/services"
)

func GetUsers(c *gin.Context) {
    users, err := services.GetAllUsers()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, users)
}

