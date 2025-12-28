package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/supabase-community/gotrue-go/types"
	"github.com/supabase-community/supabase-go"
	"github.com/thinhtn3/ip-golang.git/internal/services"
)
type ChatSessionRequest struct {
	QuestionID string `json:"question_id"`
}

type ChatSessionHandler struct {
	supabase *supabase.Client
}

func NewChatSessionHandler(supabase *supabase.Client) *ChatSessionHandler {
	return &ChatSessionHandler{supabase: supabase}
}

func (h *ChatSessionHandler) CreateSessionFromQuestion(c *gin.Context) {
	log.Println("CreateChatSession Starting")
	rawUser, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	// type assertion to get user id
	user, ok := rawUser.(*types.UserResponse)
	if !ok {
		c.JSON(500, gin.H{"message": "Internal server error"})
		return
	}

	req := ChatSessionRequest{}
	err := c.ShouldBindJSON(&req);
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid request"})
		return
	}
	log.Println("Question ID: ", req.QuestionID)
	log.Println("User: ", user.User.ID)
	chatService := services.NewChatService(h.supabase)
	chatService.CreateSession(c.Request.Context(), user.User.ID, uuid.MustParse(req.QuestionID))
	log.Println("Chat Service: ", chatService)
	c.JSON(200, gin.H{"message": "Chat session created successfully", "user": user})
}