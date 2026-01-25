# ADR-0005: Watch PRD File for Changes

## Status

Accepted

## Context

As Claude works, it updates `prd.json` to mark stories as `inProgress` or `passes: true`. The TUI needs to reflect these changes in real-time without manual refresh.

Options:

1. **Polling**: Check file periodically (simple but inefficient)
2. **After each iteration**: Only update when Claude finishes (delayed)
3. **Filesystem watching**: React to file changes immediately

## Decision

Use `fsnotify` to watch `prd.json` for filesystem changes and reload when modified.

## Implementation

```go
// internal/prd/watcher.go
type Watcher struct {
    path    string
    watcher *fsnotify.Watcher
    events  chan WatcherEvent
    lastPRD *PRD
}
```

The watcher:
1. Monitors write and create events
2. Reloads the PRD on change
3. Only emits events when status fields actually change
4. Handles file removal gracefully (attempts re-watch)

## Rationale

1. **Immediate Updates**: Story status changes appear in the TUI instantly.

2. **Efficiency**: Only reloads when the file actually changes.

3. **Filtering**: `hasStatusChanged()` prevents unnecessary UI updates for non-status changes (e.g., description edits).

4. **Graceful Degradation**: If watching fails, the app continues without real-time updates.

## Consequences

### Positive

- Real-time story progress in the TUI
- Works with external editors modifying prd.json
- Efficient (event-driven, not polling)

### Negative

- Platform differences in fsnotify behavior
- File removal requires re-watch attempt
- Additional goroutine per watched file

## Status Change Detection

Only these field changes trigger UI updates:
- `passes`: Story completion status
- `inProgress`: Story being worked on

Other changes (title, description, criteria) are ignored to reduce noise.

## References

- [fsnotify](https://github.com/fsnotify/fsnotify)
- `internal/prd/watcher.go` - Watcher implementation
- `internal/prd/watcher_test.go` - Integration tests
