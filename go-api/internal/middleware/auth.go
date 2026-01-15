package middleware

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/supabase-community/supabase-go"
)

type AuthMiddleware struct {
	supabase *supabase.Client
}

func NewAuthMiddleware(supabase *supabase.Client) *AuthMiddleware {
	return &AuthMiddleware{supabase: supabase}
}

func (m *AuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("AuthMiddlewar Starting")
		// Get access token from header and trim
		accessToken := c.GetHeader("Authorization")
		accessToken = strings.TrimPrefix(accessToken, "Bearer ")
	
		if accessToken == "" {
			c.AbortWithStatus(401)
			log.Println("Unauthorized: No access token")
			return
		}

		// Verify token (Supabase)
		// Extract user info and pass onto next handler.
		user, err := m.supabase.Auth.WithToken(accessToken).GetUser()

		if err != nil {
			c.AbortWithStatus(401)
			log.Println("Unauthorized: ", err)
			return
		}

		// Store users in context for next handler and continue
		c.Set("user", user)
		log.Println("AuthMiddleware Completed")
		c.Next()
	}
}