# Chief

**Autonomous PRD Agent** - A TUI application that orchestrates Claude to implement features from product requirements documents.

Chief takes your PRD (Product Requirements Document) and autonomously works through each user story using Claude Code, tracking progress in real-time through an interactive terminal interface.

## Quick Start

```bash
# Install via Homebrew (macOS/Linux)
brew install minicodemonkey/chief/chief

# Or via install script
curl -fsSL https://chief.codemonkey.io/install.sh | sh

# Create your first PRD
chief init

# Launch the TUI to start working
chief
```

## How It Works

Chief operates as an autonomous agent loop that orchestrates Claude Code to implement user stories:

```
                              ┌──────────────────────────────────────┐
                              │           Chief TUI                  │
                              │  ┌─────────────┬─────────────────┐   │
                              │  │  Stories    │    Details      │   │
                              │  │  ────────   │    ────────     │   │
                              │  │  ✓ US-001   │  Title: ...     │   │
                              │  │  ● US-002   │  Status: ...    │   │
                              │  │  ○ US-003   │  Criteria: ...  │   │
                              │  └─────────────┴─────────────────┘   │
                              │  [Activity: Running tool: Bash...]   │
                              └──────────────────┬───────────────────┘
                                                 │
                                                 ▼
┌────────────────────────┐              ┌───────────────────┐              ┌────────────────────────┐
│       prd.json         │◀────────────▶│    Agent Loop     │─────────────▶│     Claude Code        │
│  ────────────────────  │              │   ───────────     │              │  ───────────────────   │
│  project: "My App"     │  read/write  │  1. Read PRD      │   stream     │  - Read files          │
│  userStories:          │              │  2. Pick story    │    JSON      │  - Edit code           │
│    - id: US-001        │              │  3. Run Claude    │              │  - Run commands        │
│      passes: true      │              │  4. Check status  │              │  - Run tests           │
│    - id: US-002        │              │  5. Repeat        │              │  - Commit changes      │
│      inProgress: true  │              └───────────────────┘              └────────────────────────┘
└────────────────────────┘                       │
                                                 │
                                                 ▼
                                        ┌───────────────────┐
                                        │    progress.txt   │
                                        │   ─────────────   │
                                        │  ## US-001        │
                                        │  - Implemented... │
                                        │  - Learnings...   │
                                        └───────────────────┘
```

### Core Concepts

1. **PRD (Product Requirements Document)**: A JSON file containing user stories with acceptance criteria
2. **User Stories**: Individual features to implement, each with ID, title, description, and acceptance criteria
3. **Agent Loop**: Chief repeatedly invokes Claude Code until all stories are complete
4. **Progress Tracking**: Each story's status is tracked in `prd.json` and detailed logs in `progress.txt`

### Workflow

1. **Create PRD** - `chief init` launches Claude to help you define your requirements in `prd.md`
2. **Convert** - Chief automatically converts `prd.md` to the structured `prd.json` format
3. **Launch TUI** - `chief` starts the interactive dashboard
4. **Start Loop** - Press `s` to begin the autonomous implementation
5. **Monitor** - Watch Claude work through stories, view logs, and control execution
6. **Complete** - Chief plays a sound when all stories pass

## Installation

### Homebrew (Recommended)

```bash
# macOS (Intel or Apple Silicon) and Linux
brew install minicodemonkey/chief/chief
```

### Install Script

```bash
# Auto-detects OS and architecture
curl -fsSL https://chief.codemonkey.io/install.sh | sh

# Specify version
curl -fsSL https://chief.codemonkey.io/install.sh | sh -s -- --version v1.0.0

# Custom install directory
CHIEF_INSTALL_DIR=/opt/bin curl -fsSL https://chief.codemonkey.io/install.sh | sh
```

### From Source

```bash
# Clone and build
git clone https://github.com/minicodemonkey/chief.git
cd chief
go build -o chief ./cmd/chief

# Or with version info
go build -ldflags "-X main.Version=$(git describe --tags)" -o chief ./cmd/chief
```

### Prerequisites

- **Claude Code CLI** (`claude`) must be installed and authenticated
- Go 1.21+ (for building from source)

## Usage

### Commands

| Command | Description |
|---------|-------------|
| `chief` | Launch TUI with default PRD (`.chief/prds/main/`) |
| `chief <name>` | Launch TUI with named PRD (`.chief/prds/<name>/`) |
| `chief init` | Create a new PRD interactively |
| `chief init <name>` | Create a named PRD |
| `chief init <name> "context"` | Create PRD with context hint |
| `chief edit` | Edit the default PRD |
| `chief edit <name>` | Edit a named PRD |
| `chief status` | Show progress for default PRD |
| `chief status <name>` | Show progress for named PRD |
| `chief list` | List all PRDs with progress |
| `chief --help` | Show help message |
| `chief --version` | Show version |

### Options

| Flag | Description |
|------|-------------|
| `--max-iterations N`, `-n N` | Set maximum iterations (default: 10) |
| `--no-sound` | Disable completion sound notifications |
| `--verbose` | Show raw Claude output in log |
| `--merge` | Auto-merge progress on conversion conflicts |
| `--force` | Auto-overwrite on conversion conflicts |

### Examples

```bash
# Create a PRD with context
chief init auth "JWT authentication for REST API"

# Launch with custom iteration limit
chief -n 20 auth

# Check status without TUI
chief status auth
# Output:
# Auth System
# 3/5 stories complete
#
# Incomplete stories:
#   US-004: Password Reset
#   US-005: OAuth Integration (in progress)

# List all PRDs
chief list
# Output:
# auth: Auth System (3/5, 60%)
# main: My Project (10/10, 100%)
```

## Keyboard Shortcuts

### Loop Control

| Key | Action |
|-----|--------|
| `s` | Start/Resume the agent loop |
| `p` | Pause (stops after current iteration) |
| `x` | Stop immediately |

### Views

| Key | Action |
|-----|--------|
| `t` | Toggle log view |
| `l` | Open PRD picker |
| `?` | Show help overlay |

### Navigation (Dashboard)

| Key | Action |
|-----|--------|
| `j` / `↓` | Next story |
| `k` / `↑` | Previous story |

### Navigation (Log View)

| Key | Action |
|-----|--------|
| `j` / `↓` | Scroll down |
| `k` / `↑` | Scroll up |
| `Ctrl+D` | Page down |
| `Ctrl+U` | Page up |
| `g` | Go to top |
| `G` | Go to bottom |

### PRD Picker

| Key | Action |
|-----|--------|
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `Enter` | Select PRD |
| `n` | Create new PRD |
| `s/p/x` | Start/Pause/Stop selected PRD |
| `Esc` / `l` | Close picker |

### General

| Key | Action |
|-----|--------|
| `q` | Quit |
| `Ctrl+C` | Quit |
| `Esc` | Close overlay/modal |

## PRD Format Reference

### prd.json Schema

```json
{
  "project": "Project Name",
  "description": "Brief project description",
  "userStories": [
    {
      "id": "US-001",
      "title": "Feature Title",
      "description": "As a user, I want X so that Y",
      "acceptanceCriteria": [
        "Criterion 1",
        "Criterion 2",
        "All tests pass"
      ],
      "priority": 1,
      "passes": false,
      "inProgress": false
    }
  ]
}
```

### Fields

| Field | Type | Description |
|-------|------|-------------|
| `project` | string | Project name displayed in header |
| `description` | string | Brief project description |
| `userStories` | array | List of user stories |
| `userStories[].id` | string | Unique identifier (e.g., "US-001") |
| `userStories[].title` | string | Short feature title |
| `userStories[].description` | string | Full user story description |
| `userStories[].acceptanceCriteria` | array | List of criteria that must be met |
| `userStories[].priority` | number | Priority order (lower = higher priority) |
| `userStories[].passes` | boolean | Whether the story is complete |
| `userStories[].inProgress` | boolean | Whether Claude is currently working on it |

### Story Selection Logic

Chief selects the next story to work on using this logic:

1. **Interrupted story first**: If any story has `inProgress: true`, it's selected (handles restarts)
2. **Lowest priority next**: Among stories with `passes: false`, the one with lowest `priority` value is selected
3. **Complete**: When all stories have `passes: true`, the loop ends

### prd.md Format

When using `chief init` or `chief edit`, you write in markdown format:

```markdown
# Project Name

Brief project description.

## User Stories

### US-001: Feature Title (Priority: 1)

As a user, I want to do something so that I get value.

**Acceptance Criteria:**
- First criterion
- Second criterion
- Tests pass

### US-002: Another Feature (Priority: 2)

As a developer, I need something to enable X.

**Acceptance Criteria:**
- Criterion A
- Criterion B
```

Chief converts this to `prd.json` automatically.

## Directory Structure

```
your-project/
├── .chief/
│   └── prds/
│       ├── main/
│       │   ├── prd.md          # Human-readable PRD
│       │   ├── prd.json        # Machine-readable PRD
│       │   ├── progress.txt    # Implementation log
│       │   └── claude.log      # Raw Claude output
│       └── auth/
│           ├── prd.md
│           ├── prd.json
│           ├── progress.txt
│           └── claude.log
└── ... (your project files)
```

## Parallel Execution

Chief supports running multiple PRDs simultaneously:

1. Open the PRD picker with `l`
2. Navigate to a PRD and press `s` to start it
3. Switch to another PRD and start it too
4. Each PRD runs in its own goroutine with its own Claude process
5. The header shows `▶ +N PRDs` when multiple are running
6. Completion sound plays when ANY PRD completes

## Progress Protection

When editing a PRD that has progress:

```bash
chief edit auth
# After editing, Chief detects prd.json has progress

⚠️  Warning: prd.json has progress (3 stories with status)

How would you like to proceed?

  [m] Merge  - Keep status for matching story IDs, add new, drop removed
  [o] Overwrite - Discard all progress and use the new PRD
  [c] Cancel - Cancel conversion and keep existing prd.json

Choice [m/o/c]:
```

Use flags to skip the prompt:
- `--merge`: Auto-merge without prompting
- `--force`: Auto-overwrite without prompting

## Troubleshooting

### "Claude not found"

Ensure Claude Code CLI is installed and in your PATH:
```bash
which claude
claude --version
```

### "Permission denied" errors

Chief runs Claude with `--dangerously-skip-permissions` to enable autonomous operation. Ensure:
- You trust the PRD content
- Claude has appropriate access to your project

### No sound on completion

- Check your system audio
- Linux users need CGO and libasound2-dev for audio support
- Use `--no-sound` if audio causes issues

### PRD not updating in TUI

Chief watches `prd.json` for changes. If updates aren't appearing:
- Check file permissions
- Verify the file is valid JSON
- The watcher only triggers on status field changes

### Loop not progressing

Check `claude.log` in the PRD directory for detailed output:
```bash
tail -f .chief/prds/main/claude.log
```

### "Max iterations reached"

Increase the limit:
```bash
chief -n 50
```

Or pause and resume to continue with more iterations.

## Development

### Building

```bash
# Build for current platform
go build ./cmd/chief

# Run tests
go test ./...

# Run with verbose output
go test -v ./...
```

### Project Structure

```
chief/
├── cmd/chief/          # Main entry point
├── internal/
│   ├── cmd/            # CLI commands (init, edit, status, list)
│   ├── loop/           # Agent loop, parser, manager
│   ├── prd/            # PRD types, loader, watcher, generator
│   ├── tui/            # Bubble Tea UI components
│   └── notify/         # Audio notifications
└── embed/              # Embedded prompts and assets
```

### Release

```bash
# Create a snapshot release
goreleaser release --snapshot --clean

# Create a real release (requires GITHUB_TOKEN)
goreleaser release --clean
```

## License

MIT License - see LICENSE file for details.

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure `go test ./...` passes
5. Submit a pull request

## Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [Claude Code](https://claude.ai/code) - AI coding assistant
