package chatAIClient

import (
	"bytes"
	"context"
	"net/http"
	"encoding/json"

	"github.com/google/uuid"
	chatDomain "github.com/thinhtn3/ip-golang.git/internal/domain/chat"
	chatRepo "github.com/thinhtn3/ip-golang.git/internal/infra/supabase/chat"
)

type ChatAIClient struct {
	chatRepo *chatRepo.ChatRepo
}

func NewChatAIClient(chatRepo *chatRepo.ChatRepo) *ChatAIClient {
	return &ChatAIClient{chatRepo: chatRepo}
}

func (c *ChatAIClient) LangChainSummarizeConversation(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, summary *chatDomain.ConversationSummary, messages []chatDomain.Message) (string, error) {

	requestBody := map[string]interface{}{
		"summary": summary.Content,
		"messages": messages,
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}
	resp, err := http.Post("http://localhost:3000/summarize", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	type SummaryResponse struct {
		Content string `json:"content"`
	}
	var summaryResponse SummaryResponse
	err = json.NewDecoder(resp.Body).Decode(&summaryResponse)
	if err != nil {
		return "", err
	}
	return summaryResponse.Content, nil
}

func (c *ChatAIClient) RespondToUserMessage(ctx context.Context, messages []chatDomain.Message, summary *chatDomain.ConversationSummary) (string, string, error) {
	requestBody := map[string]interface{}{
		"summary": summary.Content,
		"messages": messages,
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		return "", "", err
	}
	resp, err := http.Post("http://localhost:3000/generate", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	type Response struct {
		Content string `json:"content"`
		Role string `json:"role"`
	}
	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", "", err
	}
	return response.Content, response.Role, nil
}