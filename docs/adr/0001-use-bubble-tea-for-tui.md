# ADR-0001: Use Bubble Tea for TUI

## Status

Accepted

## Context

Chief needs a terminal user interface (TUI) to display PRD progress, show Claude's activity, and provide keyboard controls for the agent loop. Several TUI frameworks are available for Go:

- **Bubble Tea** (charmbracelet/bubbletea) - Elm-inspired framework
- **tview** - Terminal UI library based on tcell
- **gocui** - Minimalist Go library for terminal UI
- **termui** - Dashboard/chart widgets

## Decision

We chose **Bubble Tea** with **Lip Gloss** for styling.

## Rationale

1. **Elm Architecture**: Bubble Tea's model-update-view pattern provides predictable state management, essential for complex UIs with multiple views and async events.

2. **Composability**: Components like viewports, log viewers, and modals can be built as separate models and composed together.

3. **Ecosystem**: Lip Gloss provides elegant styling, Bubbles provides pre-built components, and the charmbracelet ecosystem is actively maintained.

4. **Async Handling**: The `tea.Cmd` pattern handles async operations (file watching, loop events) cleanly without manual goroutine management.

5. **Testing**: The functional architecture makes state transitions easy to test without mocking terminal output.

## Consequences

### Positive

- Clean separation of concerns (model, update, view)
- Easy to add new views (log view, picker, help overlay)
- Built-in support for alt-screen, mouse events, and window resize
- Active community and good documentation

### Negative

- Learning curve for developers unfamiliar with Elm architecture
- Some boilerplate for message types and update handlers
- Less suitable for highly dynamic layouts (though Chief's layout is mostly static)

## References

- [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- [Bubbles](https://github.com/charmbracelet/bubbles)
