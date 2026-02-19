package tui

import (
	"fmt"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/danmartuszewski/export-key/internal/keyitem"
)

func TestNumberInputWaitsWhenPrefixCouldMatchLongerNumber(t *testing.T) {
	m := newModel(testItems(15), "")

	m = pressRune(t, m, '1')

	if m.quitting {
		t.Fatalf("typing 1 should not auto-select when 10-15 exist")
	}
	if m.filter != "1" {
		t.Fatalf("expected filter to be 1, got %q", m.filter)
	}
}

func TestNumberInputMoveCursorToItem(t *testing.T) {
	m := newModel(testItems(15), "")

	m = pressRune(t, m, '8')

	if m.quitting {
		t.Fatalf("typing 8 should not auto-select, should move cursor instead")
	}
	// Cursor should point to item index 7 (8th item, 0-based)
	if m.filtered[m.cursor] != 7 {
		t.Fatalf("expected cursor at item index 7, got %d", m.filtered[m.cursor])
	}
}

func TestTwoDigitInputMoveCursorToItem(t *testing.T) {
	m := newModel(testItems(15), "")

	m = pressRune(t, m, '1')
	m = pressRune(t, m, '2')

	if m.quitting {
		t.Fatalf("typing 12 should not auto-select, should move cursor instead")
	}
	if m.filter != "12" {
		t.Fatalf("expected filter to be 12, got %q", m.filter)
	}
	// Cursor should point to item index 11 (12th item, 0-based)
	if m.filtered[m.cursor] != 11 {
		t.Fatalf("expected cursor at item index 11, got %d", m.filtered[m.cursor])
	}
}

func TestNumberThenEnterSelectsItem(t *testing.T) {
	m := newModel(testItems(20), "")

	m = pressRune(t, m, '2')
	m = pressRune(t, m, '0')
	m = pressKey(t, m, tea.KeyEnter)

	if !m.quitting {
		t.Fatalf("Enter after typing 20 should select the item")
	}
	if m.result.Index != 19 {
		t.Fatalf("expected index 19, got %d", m.result.Index)
	}
}

func pressRune(t *testing.T, m model, r rune) model {
	t.Helper()

	next, _ := m.Update(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{r},
	})

	updated, ok := next.(model)
	if !ok {
		t.Fatalf("expected model type %T, got %T", m, next)
	}

	return updated
}

func pressKey(t *testing.T, m model, key tea.KeyType) model {
	t.Helper()

	next, _ := m.Update(tea.KeyMsg{Type: key})

	updated, ok := next.(model)
	if !ok {
		t.Fatalf("expected model type %T, got %T", m, next)
	}

	return updated
}

func testItems(count int) []keyitem.KeyItem {
	items := make([]keyitem.KeyItem, count)
	for i := 1; i <= count; i++ {
		envVar := fmt.Sprintf("KEY_%02d", i)
		items[i-1] = keyitem.KeyItem{
			Title:   envVar + "-project",
			EnvVar:  envVar,
			Projects: []string{"project"},
		}
	}

	return items
}
