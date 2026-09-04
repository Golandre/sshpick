// Command sshpick is a small TUI helper: filter hostnames and run ssh.
//
// Project layout (typical for a Go CLI):
//
//	main.go              — program entry (package main)
//	internal/hosts       — read ~/.ssh/ssh_hosts
//	internal/filter      — substring matching
//	internal/tui         — Bubble Tea UI
//	internal/sshcmd      — launch OpenSSH
//
// package main is special: it builds an executable, not a library.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/azimut/sshpick/internal/hosts"
	"github.com/azimut/sshpick/internal/sshcmd"
	"github.com/azimut/sshpick/internal/tui"
)

func main() {
	// os.Exit with a non-zero code tells the shell the program failed.
	// We keep main thin and put logic in run() so we can return errors cleanly.
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "sshpick: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	path, err := hostsPath()
	if err != nil {
		return err
	}

	list, err := hosts.Load(path)
	if err != nil {
		return err
	}

	chosen, err := tui.Run(list)
	if err != nil {
		return fmt.Errorf("ui: %w", err)
	}
	// Empty chosen means the user pressed Esc / Ctrl+C — quiet success exit.
	if chosen == "" {
		return nil
	}

	return sshcmd.Run(chosen)
}

// hostsPath resolves ~/.ssh/ssh_hosts using the current user's home directory.
func hostsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home directory: %w", err)
	}
	// filepath.Join builds an OS-correct path (slashes on Unix).
	return filepath.Join(home, ".ssh", "ssh_hosts"), nil
}
