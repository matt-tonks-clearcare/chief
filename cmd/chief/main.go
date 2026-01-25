package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/minicodemonkey/chief/internal/cmd"
	"github.com/minicodemonkey/chief/internal/notify"
	"github.com/minicodemonkey/chief/internal/prd"
	"github.com/minicodemonkey/chief/internal/tui"
)

// Version is set at build time via ldflags
var Version = "dev"

// TUIOptions holds the parsed command-line options for the TUI
type TUIOptions struct {
	PRDPath       string
	MaxIterations int
	NoSound       bool
	Verbose       bool
	Merge         bool
	Force         bool
}

func main() {
	// Handle subcommands first
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "init":
			runInit()
			return
		case "edit":
			runEdit()
			return
		case "status":
			runStatus()
			return
		case "list":
			runList()
			return
		case "help":
			printHelp()
			return
		case "--help", "-h":
			printHelp()
			return
		case "--version", "-v":
			fmt.Printf("chief version %s\n", Version)
			return
		}
	}

	// Parse flags for TUI mode
	opts := parseTUIFlags()

	// Handle special flags that were parsed
	if opts == nil {
		// Already handled (--help or --version)
		return
	}

	// Run the TUI
	runTUIWithOptions(opts)
}

// parseTUIFlags parses command-line flags for TUI mode
func parseTUIFlags() *TUIOptions {
	opts := &TUIOptions{
		PRDPath:       ".chief/prds/main/prd.json",
		MaxIterations: 10,
		NoSound:       false,
		Verbose:       false,
		Merge:         false,
		Force:         false,
	}

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		switch {
		case arg == "--help" || arg == "-h":
			printHelp()
			return nil
		case arg == "--version" || arg == "-v":
			fmt.Printf("chief version %s\n", Version)
			return nil
		case arg == "--no-sound":
			opts.NoSound = true
		case arg == "--verbose":
			opts.Verbose = true
		case arg == "--merge":
			opts.Merge = true
		case arg == "--force":
			opts.Force = true
		case arg == "--max-iterations" || arg == "-n":
			// Next argument should be the number
			if i+1 < len(os.Args) {
				i++
				n, err := strconv.Atoi(os.Args[i])
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: invalid value for %s: %s\n", arg, os.Args[i])
					os.Exit(1)
				}
				if n < 1 {
					fmt.Fprintf(os.Stderr, "Error: --max-iterations must be at least 1\n")
					os.Exit(1)
				}
				opts.MaxIterations = n
			} else {
				fmt.Fprintf(os.Stderr, "Error: %s requires a value\n", arg)
				os.Exit(1)
			}
		case strings.HasPrefix(arg, "--max-iterations="):
			val := strings.TrimPrefix(arg, "--max-iterations=")
			n, err := strconv.Atoi(val)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: invalid value for --max-iterations: %s\n", val)
				os.Exit(1)
			}
			if n < 1 {
				fmt.Fprintf(os.Stderr, "Error: --max-iterations must be at least 1\n")
				os.Exit(1)
			}
			opts.MaxIterations = n
		case strings.HasPrefix(arg, "-n="):
			val := strings.TrimPrefix(arg, "-n=")
			n, err := strconv.Atoi(val)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: invalid value for -n: %s\n", val)
				os.Exit(1)
			}
			if n < 1 {
				fmt.Fprintf(os.Stderr, "Error: -n must be at least 1\n")
				os.Exit(1)
			}
			opts.MaxIterations = n
		case strings.HasPrefix(arg, "-"):
			// Unknown flag
			fmt.Fprintf(os.Stderr, "Error: unknown flag: %s\n", arg)
			fmt.Fprintf(os.Stderr, "Run 'chief --help' for usage.\n")
			os.Exit(1)
		default:
			// Positional argument: PRD name or path
			if strings.HasSuffix(arg, ".json") || strings.HasSuffix(arg, "/") {
				opts.PRDPath = arg
			} else {
				// Treat as PRD name
				opts.PRDPath = fmt.Sprintf(".chief/prds/%s/prd.json", arg)
			}
		}
	}

	return opts
}

func runInit() {
	opts := cmd.InitOptions{}

	// Parse arguments: chief init [name] [context...]
	if len(os.Args) > 2 {
		opts.Name = os.Args[2]
	}
	if len(os.Args) > 3 {
		opts.Context = strings.Join(os.Args[3:], " ")
	}

	if err := cmd.RunInit(opts); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runEdit() {
	opts := cmd.EditOptions{}

	// Parse arguments: chief edit [name] [--merge] [--force]
	for i := 2; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch arg {
		case "--merge":
			opts.Merge = true
		case "--force":
			opts.Force = true
		default:
			// If not a flag, treat as PRD name (first non-flag arg)
			if opts.Name == "" && !strings.HasPrefix(arg, "-") {
				opts.Name = arg
			}
		}
	}

	if err := cmd.RunEdit(opts); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runStatus() {
	opts := cmd.StatusOptions{}

	// Parse arguments: chief status [name]
	if len(os.Args) > 2 && !strings.HasPrefix(os.Args[2], "-") {
		opts.Name = os.Args[2]
	}

	if err := cmd.RunStatus(opts); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runList() {
	opts := cmd.ListOptions{}

	if err := cmd.RunList(opts); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runTUIWithOptions(opts *TUIOptions) {
	prdPath := opts.PRDPath
	prdDir := filepath.Dir(prdPath)

	// Check if prd.md is newer than prd.json and run conversion if needed
	needsConvert, err := prd.NeedsConversion(prdDir)
	if err != nil {
		fmt.Printf("Warning: failed to check conversion status: %v\n", err)
	} else if needsConvert {
		fmt.Println("prd.md is newer than prd.json, running conversion...")
		convertOpts := prd.ConvertOptions{
			PRDDir: prdDir,
			Merge:  opts.Merge,
			Force:  opts.Force,
		}
		if err := prd.Convert(convertOpts); err != nil {
			fmt.Printf("Error converting PRD: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Conversion complete.")
	}

	app, err := tui.NewAppWithOptions(prdPath, opts.MaxIterations)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("\nTo get started, create a PRD first:")
		fmt.Println("  chief init              # Create default PRD")
		fmt.Println("  chief init <name>       # Create named PRD")
		os.Exit(1)
	}

	// Set verbose mode if requested
	if opts.Verbose {
		app.SetVerbose(true)
	}

	// Initialize sound notifier (unless disabled)
	if !opts.NoSound {
		notifier, err := notify.GetNotifier()
		if err != nil {
			// Log warning but don't crash - audio is optional
			log.Printf("Warning: audio initialization failed: %v", err)
		} else {
			// Set completion callback to play sound when any PRD completes
			app.SetCompletionCallback(func(prdName string) {
				notifier.PlayCompletion()
			})
		}
	}

	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println(`Chief - Autonomous PRD Agent

Usage:
  chief [options] [<name>|<path/to/prd.json>]
  chief <command> [arguments]

Commands:
  init [name] [context]     Create a new PRD interactively
  edit [name] [options]     Edit an existing PRD interactively
  status [name]             Show progress for a PRD (default: main)
  list                      List all PRDs with progress
  help                      Show this help message

Global Options:
  --max-iterations N, -n N  Set maximum iterations (default: 10)
  --no-sound                Disable completion sound notifications
  --verbose                 Show raw Claude output in log
  --merge                   Auto-merge progress on conversion conflicts
  --force                   Auto-overwrite on conversion conflicts
  --help, -h                Show this help message
  --version, -v             Show version number

Edit Options:
  --merge                   Auto-merge progress on conversion conflicts
  --force                   Auto-overwrite on conversion conflicts

Positional Arguments:
  <name>                    PRD name (loads .chief/prds/<name>/prd.json)
  <path/to/prd.json>        Direct path to a prd.json file

Examples:
  chief                     Launch TUI with default PRD (.chief/prds/main/)
  chief auth                Launch TUI with named PRD (.chief/prds/auth/)
  chief ./my-prd.json       Launch TUI with specific PRD file
  chief -n 20               Launch with 20 max iterations
  chief --max-iterations=5 auth
                            Launch auth PRD with 5 max iterations
  chief --no-sound          Launch TUI without audio notifications
  chief --verbose           Launch with raw Claude output visible
  chief init                Create PRD in .chief/prds/main/
  chief init auth           Create PRD in .chief/prds/auth/
  chief init auth "JWT authentication for REST API"
                            Create PRD with context hint
  chief edit                Edit PRD in .chief/prds/main/
  chief edit auth           Edit PRD in .chief/prds/auth/
  chief edit auth --merge   Edit and auto-merge progress
  chief status              Show progress for default PRD
  chief status auth         Show progress for auth PRD
  chief list                List all PRDs with progress
  chief --version           Show version number`)
}
