package handler

import (
	"context"
	"log"
	"strings"

	sgptrole "github.com/hirosassa/sgpt/role"
	"github.com/openai/openai-go"
	"github.com/urfave/cli/v3"
)

type DefaultHandler struct {
	client *openai.Client
	role   sgptrole.SystemRole
}

func NewDefaultHandler(cmd *cli.Command) (*DefaultHandler, error) {
	role, err := sgptrole.CheckGet(cmd.Bool("shell"), cmd.Bool("describe-shell"), cmd.Bool("code"))
	if err != nil {
		return nil, err
	}

	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return &DefaultHandler{
		client: client,
		role:   *role,
	}, nil
}

func (h *DefaultHandler) getCompletion(ctx context.Context, params openai.ChatCompletionNewParams) (openai.ChatCompletionMessage, error) {
	chatCompletion, err := h.client.Chat.Completions.New(ctx, params)
	if err != nil {
		log.Println(err)
		return openai.ChatCompletionMessage{}, err
	}
	return chatCompletion.Choices[0].Message, nil
}

func (h *DefaultHandler) makeParams(prompt string) openai.ChatCompletionNewParams {
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(h.role.Role),
		openai.UserMessage(strings.TrimSpace(prompt)),
	}

	return openai.ChatCompletionNewParams{
		Messages: openai.F(messages),
		Model:    openai.F(openai.ChatModelGPT4o), // todo: make this configurable
	}
}

func (h *DefaultHandler) Handle(ctx context.Context, cmd *cli.Command, prompt string) (string, error) {
	params := h.makeParams(prompt)
	message, err := h.getCompletion(ctx, params)
	if err != nil {
		return "", err
	}
	return message.Content, nil
}
