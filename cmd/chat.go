package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/therealironduck/flufu/internal/ai"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "A simple chat to test out the local LLM features",
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()

		aiInstance := ai.New()

		go func() {
			if err := aiInstance.Init(ctx); err != nil {
				fmt.Fprintf(os.Stderr, "failed to init ai: %v\n", err)
				cancel()
			}
		}()

		fmt.Println("Loading model...")
		select {
		case <-aiInstance.Ready():
			fmt.Println("\n Ready!")
		case <-ctx.Done():
			return ctx.Err()
		}

		return runChatLoop(ctx, aiInstance)
	},
}

func runChatLoop(ctx context.Context, aiInstance *ai.AiInstance) error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Type your message (Ctrl+D to quit):")

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		prompt := aiInstance.MakePrompt(input)
		response, err := aiInstance.Generate(prompt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "generation error: %v\n", err)
			continue
		}

		fmt.Printf("\n%s\n\n", strings.TrimSpace(response))
	}

	return scanner.Err()
}

func init() {
	rootCmd.AddCommand(chatCmd)
}
