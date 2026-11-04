// Package filter narrows a host list as the user types.
package filter

import "strings"

// Match returns hosts that contain every whitespace-separated token of query
// as a case-insensitive substring.
//
// An empty (or whitespace-only) query returns nil on purpose: the TUI should
// show no list until the user types at least one character.
//
// We return a new slice; we never modify the input hosts slice.
func Match(query string, hosts []string) []string {
	// strings.Fields splits on any whitespace and drops empty tokens, so
	// "sk1  tapp" behaves the same as "sk1 tapp".
	fields := strings.Fields(query)
	if len(fields) == 0 {
		return nil
	}

	// Fold each token to lower case once so "NSK1" matches "nsk1-ex-tav01".
	tokens := make([]string, len(fields))
	for i, f := range fields {
		tokens[i] = strings.ToLower(f)
	}

	var out []string
	for _, h := range hosts {
		lh := strings.ToLower(h)
		match := true
		for _, t := range tokens {
			if !strings.Contains(lh, t) {
				match = false
				break
			}
		}
		if match {
			out = append(out, h)
		}
	}
	return out
}
