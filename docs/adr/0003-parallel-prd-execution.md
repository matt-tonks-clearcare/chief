# ADR-0003: Support Parallel PRD Execution

## Status

Accepted

## Context

Users may want to work on multiple features simultaneously. Each PRD represents an independent work stream with its own Claude process, progress, and logs.

Options considered:

1. **Single PRD at a time**: Simpler but limits productivity
2. **Tab-based switching**: Visual complexity, still one active loop
3. **Parallel execution with manager**: Multiple loops running independently

## Decision

Implement a `Manager` that can run multiple `Loop` instances concurrently, each in its own goroutine with its own Claude process.

## Design

```
┌─────────────────────────────────────────────┐
│                  Manager                     │
│  ┌──────────────────────────────────────┐   │
│  │        instances map[string]*Instance │   │
│  │                                       │   │
│  │  "auth" ──▶ Instance {               │   │
│  │              Loop: *Loop              │   │
│  │              State: Running           │   │
│  │              ctx: context.Context     │   │
│  │            }                          │   │
│  │                                       │   │
│  │  "main" ──▶ Instance {               │   │
│  │              Loop: *Loop              │   │
│  │              State: Paused            │   │
│  │            }                          │   │
│  └──────────────────────────────────────┘   │
│                                              │
│  Events() <-chan ManagerEvent               │
│  ─────────────────────────────              │
│  Forwards events from all loops with        │
│  PRDName attached for routing               │
└─────────────────────────────────────────────┘
```

## Rationale

1. **Independence**: Each PRD has its own context for cancellation, ensuring clean shutdown.

2. **Unified Event Stream**: The manager forwards events from all loops through a single channel with PRD names attached, simplifying TUI event handling.

3. **State Tracking**: Each instance tracks its own state (Ready, Running, Paused, etc.), visible in the PRD picker.

4. **Resource Cleanup**: `StopAll()` uses `sync.WaitGroup` to ensure all loops terminate before exit.

## Consequences

### Positive

- Work on multiple features simultaneously
- Each PRD's state visible in the picker
- Completion callbacks fire for any PRD (audio notification)
- Clean separation between viewing and running (can view one PRD while another runs)

### Negative

- Higher resource usage (multiple Claude processes)
- More complex state management
- Need to route events by PRD name

## Thread Safety

- `Manager.instances` protected by `sync.RWMutex`
- Each `LoopInstance` has its own `sync.Mutex`
- Events sent through buffered channels (capacity 100)

## References

- `internal/loop/manager.go` - Manager implementation
- `internal/loop/manager_test.go` - Concurrent access tests
