package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// press feeds one key message through Update and returns the updated model.
func press(t *testing.T, m Model, k tea.KeyMsg) Model {
	t.Helper()
	nm, _ := m.Update(k)
	got, ok := nm.(Model)
	if !ok {
		t.Fatalf("Update returned %T, want Model", nm)
	}
	return got
}

func TestUpdate_spaceIsCaptured(t *testing.T) {
	m := New([]string{"sk1-ex-etapp", "sk1-ex-tapp"})

	m = press(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("sk1")})
	m = press(t, m, tea.KeyMsg{Type: tea.KeySpace})
	m = press(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("tapp")})

	if m.query != "sk1 tapp" {
		t.Fatalf("query = %q, want %q", m.query, "sk1 tapp")
	}
	if len(m.matches) != 2 {
		t.Fatalf("matches = %v, want both hosts", m.matches)
	}
}

func TestUpdate_backspaceRemovesLastRune(t *testing.T) {
	m := New([]string{"sk1-ex-etapp"})
	m = press(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("sk1")})

	// Both DEL (KeyBackspace) and Ctrl+H (KeyCtrlH) must delete a character.
	m = press(t, m, tea.KeyMsg{Type: tea.KeyBackspace})
	if m.query != "sk" {
		t.Fatalf("after KeyBackspace query = %q, want %q", m.query, "sk")
	}
	m = press(t, m, tea.KeyMsg{Type: tea.KeyCtrlH})
	if m.query != "s" {
		t.Fatalf("after KeyCtrlH query = %q, want %q", m.query, "s")
	}

	// Deleting past the start must be a no-op, not a panic.
	m = press(t, m, tea.KeyMsg{Type: tea.KeyBackspace})
	m = press(t, m, tea.KeyMsg{Type: tea.KeyBackspace})
	if m.query != "" {
		t.Fatalf("query = %q, want empty", m.query)
	}
}
