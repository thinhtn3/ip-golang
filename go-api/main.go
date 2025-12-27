package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/thinhtn3/ip-golang.git/config"
	"github.com/thinhtn3/ip-golang.git/internal/handlers"
	"github.com/thinhtn3/ip-golang.git/internal/middleware"
)

func main() {
	//enable cors for localhost:5173
	
	config.Load()
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))


	user := router.Group("/user")
	user.Use(middleware.AuthMiddleware())
	{
		user.POST("/profile", handlers.GetProfile)
	}

	//health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OK"})
	})

	router.Run(":" + config.AppConfig.Port)
}