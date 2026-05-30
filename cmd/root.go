package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/therealironduck/flufu/internal/orchestrator"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "flufu",
	Short: "Wrapper for AI coding assistants that add a small coding buddy",
	Run: func(cmd *cobra.Command, _ []string) {
		orchestrator.Start(cmd.Context())
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
