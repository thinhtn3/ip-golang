package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/supabase-community/supabase-go"
	"github.com/thinhtn3/ip-golang.git/internal/models"
)

// CONSTS
var ForbiddenError = errors.New("Forbidden: User does not own session")
var InternalServerError = errors.New("Internal server error")

// DEPENDENCY INJECTION //
type ChatService struct {
	supabase *supabase.Client
}

//constructor for ChatService
func NewChatService(supabase *supabase.Client) *ChatService {
	return &ChatService{supabase: supabase}
}

//receiver function to create chat session
func (s *ChatService) CreateSession(c context.Context, userID uuid.UUID, questionID uuid.UUID) (*models.ChatSession, error) {
	session, err := s.GetSession(c, userID, questionID)

	if (err != nil) {
		return nil, err
	}
	if (session != nil) {
		return session, nil
	}

	//create chat session object
	chat := models.ChatSession{
		ID: uuid.New(),
		UserID: userID,
		QuestionID: questionID,
		QuestionName: "", //TODO: Get question name from question_bank table
		CreatedAt: time.Now().UTC(),
	}

	//insert chat session into supabase
	_, _, err = s.supabase.
		From("chat_sessions").
		Insert(chat, true, "", "", "").
		Execute()
	
	if err != nil {
		log.Println("Error creating chat session: ", err)
		//TODO: Return error to handler
	}

	//get session after creation
	created, err := s.GetSession(c, userID, questionID)
	if (err != nil) {
		return nil, err
	}

	log.Println("Found created session", created.ID, created.QuestionID)
	return created, nil
}

// GET SESSION ID BY USER ID AND QUESTION ID //
func (s *ChatService) GetSession(c context.Context, userID uuid.UUID, questionID uuid.UUID) (*models.ChatSession, error) {
	sessions := []models.ChatSession{}
	//Return slice of rows which matches userId and questionId (because executeTo returns a slice of rows)
	_, err := s.supabase.From("chat_sessions").Select("*", "", false).Eq("user_id", userID.String()).Eq("question_id", questionID.String()).ExecuteTo(&sessions)
	if (err != nil) {
		return nil, err
	}
	
	if len(sessions) == 0 {
		return nil, nil
	}

	return &sessions[0], nil

}

// SENDING MESSAGES //
func (s *ChatService) SendMessage(c context.Context, userID uuid.UUID, sessionID uuid.UUID, message string, role string) (*models.Message, error) {
	//check userID owns sessionID
	err := s.VerifySessionOwnership(&c, userID, sessionID)
	if (err != nil) {
		return nil, err
	}

	userMessage := models.Message{
		ID: uuid.New(),
		UserID: userID,
		ChatSessionID: sessionID,
		Role: "user",
		Message: message,
	}


	//User request is the last 10 messages in the chat session for langchain
	userRequest, err := s.GetMessages(c, userID, sessionID, 10)
	if (err != nil) {
		return nil, err
	}

	//Request body is a map with key "body" and value is the user's past 10 messages
	requestBody := map[string][]models.Message{
		"body": userRequest,
	}
	body, err := json.Marshal(requestBody)

	resp, err := http.Post("http://localhost:3000/generate", "application/json", bytes.NewBuffer(body))

	if (err != nil) {
		log.Println("Error calling AI service: ", err)
	}

	defer resp.Body.Close()

	s.supabase.From("messages").Insert(userMessage, false, "", "", "").Execute()
	return &userMessage, nil
}

// GET MESSAGES //
func (s *ChatService) GetMessages(c context.Context, userID uuid.UUID, sessionID uuid.UUID, limit int) ([]models.Message, error) {
	err := s.VerifySessionOwnership(&c, userID, sessionID)
	if (err != nil) {
		return nil, err
	}
	
	chatMessages := []models.Message{}
	if limit > 0 {
		_, err = s.supabase.From("messages").Select("*", "", false).Eq("chat_session_id", sessionID.String()).Limit(limit, "").ExecuteTo(&chatMessages)
	} else {
		_, err = s.supabase.From("messages").Select("*", "", false).Eq("chat_session_id", sessionID.String()).ExecuteTo(&chatMessages)
	}

	if (err != nil) {
		return nil, err
	}
	return chatMessages, nil
}


func (s *ChatService) VerifySessionOwnership(c *context.Context, userID uuid.UUID, sessionID uuid.UUID) error {
	rows := []models.Row{}
	_, err := s.supabase.
		From("chat_sessions").
		Select("*", "", false).
		Eq("user_id", userID.String()).
		Eq("id", sessionID.String()).
		ExecuteTo(&rows)
	if (err != nil) {
		return InternalServerError
	}
	if len(rows) == 0 {
		return ForbiddenError
	}
	return nil //no error, session is owned by user
}

