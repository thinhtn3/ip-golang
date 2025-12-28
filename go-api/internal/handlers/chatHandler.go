`package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/supabase-community/supabase-go"
)

type ChatSessionRequest struct {
	supabase *supabase.Client
}

func NewChatSessionRequest(supabase *supabase.Client) *ChatSessionRequest {
	return &ChatSessionRequest{supabase: supabase}
}

func CreateChatSession(c *gin.Context) {
	log.Println("CreateChatSession Starting")
	_, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}


}`