# sshpick Design

## Goal

A small Go TUI that reads hostnames from `~/.ssh/ssh_hosts`, filters them as the user types (substring), and runs `ssh <host>` after Enter.

## Scope (v1)

**In**
- Read-only host list from `~/.ssh/ssh_hosts` (one hostname per line)
- Live substring filter (case-insensitive); empty query shows no list
- TUI: type to filter, arrows to select, Enter to connect, Esc/Ctrl+C to quit
- Always require Enter (no auto-connect on single match)
- Launch real `ssh` binary with the selected host

**Out**
- Writing or editing the hosts file
- Parsing `~/.ssh/config`
- Fuzzy matching
- Print-only / copy mode
- Custom ssh flags

## Hosts file

- Path: `~/.ssh/ssh_hosts`
- Format: one hostname per line (e.g. `nsk1-ex-tav01`)
- Blank lines ignored
- Tool never creates or overwrites this file

## Components

| Piece | Responsibility |
|---|---|
| `internal/hosts` | Load hostnames from disk |
| `internal/filter` | Substring match over the list |
| `internal/tui` | Interactive UI (Bubble Tea) |
| `internal/sshcmd` | Exec `ssh <host>` |
| `main` | Wire pieces together and handle top-level errors |

## Errors

- Missing / unreadable `ssh_hosts` → message + non-zero exit
- Empty host list → message + non-zero exit
- Enter with no selection → stay in TUI
- `ssh` not found → error after selection

## Binary name

`sshpick`
