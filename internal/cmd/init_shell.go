package cmd

import (
	"fmt"

	"github.com/danmartuszewski/export-key/internal/shell"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:       "init [shell]",
	Short:     "Output shell integration code",
	Long:      "Output shell function for the specified shell.\nAdd eval \"$(export-key init zsh)\" to your ~/.zshrc",
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"zsh", "bash", "fish"},
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "zsh":
			fmt.Print(shell.ZshInit)
		case "bash":
			fmt.Print(shell.BashInit)
		case "fish":
			fmt.Print(shell.FishInit)
		default:
			return fmt.Errorf("unsupported shell: %s (supported: zsh, bash, fish)", args[0])
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
