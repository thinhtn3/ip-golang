package chat

import "errors"

var (
	// ErrSessionNotFound is returned when a chat session doesn't exist
	ErrSessionNotFound = errors.New("chat session not found")
	
	// ErrForbidden is returned when user doesn't own the session
	ErrForbidden = errors.New("forbidden: user does not own session")
	
	// ErrInvalidSessionID is returned when session ID format is invalid
	ErrInvalidSessionID = errors.New("invalid session ID format")
	
	// ErrInvalidMessage is returned when message content is invalid
	ErrInvalidMessage = errors.New("invalid message content")
	
	// ErrAIServiceUnavailable is returned when AI service fails
	ErrAIServiceUnavailable = errors.New("AI service unavailable")

	// ErrInternalServerError is returned when an internal server error occurs
	ErrInternalServerError = errors.New("internal server error")
)