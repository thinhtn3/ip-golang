package middleware

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thinhtn3/ip-golang.git/config"
)


func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		log.Println("AuthMiddleware")
		// Get access token from header and trim
		accessToken := c.GetHeader("Authorization")
		accessToken = strings.TrimPrefix(accessToken, "Bearer ")
	
		// Verify token (Supabase)
		user, err := config.SupabaseClient.Auth.WithToken(accessToken).GetUser()
		if accessToken == "" || err != nil {
			c.AbortWithStatus(401)
			log.Println("Unauthorized: ", err)
			return
		}

		// Store users in context for next handler and continue
		c.Set("user", user)
		c.Next()
	}
}