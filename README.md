# sshpick

Small Go TUI: type to filter hostnames from `~/.ssh/ssh_hosts`, pick one, run `ssh`.

## Setup

Create the hosts file (the tool never creates it for you):

```bash
mkdir -p ~/.ssh
printf '%s\n' 'nsk-tav01' 'sk2-web03' >> ~/.ssh/ssh_hosts
chmod 600 ~/.ssh/ssh_hosts
```

## Build & run

```bash
cd sshpick
go test ./...
go build -o sshpick .
./sshpick
```

## Keys

| Key | Action |
|---|---|
| printable chars | filter (substring, case-insensitive) |
| Backspace | delete last filter character |
| ↑ / ↓ | move selection |
| Enter | `ssh` to highlighted host |
| Esc / Ctrl+C | quit |

## Layout

See comments in `main.go` and under `internal/` — written for learning Go + Bubble Tea.
