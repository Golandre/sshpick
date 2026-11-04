package hosts

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoad_skipsBlankLines checks that empty lines do not become hosts.
func TestLoad_skipsBlankLines(t *testing.T) {
	// t.TempDir() creates a temporary directory that Go cleans up after the test.
	dir := t.TempDir()
	path := filepath.Join(dir, "ssh_hosts")

	content := "nsk1-tav01\n\n  \ntpk2-web03\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	want := []string{"nsk1-tav01", "tpk2-web03"}
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d (%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("got[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestLoad_missingFile(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "no-such-file"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_emptyFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ssh_hosts")
	if err := os.WriteFile(path, []byte("\n\n"), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for empty host list")
	}
}
