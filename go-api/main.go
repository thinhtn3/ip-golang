package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/thinhtn3/ip-golang.git/config"
	"github.com/thinhtn3/ip-golang.git/internal/middleware"
	chatService "github.com/thinhtn3/ip-golang.git/internal/services/chat"
	chatHandler "github.com/thinhtn3/ip-golang.git/internal/handlers/chat"
	authHandler "github.com/thinhtn3/ip-golang.git/internal/handlers/auth"
	chatrepo "github.com/thinhtn3/ip-golang.git/internal/infra/supabase/chat"
	chatAIClient "github.com/thinhtn3/ip-golang.git/internal/infra/llm/chat"
)

func main() {
	cfg := config.Load()
	//init client and service layer for dependency injection
	supabaseClient := config.InitSupabase(cfg.SupabaseURL, cfg.SupabaseServiceKey)
	repo := chatrepo.NewChatRepo(supabaseClient)
	aiClient := chatAIClient.NewChatAIClient(repo)
	chatSvc := chatService.NewChatService(repo, aiClient)
	chatHdl := chatHandler.NewChatSessionHandler(chatSvc)


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
		user.POST("/profile", authHandler.GetProfile)
	}

	chat := router.Group("/chat")
	chat.Use(middleware.NewAuthMiddleware(supabaseClient).Handle())
	{
		chat.POST("/create", chatHdl.CreateSessionFromQuestion)
		chat.POST("/sessions/:sessionId/messages", chatHdl.SendMessage)
		chat.GET("/sessions/:sessionId/messages", chatHdl.InitialLoadMessages)
	}

	//health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OK"})
	})

	router.Run(":" + cfg.Port)
}