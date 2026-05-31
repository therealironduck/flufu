package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/therealironduck/flufu/internal/pet"
)

var msgCmd = &cobra.Command{
	Use:   "msg <pid>",
	Short: "Send information to the Flufu socket for the pet to display. Should not be used manually!",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("pid must be a number: %w", err)
		}

		// Errors are okay here.
		_ = pet.Send(cmd.Context(), pid)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(msgCmd)
}
