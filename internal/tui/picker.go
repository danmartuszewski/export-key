package tui

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/danmartuszewski/export-key/internal/keyitem"
)

const pageSize = 10

// PickerResult contains what the user selected.
type PickerResult struct {
	// Item is set when a single item is selected (enter without space selections).
	Item keyitem.KeyItem
	// Items is set when multiple items are selected via space.
	Items []keyitem.KeyItem
	// Index of the single selected item.
	Index    int
	Canceled bool
}

// RunPicker shows the interactive TUI picker on stderr and returns the selection.
func RunPicker(items []keyitem.KeyItem, version string) (PickerResult, error) {
	EnsureStderrRenderer()

	m := newModel(items, version)

	p := tea.NewProgram(m,
		tea.WithOutput(os.Stderr),
		tea.WithAltScreen(),
	)

	finalModel, err := p.Run()
	if err != nil {
		return PickerResult{Canceled: true}, err
	}

	fm := finalModel.(model)
	return fm.result, nil
}

type model struct {
	items    []keyitem.KeyItem
	filtered []int // indices into items
	cursor   int
	filter   string
	version  string
	result   PickerResult
	quitting bool
	selected map[int]bool // indices of space-selected items

	maxEnvLen   int
	maxTitleLen int
	termHeight  int
}

func newModel(items []keyitem.KeyItem, version string) model {
	indices := make([]int, len(items))
	maxEnv, maxTitle := 0, 0
	for i, item := range items {
		indices[i] = i
		if len(item.EnvVar) > maxEnv {
			maxEnv = len(item.EnvVar)
		}
		if len(item.Title) > maxTitle {
			maxTitle = len(item.Title)
		}
	}

	return model{
		items:       items,
		filtered:    indices,
		version:     version,
		maxEnvLen:   maxEnv,
		maxTitleLen: maxTitle,
		selected:    make(map[int]bool),
		result:      PickerResult{Canceled: true},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termHeight = msg.Height
		return m, nil

	case tea.KeyMsg:
		key := msg.String()

		// Always handle these regardless of filter state
		switch key {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			if len(m.selected) > 0 {
				// Multi-select: export all space-selected items
				var items []keyitem.KeyItem
				for i := 0; i < len(m.items); i++ {
					if m.selected[i] {
						items = append(items, m.items[i])
					}
				}
				m.result = PickerResult{Items: items}
			} else if len(m.filtered) > 0 {
				// Single select: export cursor item
				idx := m.filtered[m.cursor]
				m.result = PickerResult{
					Item:  m.items[idx],
					Index: idx,
				}
			}
			m.quitting = true
			return m, tea.Quit

		case " ":
			// Toggle selection on current item
			if len(m.filtered) > 0 {
				idx := m.filtered[m.cursor]
				if m.selected[idx] {
					delete(m.selected, idx)
				} else {
					m.selected[idx] = true
				}
				// Move cursor down after toggling
				if m.cursor < len(m.filtered)-1 {
					m.cursor++
				}
			}
			return m, nil

		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil

		case "down":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			}
			return m, nil

		case "pgup":
			m.cursor -= pageSize
			if m.cursor < 0 {
				m.cursor = 0
			}
			return m, nil

		case "pgdown":
			m.cursor += pageSize
			if m.cursor >= len(m.filtered) {
				m.cursor = max(0, len(m.filtered)-1)
			}
			return m, nil

		case "home":
			m.cursor = 0
			return m, nil

		case "end":
			m.cursor = max(0, len(m.filtered)-1)
			return m, nil

		case "backspace":
			if len(m.filter) > 0 {
				m.filter = m.filter[:len(m.filter)-1]
				m.applyFilter()
			}
			return m, nil
		}

		// Vim-style navigation only when not filtering
		if m.filter == "" {
			switch key {
			case "q":
				m.quitting = true
				return m, tea.Quit
			case "k":
				if m.cursor > 0 {
					m.cursor--
				}
				return m, nil
			case "j":
				if m.cursor < len(m.filtered)-1 {
					m.cursor++
				}
				return m, nil
			case "g":
				m.cursor = 0
				return m, nil
			case "G":
				m.cursor = max(0, len(m.filtered)-1)
				return m, nil
			}
		}

		// Everything else: add to filter
		if len(key) == 1 {
			m.filter += key

			// If the filter is a number, jump the cursor to that item
			// instead of fuzzy filtering. User confirms with Enter.
			if num, err := strconv.Atoi(m.filter); err == nil {
				idx := num - 1 // user sees 1-based numbers
				if idx >= 0 && idx < len(m.items) {
					// Find this item's position in the filtered list
					for i, fi := range m.filtered {
						if fi == idx {
							m.cursor = i
							break
						}
					}
				}
				return m, nil
			}

			m.applyFilter()
		}
	}

	return m, nil
}

func (m *model) applyFilter() {
	if m.filter == "" {
		m.filtered = make([]int, len(m.items))
		for i := range m.items {
			m.filtered[i] = i
		}
		m.cursor = 0
		return
	}

	lower := strings.ToLower(m.filter)
	m.filtered = nil
	for i, item := range m.items {
		if fuzzyMatch(lower, strings.ToLower(item.Title)) ||
			fuzzyMatch(lower, strings.ToLower(item.EnvVar)) ||
			fuzzyMatch(lower, strings.ToLower(item.Project)) {
			m.filtered = append(m.filtered, i)
		}
	}

	if m.cursor >= len(m.filtered) {
		m.cursor = max(0, len(m.filtered)-1)
	}
}

func fuzzyMatch(pattern, text string) bool {
	if pattern == "" {
		return true
	}
	return strings.Contains(text, pattern)
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	// Header
	b.WriteString("\n  ")
	b.WriteString(titleStyle.Render("export-key"))
	if m.version != "" {
		b.WriteString("  ")
		b.WriteString(versionStyle.Render(m.version))
	}
	if len(m.selected) > 0 {
		b.WriteString("  ")
		b.WriteString(selectedCountStyle.Render(fmt.Sprintf("%d selected", len(m.selected))))
	}
	b.WriteString("\n")

	// Filter prompt
	b.WriteString("  ")
	if m.filter != "" {
		b.WriteString(promptStyle.Render("/ "))
		b.WriteString(filterInputStyle.Render(m.filter))
	} else {
		b.WriteString(promptStyle.Render("/ Type to filter (or enter a number)"))
	}
	b.WriteString("\n\n")

	// Items
	for i, idx := range m.filtered {
		item := m.items[idx]
		isCursor := i == m.cursor
		isChecked := m.selected[idx]

		// Cursor
		if isCursor {
			b.WriteString(cursorStyle.Render(" > "))
		} else {
			b.WriteString("   ")
		}

		// Checkbox
		if isChecked {
			b.WriteString(checkStyle.Render("[x] "))
		} else {
			b.WriteString(uncheckStyle.Render("[ ] "))
		}

		// Number (1-based)
		num := fmt.Sprintf("%d)", idx+1)
		if isCursor {
			b.WriteString(selectedNumberStyle.Render(num))
		} else {
			b.WriteString(numberStyle.Render(num))
		}
		b.WriteString(" ")

		// Env var name
		envPad := fmt.Sprintf("%-*s", m.maxEnvLen, item.EnvVar)
		if isCursor {
			b.WriteString(selectedEnvVarStyle.Render(envPad))
		} else {
			b.WriteString(envVarStyle.Render(envPad))
		}
		b.WriteString("  ")

		// Full title (dim the project suffix after the first -)
		if item.Project != "" {
			base := item.EnvVar
			suffix := "-" + item.Project
			padLen := m.maxTitleLen - len(item.Title)
			pad := ""
			if padLen > 0 {
				pad = strings.Repeat(" ", padLen)
			}
			if isCursor {
				b.WriteString(selectedTitleColumnStyle.Render(base))
				b.WriteString(selectedTitleSuffixStyle.Render(suffix))
				b.WriteString(pad)
			} else {
				b.WriteString(titleColumnStyle.Render(base))
				b.WriteString(titleSuffixStyle.Render(suffix))
				b.WriteString(pad)
			}
		} else {
			titlePad := fmt.Sprintf("%-*s", m.maxTitleLen, item.Title)
			if isCursor {
				b.WriteString(selectedTitleColumnStyle.Render(titlePad))
			} else {
				b.WriteString(titleColumnStyle.Render(titlePad))
			}
		}

		// Project tag
		if item.Project != "" {
			tag := fmt.Sprintf("  [%s]", item.Project)
			if isCursor {
				b.WriteString(selectedProjectStyle.Render(tag))
			} else {
				b.WriteString(projectStyle.Render(tag))
			}
		}

		b.WriteString("\n")
	}

	if len(m.filtered) == 0 {
		b.WriteString("   ")
		b.WriteString(promptStyle.Render("No matches"))
		b.WriteString("\n")
	}

	// Help footer
	b.WriteString("\n  ")
	b.WriteString(helpStyle.Render("↑/k up  ↓/j down  g/G top/bottom  space select  enter export  q quit"))
	b.WriteString("\n")

	content := b.String()

	// Bottom-align: pad with empty lines so content sits at the bottom
	// of the terminal instead of the top.
	if m.termHeight > 0 {
		contentLines := strings.Count(content, "\n")
		padding := m.termHeight - contentLines
		if padding > 0 {
			content = strings.Repeat("\n", padding) + content
		}
	}

	return content
}
