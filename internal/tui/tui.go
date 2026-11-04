// Package tui is the interactive terminal UI built with Bubble Tea.
//
// Bubble Tea uses The Elm Architecture:
//   - Model  = all UI state
//   - Update = handle one message (key press, etc.), return new state
//   - View   = render state as a string for the terminal
//
// The program loop: View → wait for Msg → Update → View → …
package tui

import (
	"fmt"
	"strings"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/azimut/sshpick/internal/filter"
)

// Model holds everything the UI needs between key presses.
// Exported fields are not required here; keeping them unexported is fine
// because only this package's Run function uses the model.
type Model struct {
	allHosts []string // full list from the file (never mutated)
	query    string   // what the user has typed so far
	matches  []string // filtered view of allHosts
	cursor   int      // index into matches (which row is highlighted)
	chosen   string   // set when user presses Enter on a valid row
	quitting bool
}

// New builds the initial model. matches start empty until the user types.
func New(hosts []string) Model {
	return Model{
		allHosts: hosts,
		matches:  nil,
		cursor:   0,
	}
}

// Chosen returns the host selected with Enter, or "" if the user quit.
func (m Model) Chosen() string {
	return m.chosen
}

// Init is required by tea.Model. We have no async startup work, so return nil.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles one event. It returns the (possibly updated) model and an
// optional command (side effect). tea.Quit ends the program loop.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// tea.KeyMsg is sent for keyboard input.
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit

		case tea.KeyEnter:
			// Only connect if the highlighted row exists.
			if m.cursor >= 0 && m.cursor < len(m.matches) {
				m.chosen = m.matches[m.cursor]
				return m, tea.Quit
			}
			// No selection → stay in the UI (do nothing).
			return m, nil

		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil

		case tea.KeyDown:
			if m.cursor < len(m.matches)-1 {
				m.cursor++
			}
			return m, nil

		case tea.KeyBackspace, tea.KeyCtrlH:
			// KeyBackspace is the DEL byte (0x7f); KeyCtrlH is the 0x08 byte
			// some terminals send for the Backspace key. Handle both so the
			// user can always delete a typo.
			if len(m.query) > 0 {
				// Remove the last rune (not just the last byte) so multi-byte
				// Unicode characters are deleted cleanly.
				_, size := utf8.DecodeLastRuneInString(m.query)
				m.query = m.query[:len(m.query)-size]
				m.refilter()
			}
			return m, nil

		case tea.KeySpace:
			// Spaces arrive as their own key type (not KeyRunes); capture
			// them so multi-word filters like "sk1 tapp" work as expected.
			m.query += " "
			m.refilter()
			return m, nil

		case tea.KeyRunes:
			// Printable characters arrive as runes (Unicode code points).
			m.query += string(msg.Runes)
			m.refilter()
			return m, nil
		}
	}
	return m, nil
}

// refilter applies the current query and keeps the cursor in range.
func (m *Model) refilter() {
	m.matches = filter.Match(m.query, m.allHosts)
	if m.cursor >= len(m.matches) {
		m.cursor = len(m.matches) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

// View draws the current screen. Bubble Tea prints this string each frame.
func (m Model) View() string {
	if m.quitting && m.chosen == "" {
		return ""
	}

	var b strings.Builder
	b.WriteString("sshpick — type to filter hosts, ↑/↓ select, Enter connect, Esc quit\n\n")
	b.WriteString(fmt.Sprintf("Filter: %s█\n\n", m.query))

	if m.query == "" {
		b.WriteString("(type at least one character to see matches)\n")
		return b.String()
	}

	if len(m.matches) == 0 {
		b.WriteString("No matches.\n")
		return b.String()
	}

	for i, h := range m.matches {
		cursor := "  "
		if i == m.cursor {
			cursor = "> " // highlight the selected row
		}
		b.WriteString(cursor + h + "\n")
	}
	return b.String()
}

// Run starts the Bubble Tea program and returns the chosen host (or "").
func Run(hosts []string) (string, error) {
	m := New(hosts)
	p := tea.NewProgram(m)
	final, err := p.Run()
	if err != nil {
		return "", err
	}
	// Type assert: Run returns tea.Model; we know it is our Model.
	out, ok := final.(Model)
	if !ok {
		return "", fmt.Errorf("unexpected model type %T", final)
	}
	return out.Chosen(), nil
}
