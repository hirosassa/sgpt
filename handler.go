package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/urfave/cli/v2"
)

type Handler interface {
	Handle(ctx context.Context, prompt string) (string, error)
	MakeMessages(ctx context.Context, prompt string) ([]openai.ChatCompletionMessageParamUnion, error)
}

type ChatHandler struct {
	client *openai.Client
	role   SystemRole
}

func NewChatHandler(role SystemRole) (*ChatHandler, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}
	return &ChatHandler{
		client: client,
		role:   role,
	}, nil
}

func getClient() (*openai.Client, error) {
	apiKey := os.Getenv("SGPT_OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("please set api key to SPGT_OPENAI_API_KEY")
	}

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	return client, nil
}

func (h *ChatHandler) getCompletion(ctx *cli.Context, messages []openai.ChatCompletionMessageParamUnion) (string, error) {
	chatCompletion, err := h.client.Chat.Completions.New(ctx.Context, openai.ChatCompletionNewParams{
		Messages: openai.F(messages),
		Model:    openai.F(openai.ChatModelGPT4o), // todo: make this configurable
	})
	if err != nil {
		log.Println(err)
		return "", err
	}
	return chatCompletion.Choices[0].Message.Content, nil
}

func (h *ChatHandler) MakeMessages(ctx *cli.Context, prompt string) ([]openai.ChatCompletionMessageParamUnion, error) {
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(h.role.role),
		openai.UserMessage(strings.TrimSpace(prompt)),
	}
	return messages, nil
}

func (h *ChatHandler) Handle(ctx *cli.Context, prompt string) (string, error) {
	messages, err := h.MakeMessages(ctx, prompt)
	if err != nil {
		return "", err
	}
	return h.getCompletion(ctx, messages)
}

type DefaultHandler struct {
	client *openai.Client
	role   SystemRole
}

func NewDefaultHandler(role SystemRole) (*DefaultHandler, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}
	return &DefaultHandler{
		client: client,
		role:   role,
	}, nil
}

func (h *DefaultHandler) getCompletion(ctx *cli.Context, messages []openai.ChatCompletionMessageParamUnion) (string, error) {
	chatCompletion, err := h.client.Chat.Completions.New(ctx.Context, openai.ChatCompletionNewParams{
		Messages: openai.F(messages),
		Model:    openai.F(openai.ChatModelGPT4o), // todo: make this configurable
	})
	if err != nil {
		log.Println(err)
		return "", err
	}
	return chatCompletion.Choices[0].Message.Content, nil
}

func (h *DefaultHandler) MakeMessages(ctx *cli.Context, prompt string) ([]openai.ChatCompletionMessageParamUnion, error) {
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(h.role.role),
		openai.UserMessage(strings.TrimSpace(prompt)),
	}
	return messages, nil
}

func (h *DefaultHandler) Handle(ctx *cli.Context, prompt string) (string, error) {
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage(strings.TrimSpace(prompt)),
	}
	return h.getCompletion(ctx, messages)
}
