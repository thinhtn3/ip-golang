package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/supabase-community/postgrest-go"
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

// CONSTRUCTOR //
func NewChatService(supabase *supabase.Client) *ChatService {
	return &ChatService{supabase: supabase}
}

// CREATE SESSION //
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
		Archived: false,
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
	_, err := s.supabase.
		From("chat_sessions").
		Select("*", "", false).
		Eq("user_id", userID.String()).
		Eq("archived", "false").
		Eq("question_id", questionID.String()).
		ExecuteTo(&sessions)
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
		CreatedAt: time.Now().UTC(),
	}
	s.supabase.From("messages").Insert(userMessage, false, "", "", "").Execute()
	//increase message count by 1 in chat_sessions table
	s.supabase.Rpc("increment_message_count", "", map[string]interface{}{
		"session_id": sessionID.String(),
	})

	// //User request is the last 10 messages in the chat session for langchain
	// userRequest, err := s.GetMessages(c, userID, sessionID, 10)
	// if (err != nil) {
	// 	return nil, err
	// }
	// //Request body is a map with key "body" and value is the user's past 10 messages
	// requestBody := map[string][]models.Message{
	// 	"body": userRequest,
	// }
	// body, err := json.Marshal(requestBody)
	// resp, err := http.Post("http://localhost:3000/generate", "application/json", bytes.NewBuffer(body))
	// if (err != nil) {
	// 	log.Println("Error calling AI service: ", err)
	// }

	//read string response from post request body
	var aiResponse struct {
		Content string `json:"content"`
		Role string `json:"role"`
	}
	// err = json.NewDecoder(resp.Body).Decode(&aiResponse)
	// if (err != nil) {
	// 	log.Println("Error decoding AI response: ", err)
	// }
	//insert AI response into database
	aiMessage := models.Message{
		ID: uuid.New(),
		UserID: userID,
		ChatSessionID: sessionID,
		Role: aiResponse.Role,
		Message: aiResponse.Content,
		CreatedAt: time.Now().UTC(),
	}
	s.supabase.From("messages").Insert(aiMessage, false, "", "", "").Execute()
	s.supabase.Rpc("increment_message_count", "", map[string]interface{}{
		"session_id": sessionID.String(),
	})

	// Check if message count is a multiple of 10, if so, summarize conversation
	fmt.Println("Before count")
	session := []models.ChatSession{}
	_, err = s.supabase.From("chat_sessions").Select("message_count", "", false).Eq("id", sessionID.String()).ExecuteTo(&session)
	if (err != nil) {
		return nil, err
	}
	if session[0].MessageCount % 10 == 0 && session[0].MessageCount > 0 {
		s.SummarizeConversation(c, userID, sessionID)
	}
	fmt.Println("After count")
	// defer resp.Body.Close()

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
		//order by the most recent 10 in ascending order to get the most recent 10 messages (fetch message to send to langchain for context)
		_, err = s.supabase.From("messages").Select("*", "", false).Eq("chat_session_id", sessionID.String()).Order("created_at", &postgrest.OrderOpts{Ascending: false}).Limit(limit, "").ExecuteTo(&chatMessages)
		slices.Reverse(chatMessages) //reverse slice after ascending false so the latest message is at end
	} else {
		//fetch all messages for the chat session (initial load)
		_, err = s.supabase.From("messages").Select("*", "", false).Eq("chat_session_id", sessionID.String()).ExecuteTo(&chatMessages)
	}

	if (err != nil) {
		return nil, err
	}
	return chatMessages, nil
}

// VERIFY SESSION OWNERSHIP //
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

// SUMMARIZE CONVERSATION //
func (s *ChatService) SummarizeConversation(c context.Context, userID uuid.UUID, sessionID uuid.UUID) (*models.ConversationSummary, error) {
	summary, err := s.GetSummary(c, userID, sessionID)
	if (err != nil) {
		return nil, err
	}

	// If no summary found, create a new one
	if summary == nil {
		summary = &models.ConversationSummary{
			ID: uuid.New(),
			ChatSessionID: sessionID,
			Content: "No summary found",
			UpdatedAt: time.Now().UTC(),
			LastMessageID: uuid.Nil,
		}
	}

	var createdAt time.Time
	lastMessages := []models.Message{}

	if summary.LastMessageID != uuid.Nil {
		// If this is a new summary, get the first message and retrieve created_at Time to store
		_, err = s.supabase.From("messages").Select("*", "", false).Eq("id", summary.LastMessageID.String()).ExecuteTo(&lastMessages)
		if (err != nil) {
			return nil, err
		}
		createdAt = lastMessages[0].CreatedAt
	}

	// get the 10 messages AFTER the createdAt time
	messages := []models.Message{}
	_, err = s.supabase.From("messages").Select("*", "", false).Eq("chat_session_id", sessionID.String()).Order("created_at", &postgrest.OrderOpts{Ascending: true}).Gte("created_at", createdAt.Format(time.RFC3339)).Limit(10, "").ExecuteTo(&messages)
	if (err != nil) {
		return nil, err
	}

	// send summary and messages slice to localhost 300 summarize post
	requestBody := map[string]interface{}{
		"summary": summary.Content,
		"messages": messages,
	}

	body, err := json.Marshal(requestBody)
	resp, err := http.Post("http://localhost:3000/summarize", "application/json", bytes.NewBuffer(body))
	if (err != nil) {
		return nil, err
	}

	var summaryResponse models.SummaryResponse
	err = json.NewDecoder(resp.Body).Decode(&summaryResponse)
	if (err != nil) {
		fmt.Println("Error decoding summary response: ", err)
		return nil, err
	}

	//update summary with new summary
	summary.Content = summaryResponse.Content
	summary.UpdatedAt = time.Now().UTC()
	summary.LastMessageID = messages[len(messages)-1].ID
	
	// upsert summary into database
	_, _, err = s.supabase.From("conversation_summaries").Upsert(summary, "", "", "").Eq("id", summary.ID.String()).Execute()
	if (err != nil) {
		return nil, err
	}

	defer resp.Body.Close()
	return summary, nil
}

func (s *ChatService) GetSummary(c context.Context, userID uuid.UUID, sessionID uuid.UUID) (*models.ConversationSummary, error) {
	summaries := []models.ConversationSummary{}
	_, err := s.supabase.From("conversation_summaries").Select("*", "", false).Eq("chat_session_id", sessionID.String()).ExecuteTo(&summaries)
	if (err != nil) {
		return nil, err
	}
	if len(summaries) == 0 {
		return nil, nil
	}
	return &summaries[0], nil
}