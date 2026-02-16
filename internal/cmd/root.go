package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "export-key",
	Short:         "Export API keys from 1Password or dotenv as environment variables",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	// Send help and usage output to stderr so it doesn't get eval'd
	// by the shell wrapper function.
	rootCmd.SetOut(os.Stderr)
	rootCmd.SetErr(os.Stderr)
}

func Execute() error {
	return rootCmd.Execute()
}
