package auth

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {
	log.Println("GetProfile")
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Profile fetched successfully", "user": user})
}