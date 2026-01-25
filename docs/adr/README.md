# Architecture Decision Records

This directory contains Architecture Decision Records (ADRs) for the Chief project.

## What is an ADR?

An ADR is a document that captures an important architectural decision made along with its context and consequences. ADRs help new team members understand why certain decisions were made.

## ADR Index

| Number | Title | Status |
|--------|-------|--------|
| [0001](0001-use-bubble-tea-for-tui.md) | Use Bubble Tea for TUI | Accepted |
| [0002](0002-stream-json-parsing.md) | Parse Claude's Stream-JSON Output | Accepted |
| [0003](0003-parallel-prd-execution.md) | Support Parallel PRD Execution | Accepted |
| [0004](0004-embedded-prompts.md) | Embed Prompts in Binary | Accepted |
| [0005](0005-prd-file-watching.md) | Watch PRD File for Changes | Accepted |
| [0006](0006-audio-notifications.md) | Audio Completion Notifications | Accepted |

## Creating a New ADR

1. Copy the template below
2. Number sequentially (e.g., `0007-decision-name.md`)
3. Fill in the sections
4. Add to the index above

## ADR Template

```markdown
# ADR-NNNN: Title

## Status

Proposed | Accepted | Deprecated | Superseded

## Context

What is the issue we're addressing?

## Decision

What is the change we're proposing?

## Rationale

Why is this the best choice?

## Consequences

### Positive
- Benefits

### Negative
- Drawbacks

## References

- Links to relevant resources
```
