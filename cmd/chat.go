package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/therealironduck/flufu/internal/ai"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "A simple chat to test out the local LLM features",
	RunE: func(cmd *cobra.Command, _ []string) error {
		err := ai.Init(cmd.Context())
		if err != nil {
			return fmt.Errorf("failed to init ai: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
}
