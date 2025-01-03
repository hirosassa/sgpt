package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"strings"

	"github.com/hirosassa/sgpt/handler"
	"github.com/urfave/cli/v3"
)

const (
	ExitCodeOK    int = 0
	ExitCodeError int = iota
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	os.Exit(core())
}

func core() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cmd := NewCmd()
	if err := cmd.Run(ctx, os.Args); err != nil {
		fmt.Println(err)
		return ExitCodeError
	}
	return ExitCodeOK
}

func NewCmd() *cli.Command {
	cmd := &cli.Command{
		Name:  "sgpt",
		Usage: "A command-line productivity tool powered by AI large language models (LLMs)",
		MutuallyExclusiveFlags: []cli.MutuallyExclusiveFlags{
			// following flags are mutually exclusive
			{
				Flags: [][]cli.Flag{
					{
						&cli.BoolFlag{
							Name:    "shell",
							Aliases: []string{"s"},
							Usage:   "Generate and execute shell commands.",
						},
					},
					{
						&cli.BoolFlag{
							Name:    "code",
							Aliases: []string{"c"},
							Usage:   "Generate only code.",
						},
					},
					{
						&cli.BoolFlag{
							Name:    "describe-shell",
							Aliases: []string{"d"},
							Usage:   "Describe a shell command.",
						},
					},
					{},
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "chat",
				Usage: "Follow conversation with id, \" 'use \"temp\" for quick session.",
			},
		},
		Action: run,
		// todo: try enabling this feature for stdin input.
		// ReadArgsFromStdin: true,
	}
	return cmd
}

func run(ctx context.Context, cmd *cli.Command) error {
	stat, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	var stdin []byte
	if stat.Size() > 0 { // if stdin is not empty
		stdin, err = io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read from stdin: %w", err)
		}
	}
	prompt := cmd.Args().First() + "\n" + strings.TrimSpace(string(stdin))
	slog.Debug("get prompt", slog.String("prompt", prompt))

	var h handler.Handler
	chatID := cmd.String("chat")
	switch chatID {
	case "":
		h, err = handler.NewDefaultHandler(cmd)
		if err != nil {
			return fmt.Errorf("failed to create chat handler: %w", err)
		}
	default:
		h, err = handler.NewChatHandler(cmd, chatID)
		if err != nil {
			return fmt.Errorf("failed to create chat handler: %w", err)
		}
	}

	res, err := h.Handle(ctx, cmd, prompt)
	if err != nil {
		return fmt.Errorf("failed to communicate OpenAI API: %w", err)
	}

	fmt.Println(res)
	return nil
}
