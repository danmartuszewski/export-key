package backend

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestDotenvListItems(t *testing.T) {
	path := writeTempEnv(t, `OPENAI_API_KEY=sk-123
STRIPE_KEY=sk_test_456
AWS_KEY=AKIA789
`)
	d := NewDotenv([]string{path})
	items, err := d.ListItems()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 3 {
		t.Fatalf("got %d items, want 3", len(items))
	}
	// Verify insertion order is preserved
	want := []string{"OPENAI_API_KEY", "STRIPE_KEY", "AWS_KEY"}
	for i, w := range want {
		if items[i] != w {
			t.Errorf("items[%d] = %q, want %q", i, items[i], w)
		}
	}
}

func TestDotenvGetSecret(t *testing.T) {
	path := writeTempEnv(t, `KEY_A=value_a
KEY_B=value_b
`)
	d := NewDotenv([]string{path})

	val, err := d.GetSecret("KEY_A")
	if err != nil {
		t.Fatal(err)
	}
	if val != "value_a" {
		t.Errorf("got %q, want %q", val, "value_a")
	}
}

func TestDotenvGetSecretNotFound(t *testing.T) {
	path := writeTempEnv(t, "KEY=val\n")
	d := NewDotenv([]string{path})

	_, err := d.GetSecret("MISSING")
	if err == nil {
		t.Error("expected error for missing key")
	}
}

func TestDotenvUnquoting(t *testing.T) {
	path := writeTempEnv(t, `DOUBLE="hello world"
SINGLE='hello world'
NONE=hello world
EMPTY=""
`)
	d := NewDotenv([]string{path})

	tests := []struct {
		key  string
		want string
	}{
		{"DOUBLE", "hello world"},
		{"SINGLE", "hello world"},
		{"NONE", "hello world"},
		{"EMPTY", ""},
	}

	for _, tt := range tests {
		val, err := d.GetSecret(tt.key)
		if err != nil {
			t.Errorf("GetSecret(%q): %v", tt.key, err)
			continue
		}
		if val != tt.want {
			t.Errorf("GetSecret(%q) = %q, want %q", tt.key, val, tt.want)
		}
	}
}

func TestDotenvCommentsAndBlanks(t *testing.T) {
	path := writeTempEnv(t, `# This is a comment
KEY_A=a

# Another comment

KEY_B=b
`)
	d := NewDotenv([]string{path})
	items, err := d.ListItems()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("got %d items, want 2 (comments/blanks should be skipped)", len(items))
	}
}

func TestDotenvMissingFile(t *testing.T) {
	d := NewDotenv([]string{"/nonexistent/.env"})
	_, err := d.ListItems()
	if err == nil {
		t.Error("expected error for missing file with no keys")
	}
}

func TestDotenvMultipleFiles(t *testing.T) {
	dir := t.TempDir()
	path1 := filepath.Join(dir, "a.env")
	path2 := filepath.Join(dir, "b.env")

	if err := os.WriteFile(path1, []byte("KEY_A=from_a\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path2, []byte("KEY_B=from_b\n"), 0644); err != nil {
		t.Fatal(err)
	}

	d := NewDotenv([]string{path1, path2})
	items, err := d.ListItems()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("got %d items, want 2", len(items))
	}

	valA, _ := d.GetSecret("KEY_A")
	valB, _ := d.GetSecret("KEY_B")
	if valA != "from_a" || valB != "from_b" {
		t.Errorf("got A=%q B=%q", valA, valB)
	}
}

func TestDotenvDuplicateKey(t *testing.T) {
	path := writeTempEnv(t, `KEY=first
KEY=second
`)
	d := NewDotenv([]string{path})

	// Last value wins
	val, err := d.GetSecret("KEY")
	if err != nil {
		t.Fatal(err)
	}
	if val != "second" {
		t.Errorf("got %q, want %q (last value should win)", val, "second")
	}

	// But only one entry in the order
	items, _ := d.ListItems()
	if len(items) != 1 {
		t.Errorf("got %d items, want 1 (duplicate key should not create duplicate entry)", len(items))
	}
}

func TestDotenvNoEqualsSign(t *testing.T) {
	path := writeTempEnv(t, `VALID=value
INVALID_LINE_NO_EQUALS
ALSO_VALID=123
`)
	d := NewDotenv([]string{path})
	items, err := d.ListItems()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("got %d items, want 2 (line without = should be skipped)", len(items))
	}
}

func TestDotenvEmptyValue(t *testing.T) {
	path := writeTempEnv(t, "KEY=\n")
	d := NewDotenv([]string{path})
	val, err := d.GetSecret("KEY")
	if err != nil {
		t.Fatal(err)
	}
	if val != "" {
		t.Errorf("got %q, want empty string", val)
	}
}

func TestDotenvValueWithEquals(t *testing.T) {
	path := writeTempEnv(t, "DB_URL=postgres://user:pass@host/db?sslmode=require\n")
	d := NewDotenv([]string{path})
	val, err := d.GetSecret("DB_URL")
	if err != nil {
		t.Fatal(err)
	}
	want := "postgres://user:pass@host/db?sslmode=require"
	if val != want {
		t.Errorf("got %q, want %q", val, want)
	}
}

func TestDotenvWhitespace(t *testing.T) {
	path := writeTempEnv(t, "  KEY  =  value  \n")
	d := NewDotenv([]string{path})
	val, err := d.GetSecret("KEY")
	if err != nil {
		t.Fatal(err)
	}
	if val != "value" {
		t.Errorf("got %q, want %q", val, "value")
	}
}

func TestDotenvCaching(t *testing.T) {
	path := writeTempEnv(t, "KEY=original\n")
	d := NewDotenv([]string{path})

	// First load
	val1, _ := d.GetSecret("KEY")
	if val1 != "original" {
		t.Fatalf("first read got %q", val1)
	}

	// Overwrite file
	if err := os.WriteFile(path, []byte("KEY=modified\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Should still return cached value
	val2, _ := d.GetSecret("KEY")
	if val2 != "original" {
		t.Errorf("cached read got %q, want %q (should be cached)", val2, "original")
	}
}
