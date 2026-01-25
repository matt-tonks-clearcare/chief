# ADR-0002: Parse Claude's Stream-JSON Output

## Status

Accepted

## Context

Chief needs to understand what Claude is doing in real-time to display activity in the TUI. Claude Code CLI offers multiple output formats:

- **text**: Plain text output, human-readable but hard to parse
- **json**: Complete JSON at the end, no streaming
- **stream-json**: Line-by-line JSON events during execution

## Decision

Use Claude's `--output-format stream-json --verbose` flags to get real-time events.

## Rationale

1. **Real-time Updates**: Stream-JSON emits events as Claude works, enabling live activity display.

2. **Structured Data**: Each line is a valid JSON object with a `type` field, making parsing reliable.

3. **Tool Visibility**: Events include tool invocations (`tool_use`) and results (`tool_result`), allowing Chief to show what Claude is doing.

4. **Completion Detection**: Assistant messages contain text that can be scanned for completion markers like `<chief-complete/>`.

## Implementation

The parser in `internal/loop/parser.go` handles three main message types:

```go
switch msg.Type {
case "system":    // Init events
case "assistant": // Text and tool_use
case "user":      // Tool results
}
```

Events are emitted to a channel and consumed by the TUI for display.

## Consequences

### Positive

- Live activity updates in the TUI
- Can detect story starts, tool usage, and completion
- Raw output logged to `claude.log` for debugging

### Negative

- Parser must handle malformed or unexpected JSON gracefully
- Large tool outputs (e.g., file contents) can produce very long lines
- Format may change between Claude CLI versions

## Mitigations

- Use `bufio.Scanner` with 1MB buffer for large lines
- Parser returns `nil` for unparseable lines (graceful degradation)
- Version-specific parsing could be added if needed

## References

- Claude Code CLI documentation
- `internal/loop/parser.go` - Event parser implementation
- `internal/loop/parser_test.go` - Parser tests with real examples
