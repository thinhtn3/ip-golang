package models

import (
	"time"

	"github.com/google/uuid"
)

// ChatSession represents the chat_sessions table
type ChatSession struct {
	ID           uuid.UUID    `json:"id"` 
	UserID       uuid.UUID    `json:"user_id"`
	QuestionID   uuid.UUID    `json:"question_id"`
	QuestionName string    `json:"question_name"`
	CreatedAt    time.Time `json:"created_at"`
}

// Question represents the question_bank table
type Question struct {
	ID         uuid.UUID    `json:"id"`
	Slug       string    `json:"slug"`
	Title      string    `json:"title"`
	URL        string    `json:"url"`
	Difficulty string    `json:"difficulty"`
	CreatedAt  time.Time `json:"created_at"`
}

// Message represents the messages table
type Message struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	ChatSessionID uuid.UUID `json:"chat_session_id"`
	Role          string `json:"role"`
	Message       string `json:"message"`
}

// Verifying session ownership
type Row struct {
	ID uuid.UUID `json:"id"`
}