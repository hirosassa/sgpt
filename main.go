package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

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
	}
	return cmd
}

func run(ctx context.Context, cmd *cli.Command) error {
	role, err := CheckGet(cmd.Bool("shell"), cmd.Bool("describe-shell"), cmd.Bool("code"))
	if err != nil {
		return err
	}
	slog.Debug("get role", slog.String("name", role.name), slog.String("role", role.role))

	chatID := cmd.String("chat")
	if chatID != "" {
		handler, err := NewChatHandler(chatID, *role)
		if err != nil {
			return fmt.Errorf("failed to create chat handler: %w", err)
		}
		res, err := handler.Handle(ctx, cmd, cmd.Args().First())
		if err != nil {
			return fmt.Errorf("failed to communicate OpenAI API: %w", err)
		}

		fmt.Println(res)
		return nil
	}

	hander, err := NewDefaultHandler(*role)
	if err != nil {
		return fmt.Errorf("failed to create chat handler: %w", err)
	}
	res, err := hander.Handle(ctx, cmd, cmd.Args().First())
	if err != nil {
		return fmt.Errorf("failed to communicate OpenAI API: %w", err)
	}

	fmt.Println(res)
	return nil
}
