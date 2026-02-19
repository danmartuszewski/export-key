package keyitem

import (
	"slices"
	"strings"
)

// KeyItem represents a secret item from any backend.
type KeyItem struct {
	// Title is the raw item title from the backend (e.g. "OPENAI_API_KEY-myapp").
	Title string
	// EnvVar is the derived environment variable name (e.g. "OPENAI_API_KEY").
	EnvVar string
	// Projects are the project names extracted from the suffix, comma-separated.
	// e.g. "myapp" or ["web", "api"] from title "KEY-web,api".
	Projects []string
}

// HasProject reports whether the item is assigned to at least one project.
func (k KeyItem) HasProject() bool {
	return len(k.Projects) > 0
}

// ProjectString returns the projects joined by comma for display.
func (k KeyItem) ProjectString() string {
	return strings.Join(k.Projects, ",")
}

// Parse creates a KeyItem from a raw title string.
// Convention: everything before the first "-" is the env var name,
// everything after is comma-separated project names.
func Parse(title string) KeyItem {
	item := KeyItem{Title: title}

	idx := strings.Index(title, "-")
	if idx == -1 {
		item.EnvVar = title
		return item
	}

	item.EnvVar = title[:idx]
	raw := title[idx+1:]
	if raw == "" {
		return item
	}

	parts := strings.Split(raw, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			item.Projects = append(item.Projects, p)
		}
	}

	return item
}

// ParseAll creates KeyItems from a slice of raw titles.
func ParseAll(titles []string) []KeyItem {
	items := make([]KeyItem, len(titles))
	for i, t := range titles {
		items[i] = Parse(t)
	}
	return items
}

// FilterByProject returns items that belong to the given project.
func FilterByProject(items []KeyItem, project string) []KeyItem {
	var filtered []KeyItem
	for _, item := range items {
		if project == "" && !item.HasProject() {
			filtered = append(filtered, item)
		} else if slices.Contains(item.Projects, project) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// EnvVarNames returns a slice of env var names from the items.
func EnvVarNames(items []KeyItem) []string {
	names := make([]string, len(items))
	for i, item := range items {
		names[i] = item.EnvVar
	}
	return names
}
