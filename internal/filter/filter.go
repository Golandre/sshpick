// Package filter narrows a host list as the user types.
package filter

import "strings"

// Match returns hosts that contain query as a substring (case-insensitive).
//
// An empty query returns nil/empty slice on purpose: the TUI should show
// no list until the user types at least one character.
//
// We return a new slice; we never modify the input hosts slice.
func Match(query string, hosts []string) []string {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil
	}

	// Fold both sides to lower case so "NSK" matches "nsk1-ex-tav01".
	q := strings.ToLower(query)

	var out []string
	for _, h := range hosts {
		if strings.Contains(strings.ToLower(h), q) {
			out = append(out, h)
		}
	}
	return out
}
