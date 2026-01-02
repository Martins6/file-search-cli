# fsc - File Search CLI

A Finder-like CLI tool for interactive file browsing in the terminal.

## Features

- **First-level browsing**: Always shows first-level directory contents (no recursion)
- **Interactive filtering**: Filter by character matching or regex patterns as you type
- **Intuitive navigation**: Navigate directories with Enter, open files in editor with 'e'
- **Keyboard shortcuts**: Vim-like navigation support
- **Customizable views**: Show/hidden files, directories only, or files only

## Prerequisites

- Go 1.21 or later
- Task runner (<https://taskfile.dev/>)

## Installation

### Build from source

```bash
# Clone the repository
cd /path/to/file-search-cli

# Download dependencies
task deps

# Build the binary
task build

# Install to GOBIN or GOPATH/bin
task install
```

### Manual installation

```bash
go install ./...
```

## Usage

### Basic usage

```bash
# Start browsing in a specific directory
fsc /path/to/dir

# Browse current directory
fsc .

# Browse home directory
fsc ~
```

### Flags

```bash
# Show hidden files
fsc /path --hidden

# Show only directories
fsc /path --dirs-only

# Show only files
fsc /path --files-only
```

### Example workflow

```
$ fsc ~/projects
┌─────────────────────────────────────┐
│ fsc - File Search CLI              │
└─────────────────────────────────────┘
📁 /home/user/projects
──────────────────────────────────────
  📁 fsc
  📁 myapp
  📁 website
  📄 README.md

Filter: app_
↑/j: up | ↓/k: down | enter: navigate | e: open editor | q: quit

[Press Enter on "myapp"]
┌─────────────────────────────────────┐
│ fsc - File Search CLI              │
└─────────────────────────────────────┘
📁 /home/user/projects/myapp
──────────────────────────────────────
  📁 src
  📄 main.go
  📄 go.mod
  📄 README.md

Filter: _
↑/j: up | ↓/k: down | enter: navigate | e: open editor | q: quit

[Type "go"]
┌─────────────────────────────────────┐
│ fsc - File Search CLI              │
└─────────────────────────────────────┘
📁 /home/user/projects/myapp
──────────────────────────────────────
  📄 go.mod

Filter: go_
↑/j: up | ↓/k: down | enter: navigate | e: open editor | q: quit

[Press 'e' on go.mod]
[Opens go.mod in your $EDITOR]
```

## Keyboard Shortcuts

| Key        | Action                           |
|------------|----------------------------------|
| ↑ / k      | Move selection up                |
| ↓ / j      | Move selection down              |
| Enter      | Navigate into selected directory |
| e          | Open selected item in editor     |
| Backspace  | Delete last character in filter |
| q / Esc    | Quit                             |
| Ctrl+C     | Quit                             |

## Filtering

### Character matching (default)

Type characters to filter files/folders by name (case-insensitive):

```
Filter: go_
Shows: go.mod, go.sum, main.go (if contains "go")
```

### Regex matching

Prefix your pattern with `/` to use regex:

```
Filter: /.*\.go$
Shows: All .go files
```

```
Filter: /[A-Z]
Shows: Files starting with uppercase letters
```

## Development

### Available tasks

```bash
# List all tasks
task

# Download dependencies
task deps

# Build binary
task build

# Clean build artifacts
task clean

# Run tests
task test

# Install to GOBIN
task install

# Build and run
task run
```

### Project structure

```
fsc/
├── cmd/
│   └── root.go          # Cobra CLI command definition
├── pkg/
│   ├── editor/
│   │   └── opener.go    # Editor integration
│   ├── filesystem/
│   │   ├── entry.go     # Entry struct
│   │   ├── scanner.go   # Directory scanning
│   │   ├── navigator.go # Navigation logic
│   │   └── filter.go    # Filtering logic
│   └── ui/
│       ├── model.go     # Bubbletea model
│       ├── view.go      # TUI rendering
│       ├── keymap.go    # Key bindings
│       └── messages.go  # Message types
├── go.mod
├── main.go              # Entry point
└── Taskfile.yaml        # Build configuration
```

## License

MIT License - see LICENSE file for details
