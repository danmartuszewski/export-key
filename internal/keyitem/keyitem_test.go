package keyitem

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		title   string
		envVar  string
		project string
	}{
		{"OPENAI_API_KEY", "OPENAI_API_KEY", ""},
		{"OPENAI_API_KEY-myapp", "OPENAI_API_KEY", "myapp"},
		{"GEMINI_API_KEY-projectA", "GEMINI_API_KEY", "projectA"},
		{"STRIPE_KEY-my-app", "STRIPE_KEY", "my-app"}, // second dash is part of project
		{"KEY-", "KEY", ""},         // trailing dash, empty project
		{"-project", "", "project"}, // leading dash, empty env var
		{"SIMPLE", "SIMPLE", ""},    // no dash at all
		{"A-B", "A", "B"},           // minimal
		{"", "", ""},                // empty string
		{"NO_DASH_HERE", "NO_DASH_HERE", ""},
		{"AWS_ACCESS_KEY-infra", "AWS_ACCESS_KEY", "infra"},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			item := Parse(tt.title)

			if item.Title != tt.title {
				t.Errorf("Title = %q, want %q", item.Title, tt.title)
			}
			if item.EnvVar != tt.envVar {
				t.Errorf("Parse(%q).EnvVar = %q, want %q", tt.title, item.EnvVar, tt.envVar)
			}
			if item.Project != tt.project {
				t.Errorf("Parse(%q).Project = %q, want %q", tt.title, item.Project, tt.project)
			}
		})
	}
}

func TestParseAll(t *testing.T) {
	titles := []string{"OPENAI_API_KEY-myapp", "AWS_KEY", "STRIPE-prod"}
	items := ParseAll(titles)

	if len(items) != 3 {
		t.Fatalf("ParseAll returned %d items, want 3", len(items))
	}

	if items[0].EnvVar != "OPENAI_API_KEY" || items[0].Project != "myapp" {
		t.Errorf("items[0] = %+v", items[0])
	}
	if items[1].EnvVar != "AWS_KEY" || items[1].Project != "" {
		t.Errorf("items[1] = %+v", items[1])
	}
	if items[2].EnvVar != "STRIPE" || items[2].Project != "prod" {
		t.Errorf("items[2] = %+v", items[2])
	}
}

func TestParseAllEmpty(t *testing.T) {
	items := ParseAll(nil)
	if len(items) != 0 {
		t.Errorf("ParseAll(nil) returned %d items, want 0", len(items))
	}
}

func TestFilterByProject(t *testing.T) {
	items := ParseAll([]string{
		"OPENAI_API_KEY-myapp",
		"STRIPE_KEY-myapp",
		"AWS_KEY",
		"GEMINI_API_KEY-infra",
		"DB_URL-myapp",
	})

	t.Run("existing project", func(t *testing.T) {
		filtered := FilterByProject(items, "myapp")
		if len(filtered) != 3 {
			t.Fatalf("FilterByProject(myapp) returned %d items, want 3", len(filtered))
		}
		envs := EnvVarNames(filtered)
		want := []string{"OPENAI_API_KEY", "STRIPE_KEY", "DB_URL"}
		for i, w := range want {
			if envs[i] != w {
				t.Errorf("envs[%d] = %q, want %q", i, envs[i], w)
			}
		}
	})

	t.Run("single match", func(t *testing.T) {
		filtered := FilterByProject(items, "infra")
		if len(filtered) != 1 {
			t.Fatalf("FilterByProject(infra) returned %d items, want 1", len(filtered))
		}
		if filtered[0].EnvVar != "GEMINI_API_KEY" {
			t.Errorf("got %q, want GEMINI_API_KEY", filtered[0].EnvVar)
		}
	})

	t.Run("no match", func(t *testing.T) {
		filtered := FilterByProject(items, "nonexistent")
		if len(filtered) != 0 {
			t.Errorf("FilterByProject(nonexistent) returned %d items, want 0", len(filtered))
		}
	})

	t.Run("empty project matches no-project items", func(t *testing.T) {
		filtered := FilterByProject(items, "")
		if len(filtered) != 1 {
			t.Fatalf("FilterByProject('') returned %d items, want 1", len(filtered))
		}
		if filtered[0].EnvVar != "AWS_KEY" {
			t.Errorf("got %q, want AWS_KEY", filtered[0].EnvVar)
		}
	})
}

func TestEnvVarNames(t *testing.T) {
	items := ParseAll([]string{"A-x", "B", "C-y"})
	names := EnvVarNames(items)

	want := []string{"A", "B", "C"}
	if len(names) != len(want) {
		t.Fatalf("len = %d, want %d", len(names), len(want))
	}
	for i, w := range want {
		if names[i] != w {
			t.Errorf("names[%d] = %q, want %q", i, names[i], w)
		}
	}
}
