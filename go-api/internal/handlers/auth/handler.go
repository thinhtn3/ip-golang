package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {
	log.Println("GetProfile")
	user, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}
	
	c.JSON(200, gin.H{"message": "Profile fetched successfully", "user": user})
}