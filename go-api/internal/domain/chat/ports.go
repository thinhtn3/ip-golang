package chat

import (
	"context"
	"github.com/google/uuid"
)

// type ChatService interface {
// 	CreateSession(c context.Context, userID uuid.UUID, questionID uuid.UUID) (*ChatSession, error)
// 	GetSession(c context.Context, userID uuid.UUID, questionID uuid.UUID) (*ChatSession, error)
// 	SendMessage(c context.Context, userID uuid.UUID, sessionID uuid.UUID, message string, role string) (*Message, error)
// 	GetMessages(c context.Context, userID uuid.UUID, sessionID uuid.UUID, limit int) ([]Message, error)
// 	VerifySessionOwnership(c context.Context, userID uuid.UUID, sessionID uuid.UUID) error
// 	SummarizeConversation(c context.Context, userID uuid.UUID, sessionID uuid.UUID) (*ConversationSummary, error)
// 	GetSummary(c context.Context, userID uuid.UUID, sessionID uuid.UUID) (*ConversationSummary, error)
// }

type Repository interface {
	// layer that handles database operations

	// sessions
	GetSessionByUserQuestion(ctx context.Context, userID uuid.UUID, questionID uuid.UUID) (*ChatSession, error)
	InsertSession(ctx context.Context, session *ChatSession) error
	GetMessageCount(ctx context.Context, sessionID uuid.UUID) (int, error)
	IncrementMessageCount(ctx context.Context, sessionID uuid.UUID) error
	VerifyOwnership(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) error

	// messages
	InsertMessage(ctx context.Context, message *Message) error
	GetMessages(ctx context.Context, sessionID uuid.UUID, limit int) ([]Message, error)

	// summaries
	UpsertSummary(ctx context.Context, summary *ConversationSummary) error
	GetSummary(ctx context.Context, sessionID uuid.UUID) (*ConversationSummary, error)
	ListMessagesAfter(ctx context.Context, sessionID uuid.UUID, messageID uuid.UUID) ([]Message, error)
}

type AIClient interface {
	SummarizeConversation(ctx context.Context, messages []Message) (string, error)
	GenerateResponse(ctx context.Context, messages []Message) (string, error)
}