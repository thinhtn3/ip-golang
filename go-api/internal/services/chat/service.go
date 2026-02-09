package chat

import (
	// "bytes"
	// "encoding/json"
	// "log"
	// "net/http"
	// "slices"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	chatDomain "github.com/thinhtn3/ip-golang.git/internal/domain/chat"
)



type ChatService struct {
	repo chatDomain.Repository
	aiClient chatDomain.AIClient
}

func NewChatService(repo chatDomain.Repository, aiClient chatDomain.AIClient) *ChatService {
	return &ChatService{repo: repo, aiClient: aiClient}
}

// ===============================
// SESSION SERVICES
// ===============================
func (r *ChatService) CreateSession(ctx context.Context, userID uuid.UUID, questionID uuid.UUID) (*chatDomain.ChatSession, error) {
	//check if session already exists
	session, err := r.GetSession(ctx, userID, questionID)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get session: %w", err)
	}

	if session != nil {
		return session, nil
	}

	//create chat session object
	newSession := chatDomain.ChatSession{
		ID: uuid.New(),
		UserID: userID,
		QuestionID: questionID,
		QuestionName: "", //TODO: Get question name from question_bank table
		CreatedAt: time.Now().UTC(),
		Archived: false,
	}

	//insert chat session into supabase
	err = r.repo.InsertSession(ctx, &newSession)

	if err != nil {
		return nil, fmt.Errorf("service: failed to insert session: %w", err)
	}

	//get session after creation
	created, err := r.repo.GetSessionByUserQuestion(ctx, userID, questionID)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get session by user question after creation: %w", err)
	}

	return created, nil
}

func (r *ChatService) GetSession(ctx context.Context, userID uuid.UUID, questionID uuid.UUID) (*chatDomain.ChatSession, error) {
	session, err := r.repo.GetSessionByUserQuestion(ctx, userID, questionID)
	if err != nil {
		return nil, chatDomain.ErrInternalServerError
	}

	if session == nil {
		return nil, nil //session not found
	}

	return session, nil //session found
}

// ===============================
// MESSAGES SERVICES
// ===============================
func (r *ChatService) SendMessage(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, message string, role string) (*chatDomain.Message, error) {
	//check userID owns sessionID
	owns, err := r.repo.VerifyOwnership(ctx, userID, sessionID)
	if err != nil {
		return nil, err
	}
	if !owns {
		return nil, chatDomain.ErrForbidden
	}

	//create user message object
	userMessage := chatDomain.Message{
		ID: uuid.New(),
		UserID: userID,
		ChatSessionID: sessionID,
		Role: "user",
		Message: message,
		CreatedAt: time.Now().UTC(),
	}
	fmt.Println("service: Creating user message object")
	err = r.repo.InsertMessage(ctx, &userMessage)
	if err != nil {
		return nil, fmt.Errorf("service: failed to insert message: %w", err)
	}
	fmt.Println("service: User Message inserted")
	r.repo.IncrementMessageCount(ctx, sessionID)
	fmt.Println("service: Message count incremented")


	var aiResponse struct {
		Content string `json:"content"`
		Role string `json:"role"`
	}
	aiMessage := chatDomain.Message{
		ID: uuid.New(),
		UserID: userID,
		ChatSessionID: sessionID,
		Role: aiResponse.Role,
		Message: aiResponse.Content,
		CreatedAt: time.Now().UTC(),
	}
	err = r.repo.InsertMessage(ctx, &aiMessage)
	if err != nil {
		return nil, fmt.Errorf("service: failed to insert AI message: %w", err)
	}
	r.repo.IncrementMessageCount(ctx, sessionID)
	fmt.Println("service: Message count incremented")

	// Check if message count is a multiple of 10, if so, summarize conversation
	fmt.Println("Before count")
	session, err := r.repo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get session by user question: %w", err)
	}
	messageCount := session.MessageCount
	fmt.Println("service: Message count: ", messageCount)
	if messageCount % 10 == 0 && messageCount > 0 {
		summary, err := r.SummarizeConversation(ctx, userID, sessionID)
		if err != nil {
			return nil, fmt.Errorf("service: failed to summarize conversation: %w", err)
		}
		//print summary content to console
		fmt.Println("service: Summary content: ", summary.Content)
	}

	return &userMessage, nil
}

func (r *ChatService) GetMessages(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, limit int) ([]chatDomain.Message, error) {
	owns, err := r.repo.VerifyOwnership(ctx, userID, sessionID)
	if err != nil {
		return nil, err
	}
	if !owns {
		return nil, chatDomain.ErrForbidden
	}
	
	chatMessages, err := r.repo.GetMessagesInitialLoad(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get messages initial load: %w", err)
	}
	return chatMessages, nil
}

// ===============================
// SUMMARIZATION SERVICES
// ===============================
func (r *ChatService) SummarizeConversation(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) (*chatDomain.ConversationSummary, error) {
	//1. Get existing summary from repository
	summary, err := r.repo.GetSummary(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	//2. If no summary, create a new one
	if summary == nil {
		summary = &chatDomain.ConversationSummary{
			ID: uuid.New(),
			ChatSessionID: sessionID,
			Content: "No summary found",
			UpdatedAt: time.Now().UTC(),
			LastMessageID: uuid.Nil,
		}
	}
	var createdAt time.Time;

	if summary.LastMessageID != uuid.Nil {
		// Last message ID is not nil, use the last message time stamp
		lastMessage, err := r.repo.GetMessageByID(ctx, summary.LastMessageID);
		if err != nil {
			return nil, err
		}
		createdAt = lastMessage.CreatedAt
	} else if summary.LastMessageID == uuid.Nil {
		// First time summarizing, use the first message time stamp in the session
		firstMessage, err := r.repo.GetMessagesInitialLoad(ctx, sessionID);
		if err != nil {
			return nil, err
		}
		createdAt = firstMessage[0].CreatedAt
	}

	//3. Get messages after the createdAt time stamp
	messagesToSummarize, err := r.repo.GetMessagesAfterCreatedAt(ctx, sessionID, createdAt);
	if err != nil {
		return nil, err
	}

	//4. Summarize the messages
	summary.Content, err = r.aiClient.LangChainSummarizeConversation(ctx, userID, sessionID, summary, messagesToSummarize);
	if err != nil {
		return nil, err
	}
	
	return summary, nil
}