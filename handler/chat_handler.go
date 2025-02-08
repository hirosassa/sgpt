package handler

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"os"
	"strings"

	sgptrole "github.com/hirosassa/sgpt/role"
	"github.com/openai/openai-go"
	"github.com/urfave/cli/v3"
)

const cacheUmask = 0o700

// The ChatSession caches chat messages and keeps track of the conversation history.
// It is designed to store cached messages in a specified directory and in JSON format.
type ChatSession struct {
	storagePath string
}

func NewChatSession(storagePath string) (*ChatSession, error) {
	if err := createDirectory(storagePath); err != nil {
		return nil, err
	}
	slog.Debug("cache directory created", slog.String("storagePath", storagePath))

	return &ChatSession{
		storagePath: storagePath,
	}, nil
}

func (c *ChatSession) Wrap(fn func(ctx context.Context, cmd *cli.Command, params openai.ChatCompletionNewParams) (openai.ChatCompletionMessage, error)) func(ctx context.Context, cmd *cli.Command, params openai.ChatCompletionNewParams) (openai.ChatCompletionMessage, error) {
	return func(ctx context.Context, cmd *cli.Command, params openai.ChatCompletionNewParams) (openai.ChatCompletionMessage, error) {
		chatID := cmd.String("chat")
		if chatID == "" {
			return fn(ctx, cmd, params)
		}

		previousParams, err := c.read(chatID)
		if err != nil {
			return openai.ChatCompletionMessage{}, err
		}

		params.Messages.Value = append(previousParams.Messages.Value, params.Messages.Value...)
		message, err := fn(ctx, cmd, params)
		if err != nil {
			return openai.ChatCompletionMessage{}, err
		}

		params.Messages.Value = append(params.Messages.Value, message)
		if err := c.write(chatID, params); err != nil {
			return openai.ChatCompletionMessage{}, err
		}
		return message, nil
	}
}

// todo: These structs are used to parse the JSON data stored in the cache.
// This is a temporary solution and should be replaced with a more robust solution.
// ref: https://github.com/openai/openai-go/issues/133
type Content struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type Message struct {
	Content interface{} `json:"content"`
	Role    string      `json:"role"`
}

type Root struct {
	Messages []Message `json:"messages"`
	Model    string    `json:"model"`
}

func (c *ChatSession) read(chatID string) (openai.ChatCompletionNewParams, error) {
	filePath := c.storagePath + "/" + chatID
	if _, err := os.Stat(filePath); err != nil {
		//lint:ignore nilerr for initial kick
		return openai.ChatCompletionNewParams{}, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return openai.ChatCompletionNewParams{}, err
	}

	var cache Root
	if err := json.Unmarshal(data, &cache); err != nil {
		return openai.ChatCompletionNewParams{}, err
	}

	messages := marshalMessages(cache)
	return openai.ChatCompletionNewParams{
		Messages: openai.F(messages),
		Model:    openai.F(openai.ChatModelGPT4o),
	}, nil
}

func marshalMessages(cache Root) []openai.ChatCompletionMessageParamUnion {
	// Construct []openai.ChatCompletionMessageParamUnion manually from cache data.
	// todo: This process is needed because currently there is no way to directly convert the JSON data
	// to the openai.ChatCompletionMessageParamUnion type.
	// ref: https://github.com/openai/openai-go/issues/133
	var messages []openai.ChatCompletionMessageParamUnion
	for _, m := range cache.Messages {
		role := m.Role

		var content string
		switch parsed := m.Content.(type) {
		case string:
			content = parsed
		case []interface{}:
			for _, item := range parsed {
				if contentMap, ok := item.(map[string]interface{}); ok {
					content, ok = contentMap["text"].(string)
					if !ok {
						content = ""
					}
				}
			}
		default:
			content = ""
		}
		switch role {
		case "user":
			messages = append(messages, openai.UserMessage(content))
		case "assistant":
			messages = append(messages, openai.AssistantMessage(content))
		default:
			messages = append(messages, openai.SystemMessage(content))
		}
	}

	return messages
}

func (c *ChatSession) write(chatID string, params openai.ChatCompletionNewParams) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}

	filePath := c.storagePath + "/" + chatID
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, cacheUmask) // always overwrite
	if err != nil {
		return err
	}

	if _, err := f.Write(data); err != nil {
		return err
	}
	return nil
}

func (c *ChatSession) invalidate(chatID string) error {
	filePath := c.storagePath + "/" + chatID
	if _, err := os.Stat(filePath); err != nil {
		// lint:ignore nilerr already invalidated
		return nil
	}
	if err := os.Remove(filePath); err != nil {
		return err
	}
	return nil
}

func (c *ChatSession) exists(chatID string) bool {
	data, err := c.read(chatID)
	if err != nil {
		return false
	}
	return len(data.Messages.Value) > 0
}

// codes below are for in the future
// func (c *ChatSession) list() ([]string, error) {
// 	files, err := os.ReadDir(c.storagePath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	filePathes := make([]string, len(files))
// 	sortFilesByModTime(files)
// 	for i, file := range files {
// 		filePathes[i] = c.storagePath + "/" + file.Name()
// 	}
// 	return filePathes, nil
// }

// func sortFilesByModTime(files []os.DirEntry) {
// 	sort.Slice(files, func(i, j int) bool {
// 		return getModTime(files[i]).After(getModTime(files[j]))
// 	})
// }

// func getModTime(file os.DirEntry) time.Time {
// 	fileInfo, err := file.Info()
// 	if err != nil {
// 		return time.Time{}
// 	}
// 	return fileInfo.ModTime()
// }

func createDirectory(storagePath string) error {
	err := os.MkdirAll(storagePath, cacheUmask)
	if err != nil {
		return err
	}
	return nil
}

type ChatHandler struct {
	client      *openai.Client
	role        sgptrole.SystemRole
	chatID      string
	chatSession *ChatSession
}

func NewChatHandler(cmd *cli.Command, chatID string) (*ChatHandler, error) {
	chatCachePath := os.ExpandEnv("$HOME/.config/shell_gpt/chat_cache")
	chatSession, err := NewChatSession(chatCachePath) // todo: make this configurable
	if err != nil {
		return nil, err
	}

	role, err := sgptrole.CheckGet(cmd.Bool("shell"), cmd.Bool("describe-shell"), cmd.Bool("code"))
	if err != nil {
		return nil, err
	}

	client, err := getClient()
	if err != nil {
		return nil, err
	}

	if chatID == "temp" {
		if err := chatSession.invalidate(chatID); err != nil {
			return nil, err
		}
	}

	return &ChatHandler{
		client:      client,
		role:        *role,
		chatID:      chatID,
		chatSession: chatSession,
	}, nil
}

func (h *ChatHandler) initiated() bool {
	return h.chatSession.exists(h.chatID)
}

func (h *ChatHandler) getCompletion(ctx context.Context, cmd *cli.Command, params openai.ChatCompletionNewParams) (openai.ChatCompletionMessage, error) {
	chatCompletion, err := h.client.Chat.Completions.New(ctx, params)
	if err != nil {
		log.Println(err)
		return openai.ChatCompletionMessage{}, err
	}
	return chatCompletion.Choices[0].Message, nil
}

func (h *ChatHandler) makeParams(prompt string) openai.ChatCompletionNewParams {
	if h.initiated() {
		messages := []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(h.role.Role),
			openai.UserMessage(strings.TrimSpace(prompt)),
		}
		return openai.ChatCompletionNewParams{
			Messages: openai.F(messages),
			Model:    openai.F(openai.ChatModelGPT4o), // todo: make this configurable
		}
	}
	return openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(strings.TrimSpace(prompt)),
		}),
		Model: openai.F(openai.ChatModelGPT4o),
	}
}

func (h *ChatHandler) Handle(ctx context.Context, cmd *cli.Command, prompt string) (string, error) {
	params := h.makeParams(prompt)

	wrappedGetCompletion := h.chatSession.Wrap(h.getCompletion)
	message, err := wrappedGetCompletion(ctx, cmd, params)
	if err != nil {
		return "", err
	}
	return message.Content, nil
}
