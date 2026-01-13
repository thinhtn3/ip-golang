package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/supabase-community/gotrue-go/types"
	"github.com/supabase-community/supabase-go"
	"github.com/thinhtn3/ip-golang.git/internal/services"
)

// DEPENDENCY INJECTION //

type ChatSessionHandler struct {
	supabase *supabase.Client
}
//constructor for ChatSessionHandler
func NewChatSessionHandler(supabase *supabase.Client) *ChatSessionHandler {
	return &ChatSessionHandler{supabase: supabase}
}


// HANDLER FUNCTION TO CREATE CHAT SESSION //
type ChatSessionRequest struct {
	QuestionID string `json:"question_id"`
}

func (h *ChatSessionHandler) CreateSessionFromQuestion(c *gin.Context) {
	rawUser, exists := c.Get("user")
	if (!exists) {
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}
	// type assertion from supabase goauth types library
	user, ok := rawUser.(*types.UserResponse)
	if !ok {
		c.JSON(500, gin.H{"message": "Internal server error"})
		return
	}

	//bind question id to request
	req := ChatSessionRequest{}
	err := c.ShouldBindJSON(&req);
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid request"})
		return
	}

	//create chat session
	chatService := services.NewChatService(h.supabase)
	session, err := chatService.CreateSession(c.Request.Context(), user.User.ID, uuid.MustParse(req.QuestionID))
	if err != nil {
		c.JSON(500, gin.H{"message": "Internal server error"})
		return
	}

	c.JSON(200, gin.H{"message": "Chat session created successfully", "session": session})
}

// SENDING MESSAGES //
type MessageRequest struct {
	Message       string `json:"message"`
	Role          string `json:"role"`
}
func (h *ChatSessionHandler) SendMessage(c *gin.Context) {
	//Get user from context
	rawUser, exists := c.Get("user")
	if (!exists) {
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}
	user, ok := rawUser.(*types.UserResponse)
	if !ok {
		c.JSON(500, gin.H{"message": "Internal server error"})
		return
	}

	//Parse sessionID string from URL
	sessionIDStr := c.Param("sessionId")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid session ID format"})
		return
	}

	//bind messages into request, which can be accessed as req.Message and req.Role
	req := MessageRequest{}
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid request"})
		return
	}

	chatService := services.NewChatService(h.supabase)
	chat, err := chatService.SendMessage(c, user.User.ID, sessionID, req.Message, req.Role);
	if err != nil {
		if err == services.ForbiddenError {
			c.JSON(403, gin.H{"message": "Forbidden: User does not own session"})
			return
		} else {
			c.JSON(500, gin.H{"message": "Internal server error"})
		}
		return
	}
	c.JSON(200, gin.H{"message": "Succesfully sent", "chat": chat})
}

// GET MESSAGES //
func (h *ChatSessionHandler) GetMessages(c *gin.Context) {
	rawUser, exists := c.Get("user")
	if (!exists) {
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}
	user, ok := rawUser.(*types.UserResponse)
	if !ok {
		c.JSON(500, gin.H{"message": "Internal server error"})
		return
	}

	sessionIDStr := c.Param("sessionId")
	sessionID, err := uuid.Parse(sessionIDStr)
	if (err != nil) {
		c.JSON(400, gin.H{"message": "Invalid session ID format"})
	}
	chatService := services.NewChatService(h.supabase)
	messages, err := chatService.GetMessages(c.Request.Context(), user.User.ID, sessionID)
	if err != nil {
		if err == services.ForbiddenError {
			c.JSON(403, gin.H{"message": "Forbidden: User does not own session"})
			return
		} else {
			c.JSON(500, gin.H{"message": "Internal server error"})
		}
		return
}
	c.JSON(200, gin.H{"message": "Successfully retrieved messages", "messages": messages})
}