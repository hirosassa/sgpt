package handler

import (
	"context"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/urfave/cli/v3"
	"google.golang.org/api/option"
)

var _ Handler = (*GeminiChatHandler)(nil)

type GeminiChatHandler struct {
	// storagePath string TODO: Implement chat session caching
	client *genai.Client
	model  string
}

func NewGeminiChatHandler(ctx context.Context, apiKey string, model string) (*GeminiChatHandler, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &GeminiChatHandler{
		client: client,
		model:  model,
	}, nil
}

// Handle
func (h *GeminiChatHandler) Handle(ctx context.Context, cmd *cli.Command, prompt string) (string, error) {
	model := h.client.GenerativeModel(h.model)
	session := model.StartChat()

	response, err := session.SendMessage(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	responseMessage := []string{}
	for _, part := range response.Candidates[0].Content.Parts {
		if str, ok := part.(genai.Text); ok {
			responseMessage = append(responseMessage, string(str))
		}
	}
	return strings.Join(responseMessage, "\n"), nil
}
