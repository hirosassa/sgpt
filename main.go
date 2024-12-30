package main

import (
	"os"
	"context"
	"io"
	"log"
	"flag"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

const cmdName = "sgpt"

type CLI struct {
	Stdin          io.Reader
	Stdout, Stderr io.Writer
}

func main() {
	cli := &CLI{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	flag.Parse()
	if err := cli.run(flag.Args()); err != nil {
		fmt.Fprintln(cli.Stderr, err)
		os.Exit(1)
	}
}

func (c *CLI) run(args []string) error {
	apiKey :=os.Getenv("SGPT_OPENAI_API_KEY")
	if apiKey == "" {
		log.Println("please set api key to SPGT_OPENAI_API_KEY")
		return fmt.Errorf("please set api key to SPGT_OPENAI_API_KEY")
	}

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	prompt := args[0]
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			 openai.UserMessage(prompt),
		}),
		Model: openai.F(openai.ChatModelGPT4o),
	})
	if err != nil {
		log.Println(err)
		return err
	}

	fmt.Fprintln(c.Stdout, chatCompletion.Choices[0].Message.Content)
	return nil
}
