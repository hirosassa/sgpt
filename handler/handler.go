package handler

import (
	"context"
	"errors"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/urfave/cli/v3"
)

type Handler interface {
	Handle(ctx context.Context, cmd *cli.Command, prompt string) (string, error)
}

func getClient() (*openai.Client, error) {
	apiKey := os.Getenv("SGPT_OPENAI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("please set api key to SPGT_OPENAI_API_KEY")
	}

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	return client, nil
}
