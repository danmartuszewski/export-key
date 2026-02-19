package keyitem

import (
	"slices"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		title    string
		envVar   string
		projects []string
	}{
		{"OPENAI_API_KEY", "OPENAI_API_KEY", nil},
		{"OPENAI_API_KEY-myapp", "OPENAI_API_KEY", []string{"myapp"}},
		{"GEMINI_API_KEY-projectA", "GEMINI_API_KEY", []string{"projectA"}},
		{"STRIPE_KEY-web,api", "STRIPE_KEY", []string{"web", "api"}},
		{"KEY-a,b,c", "KEY", []string{"a", "b", "c"}},
		{"KEY-", "KEY", nil},              // trailing dash, no projects
		{"-project", "", []string{"project"}}, // leading dash, empty env var
		{"SIMPLE", "SIMPLE", nil},         // no dash at all
		{"A-B", "A", []string{"B"}},       // minimal
		{"", "", nil},                     // empty string
		{"NO_DASH_HERE", "NO_DASH_HERE", nil},
		{"AWS_ACCESS_KEY-infra", "AWS_ACCESS_KEY", []string{"infra"}},
		{"KEY-a,,b", "KEY", []string{"a", "b"}},       // empty segment ignored
		{"KEY-a, b", "KEY", []string{"a", "b"}},       // whitespace trimmed
		{"KEY-,", "KEY", nil},                         // only commas, no projects
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
			if !slices.Equal(item.Projects, tt.projects) {
				t.Errorf("Parse(%q).Projects = %v, want %v", tt.title, item.Projects, tt.projects)
			}
		})
	}
}

func TestHasProject(t *testing.T) {
	if Parse("KEY").HasProject() {
		t.Error("KEY should not have a project")
	}
	if !Parse("KEY-app").HasProject() {
		t.Error("KEY-app should have a project")
	}
	if !Parse("KEY-a,b").HasProject() {
		t.Error("KEY-a,b should have projects")
	}
}

func TestProjectString(t *testing.T) {
	if got := Parse("KEY-web,api").ProjectString(); got != "web,api" {
		t.Errorf("ProjectString() = %q, want %q", got, "web,api")
	}
	if got := Parse("KEY").ProjectString(); got != "" {
		t.Errorf("ProjectString() = %q, want empty", got)
	}
}

func TestParseAll(t *testing.T) {
	titles := []string{"OPENAI_API_KEY-myapp", "AWS_KEY", "STRIPE-prod"}
	items := ParseAll(titles)

	if len(items) != 3 {
		t.Fatalf("ParseAll returned %d items, want 3", len(items))
	}

	if items[0].EnvVar != "OPENAI_API_KEY" || !slices.Equal(items[0].Projects, []string{"myapp"}) {
		t.Errorf("items[0] = %+v", items[0])
	}
	if items[1].EnvVar != "AWS_KEY" || items[1].HasProject() {
		t.Errorf("items[1] = %+v", items[1])
	}
	if items[2].EnvVar != "STRIPE" || !slices.Equal(items[2].Projects, []string{"prod"}) {
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
		"SHARED_KEY-myapp,infra",
	})

	t.Run("existing project", func(t *testing.T) {
		filtered := FilterByProject(items, "myapp")
		if len(filtered) != 4 { // 3 single + 1 shared
			t.Fatalf("FilterByProject(myapp) returned %d items, want 4", len(filtered))
		}
		envs := EnvVarNames(filtered)
		want := []string{"OPENAI_API_KEY", "STRIPE_KEY", "DB_URL", "SHARED_KEY"}
		for i, w := range want {
			if envs[i] != w {
				t.Errorf("envs[%d] = %q, want %q", i, envs[i], w)
			}
		}
	})

	t.Run("single match", func(t *testing.T) {
		filtered := FilterByProject(items, "infra")
		if len(filtered) != 2 { // GEMINI + SHARED
			t.Fatalf("FilterByProject(infra) returned %d items, want 2", len(filtered))
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
