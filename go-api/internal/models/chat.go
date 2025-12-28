package models

import "time"

// ChatSession represents the chat_sessions table
type ChatSession struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	QuestionID   string    `json:"question_id"`
	QuestionName string    `json:"question_name"`
	CreatedAt    time.Time `json:"created_at"`
}

// Question represents the question_bank table
type Question struct {
	ID         string    `json:"id"`
	Slug       string    `json:"slug"`
	Title      string    `json:"title"`
	URL        string    `json:"url"`
	Difficulty string    `json:"difficulty"`
	CreatedAt  time.Time `json:"created_at"`
}

// Message represents the messages table
type Message struct {
	ID            string `json:"id"`
	UserID        string `json:"user_id"`
	ChatSessionID string `json:"chat_session_id"`
	Role          string `json:"role"`
	Message       string `json:"message"`
}