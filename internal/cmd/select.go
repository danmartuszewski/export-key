package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/danmartuszewski/export-key/internal/backend"
	"github.com/danmartuszewski/export-key/internal/config"
	"github.com/danmartuszewski/export-key/internal/keyitem"
	"github.com/danmartuszewski/export-key/internal/tui"
	"github.com/spf13/cobra"
)

var projectFlag string

var selectCmd = &cobra.Command{
	Use:   "select [query|number|.project]",
	Short: "Select and export an API key",
	Long:  "Select a key interactively or by query/number. Use .project to export all keys for a project.\nOutput is eval'd by the shell wrapper.",
	RunE:  runSelect,
}

func init() {
	selectCmd.Flags().StringVarP(&projectFlag, "project", "p", "", "Export all keys for a project")
	rootCmd.AddCommand(selectCmd)
}

func runSelect(cmd *cobra.Command, args []string) error {
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

	// Project mode: export all keys for a project
	if projectFlag != "" {
		return exportProject(b, items, projectFlag)
	}

	// Dot-prefix shortcut: ek .myapp -> project mode
	if len(args) == 1 && strings.HasPrefix(args[0], ".") {
		return exportProject(b, items, args[0][1:])
	}

	// Number mode: direct select by number
	if len(args) == 1 {
		if num, err := strconv.Atoi(args[0]); err == nil {
			idx := num - 1 // 1-based to 0-based
			if idx < 0 || idx >= len(items) {
				return fmt.Errorf("number %d out of range (1-%d)", num, len(items))
			}
			return exportSingle(b, items[idx])
		}
	}

	// Query mode: filter by text then show picker if needed
	if len(args) == 1 {
		query := strings.ToLower(args[0])
		var matches []keyitem.KeyItem
		for _, item := range items {
			if strings.Contains(strings.ToLower(item.EnvVar), query) ||
				strings.Contains(strings.ToLower(item.Title), query) {
				matches = append(matches, item)
			}
		}

		if len(matches) == 1 {
			return exportSingle(b, matches[0])
		}

		if len(matches) > 1 {
			items = matches
		}
		// If no matches, fall through to show all items in picker
	}

	// Interactive TUI picker
	result, err := tui.RunPicker(items, "v"+Version)
	if err != nil {
		return err
	}

	if result.Canceled {
		return nil
	}

	// Multi-select: export all space-selected items
	if len(result.Items) > 0 {
		return exportMultiple(b, result.Items)
	}

	return exportSingle(b, result.Item)
}

func exportSingle(b backend.Backend, item keyitem.KeyItem) error {
	secret, err := b.GetSecret(item.Title)
	if err != nil {
		return err
	}

	// Output the export statement to stdout for eval
	fmt.Printf("export %s=%q\n", item.EnvVar, secret)

	// Feedback on stderr so user sees it but it doesn't get eval'd
	if item.HasProject() {
		fmt.Fprintf(os.Stderr, "Exported %s [%s]\n", item.EnvVar, item.ProjectString())
	} else {
		fmt.Fprintf(os.Stderr, "Exported %s\n", item.EnvVar)
	}

	return nil
}

func exportMultiple(b backend.Backend, items []keyitem.KeyItem) error {
	var exported []string
	for _, item := range items {
		secret, err := b.GetSecret(item.Title)
		if err != nil {
			return fmt.Errorf("fetching %s: %w", item.EnvVar, err)
		}

		fmt.Printf("export %s=%q\n", item.EnvVar, secret)
		exported = append(exported, item.EnvVar)
	}

	fmt.Fprintf(os.Stderr, "Exported %d keys: %s\n",
		len(exported), strings.Join(exported, ", "))

	return nil
}

func exportProject(b backend.Backend, items []keyitem.KeyItem, project string) error {
	projectItems := keyitem.FilterByProject(items, project)
	if len(projectItems) == 0 {
		return fmt.Errorf("no keys found for project %q", project)
	}

	var exported []string
	for _, item := range projectItems {
		secret, err := b.GetSecret(item.Title)
		if err != nil {
			return fmt.Errorf("fetching %s: %w", item.EnvVar, err)
		}

		fmt.Printf("export %s=%q\n", item.EnvVar, secret)
		exported = append(exported, item.EnvVar)
	}

	fmt.Fprintf(os.Stderr, "Exported %d keys for project %s: %s\n",
		len(exported), project, strings.Join(exported, ", "))

	return nil
}
