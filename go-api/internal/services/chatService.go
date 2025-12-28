package services

import (
	"context"
	"log"

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
func (s *ChatService) CreateChatSession(c context.Context, userID string, questionID string) {
	var questionName string
	question, err := s.supabase.From("question_bank").Select("title").Eq("id", questionID).Single().ExecuteTo(&questionName)
	if err != nil {
		log.Println("Error getting question: ", err)
	}
	log.Println("Question: ", question)

	// chat := &models.ChatSession{
	// 	UserID: userID,
	// 	QuestionID: question.ID,
	// 	QuestionName: question.Title,
	// }

	// err = s.db.QueryRowContext(c, `insert into chat_sessions (user_id, question_id, question_name) values ($1, $2, $3, $4) returning id`, chat.UserID, chat.QuestionID, chat.QuestionName, chat.CreatedAt).Scan(&chat.ID)
	// if err != nil {
	// 	return nil, err
	// }
	// return chat, nil
}