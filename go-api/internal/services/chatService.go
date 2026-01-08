package services

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/supabase-community/supabase-go"
	"github.com/thinhtn3/ip-golang.git/internal/models"
)

// CONSTS
var ForbiddenError = errors.New("Forbidden: User does not own session")
var InternalServerError = errors.New("Internal server error")

type ChatService struct {
	supabase *supabase.Client
}

func NewChatService(supabase *supabase.Client) *ChatService {
	//constructor for ChatService
	return &ChatService{supabase: supabase}
}

// Handler function to create a new chat session
func (s *ChatService) CreateSession(c context.Context, userID uuid.UUID, questionID uuid.UUID) (*models.ChatSession, error) {
	session, err := s.GetSession(c, userID, questionID)

	if (err != nil) {
		return nil, err
	}

	if (session != nil) {
		log.Println("Found existing session", session.ID, session.QuestionID)
		return session, nil
	}

	//retrieve question name from question_bank table
	chat := map[string]interface{}{
		"user_id": userID,
		"question_id": questionID,
	}

	//create chat session in supabase
	_, _, err = s.supabase.
		From("chat_sessions").
		Insert(chat, true, "", "", "").
		Execute()
	
	if err != nil {
		log.Println("Error creating chat session: ", err)
	}
	log.Println("Chat session created successfully: ", chat)

	created, err := s.GetSession(c, userID, questionID)
	if (err != nil) {
		return nil, err
	}
	log.Println("Found created session", created.ID, created.QuestionID)
	return created, nil
}

//Initial fetch of session
func (s *ChatService) GetSession(c context.Context, userID uuid.UUID, questionID uuid.UUID) (*models.ChatSession, error) {
	sessions := []models.ChatSession{}
	//Return slice of rows which matches userId and questionId
	_, err := s.supabase.From("chat_sessions").Select("*", "", false).Eq("user_id", userID.String()).Eq("question_id", questionID.String()).ExecuteTo(&sessions)
	if (err != nil) {
		return nil, err
	}

	if len(sessions) == 0 {
		return nil, nil
	}

	return &sessions[0], nil

}

// SENDING MESSAGES

func (s *ChatService) SendMessage(c context.Context, userID uuid.UUID, sessionID uuid.UUID, message string, role string) (*models.Message, error) {
	// check userID owns sessionID
	rows := []models.Row{}

	_, err := s.supabase.
		From("chat_sessions").
		Select("*", "", false).
		Eq("user_id", userID.String()).
		Eq("id", sessionID.String()).
		ExecuteTo(&rows)

	if err != nil {
		return nil, InternalServerError
	}

	if len(rows) == 0 {
		return nil, ForbiddenError
	}

	//Turn chat into interface
	chat := models.Message{
		ID: uuid.New(),
		UserID: userID,
		ChatSessionID: sessionID,
		Role: role,
		Message: message,
	}

	s.supabase.From("messages").Insert(chat, false, "", "", "").Execute()

	return &chat, nil
}