package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/danmartuszewski/export-key/internal/backend"
	"github.com/danmartuszewski/export-key/internal/config"
	"github.com/danmartuszewski/export-key/internal/keyitem"
	"github.com/danmartuszewski/export-key/internal/tui"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available keys",
	RunE:  runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	tui.EnsureStderrRenderer()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	b, err := backend.New(cfg)
	if err != nil {
		return err
	}

	titles, err := b.ListItems()
	if err != nil {
		return err
	}

	items := keyitem.ParseAll(titles)

	// Calculate column widths
	maxNum := len(fmt.Sprintf("%d", len(items)))
	maxEnv := 0
	maxTitle := 0
	for _, item := range items {
		if len(item.EnvVar) > maxEnv {
			maxEnv = len(item.EnvVar)
		}
		if len(item.Title) > maxTitle {
			maxTitle = len(item.Title)
		}
	}

	// Header
	numPad := fmt.Sprintf("%-*s", maxNum+1, "#")
	envPad := fmt.Sprintf("%-*s", maxEnv, "ENV VAR")
	titlePad := fmt.Sprintf("%-*s", maxTitle, "TITLE")
	fmt.Fprintf(os.Stderr, "  %s  %s  %s  %s\n",
		tui.ListHeaderStyle.Render(numPad),
		tui.ListHeaderStyle.Render(envPad),
		tui.ListHeaderStyle.Render(titlePad),
		tui.ListHeaderStyle.Render("PROJECT"),
	)

	// Separator
	sep := strings.Repeat("─", maxNum+1+maxEnv+maxTitle+20)
	fmt.Fprintf(os.Stderr, "  %s\n", tui.ListSeparatorStyle.Render(sep))

	// Items
	for i, item := range items {
		num := fmt.Sprintf("%-*s", maxNum+1, fmt.Sprintf("%d)", i+1))
		env := fmt.Sprintf("%-*s", maxEnv, item.EnvVar)

		// Split title: base in normal color, project suffix dimmer
		var titleRendered string
		if item.Project != "" {
			base := item.EnvVar
			suffix := "-" + item.Project
			padLen := maxTitle - len(item.Title)
			pad := ""
			if padLen > 0 {
				pad = strings.Repeat(" ", padLen)
			}
			titleRendered = tui.ListTitleStyle.Render(base) + tui.ListTitleSuffixStyle.Render(suffix) + pad
		} else {
			titleRendered = tui.ListTitleStyle.Render(fmt.Sprintf("%-*s", maxTitle, item.Title))
		}

		projectCol := ""
		if item.Project != "" {
			projectCol = tui.ListProjectTagStyle.Render(item.Project)
		}

		fmt.Fprintf(os.Stderr, "  %s  %s  %s  %s\n",
			tui.ListNumberStyle.Render(num),
			tui.ListEnvVarStyle.Render(env),
			titleRendered,
			projectCol,
		)
	}

	fmt.Fprintln(os.Stderr)
	return nil
}
