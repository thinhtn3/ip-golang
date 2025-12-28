package services

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/supabase-community/supabase-go"
)


type ChatService struct {
	supabase *supabase.Client
}

func NewChatService(supabase *supabase.Client) *ChatService {
	//constructor for ChatService
	return &ChatService{supabase: supabase}
}

// Handler function to create a new chat session
func (s *ChatService) CreateSession(c context.Context, userID uuid.UUID, questionID uuid.UUID) () {
	//retrieve question name from question_bank table
	chat := map[string]interface{}{
		"user_id": userID,
		"question_id": questionID,
	}

	//create chat session in supabase
	_, _, err := s.supabase.
		From("chat_sessions").
		Insert(chat, true, "", "", "").
		Execute()
	
	if err != nil {
		log.Println("Error creating chat session: ", err)
	}
	log.Println("Chat session created successfully: ", chat)
}