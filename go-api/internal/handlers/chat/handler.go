package chat

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	chatDomain "github.com/thinhtn3/ip-golang.git/internal/domain/chat"
	chatService "github.com/thinhtn3/ip-golang.git/internal/services/chat"
	"net/http"
)

type ChatSessionHandler struct {
	chatService *chatService.ChatService
}


// constructor for ChatSessionHandler
func NewChatSessionHandler(chatService *chatService.ChatService) *ChatSessionHandler {
	return &ChatSessionHandler{chatService: chatService}
}


// HANDLER FUNCTION TO CREATE CHAT SESSION 
func (h *ChatSessionHandler) CreateSessionFromQuestion(c *gin.Context) {
	// get user id from context
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}
	
	type ChatSessionRequest struct {
		QuestionID string `json:"question_id"`
	}

	// bind request to go type
	req := ChatSessionRequest{}
	err = bindRequest(c, &req)

	// call service
	session, err := h.chatService.CreateSession(c.Request.Context(), userID, uuid.MustParse(req.QuestionID))
	if err != nil {
		c.JSON(500, gin.H{"message": "Internal server error"})
		return
	}

	c.JSON(200, gin.H{"message": "Chat session created successfully", "session": session})
}


// SENDING MESSAGES 
func (h *ChatSessionHandler) SendMessage(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid request"})
		return
	}
	

	sessionID, err := parseSessionIDFromContext(c)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid session ID format"})
		return
	}

	type MessageRequest struct {
		Message       string `json:"message"`
		Role          string `json:"role"`
	}

	req := MessageRequest{}
	err = bindRequest(c, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	chat, err := h.chatService.SendMessage(c, userID, sessionID, req.Message, req.Role);
	if err != nil {
		if err == chatDomain.ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden: User does not own session"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Succesfully sent", "chat": chat})
}


// GET MESSAGES INITIAL LOAD
func (h *ChatSessionHandler) GetMessages(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	sessionID, err := parseSessionIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid session ID format"})
		return
	}

	messages, err := h.chatService.GetMessages(c.Request.Context(), userID, sessionID, 0)
	if err != nil {
		if err == chatDomain.ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden: User does not own session"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully retrieved messages", "messages": messages})
}