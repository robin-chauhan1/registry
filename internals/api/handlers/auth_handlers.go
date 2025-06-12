package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ondc-registry/internals/config/database"
	"ondc-registry/internals/models"
	"ondc-registry/internals/utils"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	oauthConfig      *oauth2.Config
	oauthStateString = "randomstate" // Replace with a random string for CSRF protection
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
	}

	oauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/api/auth/callback",
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
}

func GitHubLogin(c *gin.Context) {
	url := oauthConfig.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	fmt.Println("url", url)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GitHubCallback(c *gin.Context) {
	if state := c.Query("state"); state != oauthStateString {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid OAuth state"})
		return
	}

	code := c.Query("code")
	if code == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
		return
	}

	token, err := oauthConfig.Exchange(c.Request.Context(), code)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to exchange token: %v", err)})
		return
	}

	client := oauthConfig.Client(c.Request.Context(), token)
	resp, err := client.Get("https://api.github.com/user")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get user info: %v", err)})
		return
	}
	defer resp.Body.Close()

	var userData map[string]any

	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to parse user data: %v", err)})
		return
	}

	user := models.User{
		Name:  fmt.Sprintf("%v", userData["name"]),
		Email: fmt.Sprintf("%v", userData["email"]),
	}

	if err := database.GetDB().Create(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save user: %v", err)})
		return
	}

	// c.JSON(http.StatusOK, gin.H{
	// 	"message": "Successfully authenticated",
	// 	"user":    userData,
	// })
	jwtToken, err := utils.GenerateJwtToken(fmt.Sprintf("%v", user.ID), "github")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate token: %v", err)})
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, "http://localhost:8080/token?"+jwtToken)
}
