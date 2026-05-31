package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/therealironduck/flufu/internal/agent"
	"github.com/therealironduck/flufu/internal/pet"
)

const msgArgCount = 2

var msgCmd = &cobra.Command{
	Use:   "msg <pid> <transcript_path>",
	Short: "Send information to the Flufu socket for the pet to display. Should not be used manually!",
	Args:  cobra.ExactArgs(msgArgCount),
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("pid must be a number: %w", err)
		}

		message, err := agent.ReadNewestUserMessageFromTranscript(args[1])
		if err != nil {
			return nil // No error
		}

		// Errors are okay here.
		_ = pet.Send(cmd.Context(), pid, message)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(msgCmd)
}
