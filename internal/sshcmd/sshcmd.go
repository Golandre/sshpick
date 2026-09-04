// Package sshcmd starts the real OpenSSH client.
package sshcmd

import (
	"fmt"
	"os"
	"os/exec"
)

// Run replaces this process's work with `ssh host` by attaching the current
// terminal's stdin/stdout/stderr to the ssh child process and waiting for it.
//
// Why not "print the host and exit"? Because the desired UX is: pick → connect.
// Using exec.Command (instead of syscall.Exec) keeps error handling simple and
// works the same on Linux/macOS for a learner project.
func Run(host string) error {
	if host == "" {
		return fmt.Errorf("empty host")
	}

	// LookPath finds ssh in $PATH (e.g. /usr/bin/ssh).
	bin, err := exec.LookPath("ssh")
	if err != nil {
		return fmt.Errorf("ssh not found in PATH: %w", err)
	}

	cmd := exec.Command(bin, host)
	// Connect the child's I/O to our terminal so the SSH session is interactive
	// (password prompts, shell, etc. work as if you typed `ssh host` yourself).
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run starts ssh and waits until it exits. The returned error is non-nil
	// if ssh exits with a non-zero status (or fails to start).
	return cmd.Run()
}
