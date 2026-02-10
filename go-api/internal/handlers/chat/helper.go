package chat

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/supabase-community/gotrue-go/types"
)

func getUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	// Helper function to validate and return userID
	rawUser, exists := c.Get("user")
	if !exists {
		return uuid.Nil, errors.New("user not found in context")
	}
	// type assertion from supabase goauth types library
	user, ok := rawUser.(*types.UserResponse)
	if !ok {
		return uuid.Nil, errors.New("user is not a valid type")
	}
	return user.User.ID, nil
}

func parseSessionIDFromContext(c *gin.Context) (uuid.UUID, error) {
	sessionIDStr := c.Param("sessionId")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return uuid.Nil, errors.New("invalid session ID format")
	}
	return sessionID, nil
}

func bindRequest[T any](c *gin.Context, request *T) (error) {
	return c.ShouldBindJSON(request)
}