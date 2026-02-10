package chat

import (
	"context"
	"github.com/google/uuid"
	"time"
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
	GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*ChatSession, error)
	IncrementMessageCount(ctx context.Context, sessionID uuid.UUID) ()
	VerifyOwnership(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) (bool, error)

	// messages
	InsertMessage(ctx context.Context, message *Message) error
	GetMessagesInitialLoad(ctx context.Context, sessionID uuid.UUID) ([]Message, error)
	GetMessagesAfterCreatedAt(ctx context.Context, sessionID uuid.UUID, createdAt time.Time) ([]Message, error)
	GetMessageByID(ctx context.Context, messageID uuid.UUID) (*Message, error)
	GetRecentMessages(ctx context.Context, sessionID uuid.UUID, limit int) ([]Message, error)

	// summaries
	UpsertSummary(ctx context.Context, summary *ConversationSummary) error
	GetSummary(ctx context.Context, sessionID uuid.UUID) (*ConversationSummary, error)
	// ListMessagesAfter(ctx context.Context, sessionID uuid.UUID, messageID uuid.UUID, createdAt time.Time) ([]Message, error)
}

type AIClient interface {
	LangChainSummarizeConversation(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, summary *ConversationSummary, messages []Message) (string, error)
	RespondToUserMessage(ctx context.Context, messages []Message, summary *ConversationSummary) (string, string, error)
}