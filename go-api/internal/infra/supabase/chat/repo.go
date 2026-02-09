package chatRepo


import (
	"context"
	"github.com/google/uuid"
	chatDomain "github.com/thinhtn3/ip-golang.git/internal/domain/chat"
	"github.com/supabase-community/supabase-go"
	"github.com/supabase-community/postgrest-go"
	"slices"
	"time"
)

type ChatRepo struct {
	supabase *supabase.Client
}

func NewChatRepo(supabase *supabase.Client) *ChatRepo {
	return &ChatRepo{supabase: supabase}
}

// ===============================
// SESSIONS 
// ===============================
func (r *ChatRepo) GetSessionByUserQuestion(ctx context.Context, userID uuid.UUID, questionID uuid.UUID) (*chatDomain.ChatSession, error) {
	sessions := []chatDomain.ChatSession{}
	//Return slice of rows which matches userId and questionId (because executeTo returns a slice of rows)
	_, err := r.supabase.
		From("chat_sessions").
		Select("*", "", false).
		Eq("user_id", userID.String()).
		Eq("archived", "false").
		Eq("question_id", questionID.String()).
		ExecuteTo(&sessions)
	if err != nil {
		return nil, err
	}
	if len(sessions) == 0 {
		return nil, nil
	}

	return &sessions[0], nil
}

func (r *ChatRepo) InsertSession(ctx context.Context, session *chatDomain.ChatSession) error {
	_, _, err := r.supabase.
		From("chat_sessions").
		Insert(session, true, "", "", "").
		Execute()
	if err != nil {
		return err
	}
	return nil
}

func (r *ChatRepo) GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*chatDomain.ChatSession, error) {
	session := []chatDomain.ChatSession{}
	_, err := r.supabase.
		From("chat_sessions").
		Select("*", "", false).
		Eq("id", sessionID.String()).
		ExecuteTo(&session)
	if err != nil {
		return nil, err
	}
	if len(session) == 0 {
		return nil, nil
	}
	return &session[0], nil
}

func (r *ChatRepo) IncrementMessageCount(ctx context.Context, sessionID uuid.UUID) {
	r.supabase.Rpc("increment_message_count", "", map[string]interface{}{
		"session_id": sessionID.String(),
	})
}
func (r *ChatRepo) VerifyOwnership(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) (bool, error) {
	type Row struct {
		ID uuid.UUID `json:"id"`
	}
	rows := []Row{}
	_, err := r.supabase.
		From("chat_sessions").
		Select("*", "", false).
		Eq("user_id", userID.String()).
		Eq("id", sessionID.String()).
		ExecuteTo(&rows)
	if err != nil {
		return false, err
	}
	return len(rows) > 0, nil
}

// ===============================
// MESSAGES
// ===============================
func (r *ChatRepo) InsertMessage(ctx context.Context, message *chatDomain.Message) error {
	_, _, err := r.supabase.
		From("messages").
		Insert(message, true, "", "", "").
		Execute()
	if err != nil {
		return err
	}
	return nil
}

func (r *ChatRepo) GetMessagesInitialLoad(ctx context.Context, sessionID uuid.UUID) ([]chatDomain.Message, error) {
	messages := []chatDomain.Message{}
	_, err := r.supabase.
		From("messages").
		Select("*", "", false).
		Eq("chat_session_id", sessionID.String()).
		ExecuteTo(&messages)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *ChatRepo) GetMessagesAfterCreatedAt(ctx context.Context, sessionID uuid.UUID, createdAt time.Time) ([]chatDomain.Message, error) {
	messages := []chatDomain.Message{}
	_, err := r.supabase.
		From("messages").
		Select("*", "", false).
		Eq("chat_session_id", sessionID.String()).
		Gte("created_at", createdAt.Format(time.RFC3339)).
		Order("created_at", &postgrest.OrderOpts{Ascending: false}).
		Limit(10, "").
		ExecuteTo(&messages)
	if err != nil {
		return nil, err
	}
	slices.Reverse(messages) //reverse slice after ascending false so the latest message is at end
	return messages, nil
}

func (r *ChatRepo) GetMessageByID(ctx context.Context, messageID uuid.UUID) (*chatDomain.Message, error) {
	message := chatDomain.Message{}
	_, err := r.supabase.
		From("messages").
		Select("*", "", false).
		Eq("id", messageID.String()).
		ExecuteTo(&message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}


// ===============================
// SUMMARIES
// ===============================
func (r *ChatRepo) UpsertSummary(ctx context.Context, summary *chatDomain.ConversationSummary) error {
	_, _, err := r.supabase.
		From("conversation_summaries").
		Upsert(summary, "", "", "").
		Eq("id", summary.ID.String()).
		Execute()
	if err != nil {
		return err
	}
	return nil
}

func (r *ChatRepo) GetSummary(ctx context.Context, sessionID uuid.UUID) (*chatDomain.ConversationSummary, error) {
	summaries := []chatDomain.ConversationSummary{}
	_, err := r.supabase.
		From("conversation_summaries").
		Select("*", "", false).
		Eq("chat_session_id", sessionID.String()).
		ExecuteTo(&summaries)
	if err != nil {
		return nil, err
	}
	if len(summaries) == 0 {
		return nil, nil
	}
	return &summaries[0], nil
}

// func (r *ChatRepo) ListMessagesAfterTime(ctx context.Context, sessionID uuid.UUID, messageID uuid.UUID, createdAt time.Time) ([]chatDomain.Message, error) {
// 	messages := []chatDomain.Message{}
// 	_, err := r.supabase.
// 		From("messages").
// 		Select("*", "", false).
// 		Eq("chat_session_id", sessionID.String()).
// 		Order("created_at", &postgrest.OrderOpts{Ascending: true}).
// 		Gte("created_at", createdAt.Format(time.RFC3339)).
// 		Limit(10, "").
// 		ExecuteTo(&messages)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return messages, nil
// }