package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/thinhtn3/ip-golang.git/config"
	"github.com/thinhtn3/ip-golang.git/internal/handlers"
	"github.com/thinhtn3/ip-golang.git/internal/middleware"
	"github.com/thinhtn3/ip-golang.git/internal/services"
)

func main() {
	//enable cors for localhost:5173
	
	
	cfg := config.Load()
	//init client and service layer for dependency injectionwhy
	supabaseClient := config.InitSupabase(cfg.SupabaseURL, cfg.SupabaseServiceKey)
	chatService := services.NewChatService(supabaseClient)
	chatHandler := handlers.NewChatSessionHandler(chatService)


	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))


	user := router.Group("/user")
	user.Use(middleware.NewAuthMiddleware(supabaseClient).Handle())
	{
		user.POST("/profile", handlers.GetProfile)
	}

	chat := router.Group("/chat")
	chat.Use(middleware.NewAuthMiddleware(supabaseClient).Handle())
	{
		chat.POST("/create", chatHandler.CreateSessionFromQuestion)
		chat.POST("/sessions/:sessionId/messages", chatHandler.SendMessage)
		chat.GET("/sessions/:sessionId/messages", chatHandler.GetMessages)
	}

	//health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OK"})
	})

	router.Run(":" + cfg.Port)
}