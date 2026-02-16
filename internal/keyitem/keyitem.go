package keyitem

import "strings"

// KeyItem represents a secret item from any backend.
type KeyItem struct {
	// Title is the raw item title from the backend (e.g. "OPENAI_API_KEY-myapp").
	Title string
	// EnvVar is the derived environment variable name (e.g. "OPENAI_API_KEY").
	EnvVar string
	// Project is the project name extracted from the suffix (e.g. "myapp"), empty if none.
	Project string
}

// Parse creates a KeyItem from a raw title string.
// Convention: everything before the first "-" is the env var name,
// everything after is the project name.
func Parse(title string) KeyItem {
	item := KeyItem{Title: title}

	idx := strings.Index(title, "-")
	if idx == -1 {
		item.EnvVar = title
		item.Project = ""
	} else {
		item.EnvVar = title[:idx]
		item.Project = title[idx+1:]
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
		if item.Project == project {
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
