// Package hosts loads SSH hostnames from a plain text file.
//
// In Go, a directory under internal/ can only be imported by code inside
// this module (the parent tree). That keeps "library" code private to sshpick.
package hosts

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Load reads path and returns non-empty hostnames, one per line.
//
// It never creates or modifies the file — we only open it for reading.
// Returning ([]string, error) is the common Go pattern: either a useful
// value, or an error explaining what went wrong (not both for failure cases).
func Load(path string) ([]string, error) {
	// os.Open opens an existing file. If the file is missing, err is non-nil
	// (usually *os.PathError wrapping os.ErrNotExist).
	f, err := os.Open(path)
	if err != nil {
		// %w wraps the original error so callers can use errors.Is(err, os.ErrNotExist).
		return nil, fmt.Errorf("open hosts file %q: %w", path, err)
	}
	// defer runs when Load returns — even if we return early later.
	// Always Close files you Open, or you leak file descriptors.
	defer f.Close()

	var list []string

	// bufio.Scanner reads line by line without loading the whole file at once.
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// TrimSpace removes leading/trailing spaces and newlines leftovers.
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue // skip blank lines
		}
		list = append(list, line)
	}

	// Scan can fail (e.g. extremely long line, I/O error). Check after the loop.
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read hosts file %q: %w", path, err)
	}

	if len(list) == 0 {
		return nil, fmt.Errorf("hosts file %q is empty (add one hostname per line)", path)
	}

	return list, nil
}
