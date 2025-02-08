package main

import (
	"log/slog"
	"os"

	"github.com/hirosassa/sgpt/cmd"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	os.Exit(cmd.Core())
}
