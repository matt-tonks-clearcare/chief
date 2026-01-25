# ADR-0004: Embed Prompts in Binary

## Status

Accepted

## Context

Chief uses several prompts to instruct Claude:

1. **Agent prompt**: How to work through user stories
2. **Init prompt**: How to create a new PRD
3. **Edit prompt**: How to edit an existing PRD
4. **Convert prompt**: How to convert prd.md to prd.json

Options for storing prompts:

1. **External files**: Flexible but requires distribution
2. **Embedded in code**: Hard to maintain, escaping issues
3. **Go embed directive**: Embedded at compile time, no external dependencies

## Decision

Use Go's `//go:embed` directive to embed prompt files in the binary.

## Implementation

```go
// embed/embed.go
package embed

import (
    _ "embed"
    "strings"
)

//go:embed prompt.txt
var promptTemplate string

// GetPrompt returns the agent prompt with placeholders substituted.
func GetPrompt(prdPath string) string {
    return strings.ReplaceAll(promptTemplate, "{{PRD_PATH}}", prdPath)
}
```

Prompts use `{{PLACEHOLDER}}` syntax for variable substitution.

## Rationale

1. **Single Binary**: No external files to manage or distribute.

2. **Version Consistency**: Prompts are tied to the binary version.

3. **Easy Maintenance**: Prompts are plain text files, easy to edit.

4. **Type Safety**: Compile fails if embed file is missing.

## Consequences

### Positive

- Self-contained binary
- Prompts versioned with code
- No runtime file loading errors
- Works with goreleaser cross-compilation

### Negative

- Requires rebuild to change prompts
- Increases binary size (minimal, ~10KB total)
- Can't customize prompts per-installation

## Prompt Files

| File | Purpose | Placeholders |
|------|---------|--------------|
| `prompt.txt` | Agent loop instructions | `{{PRD_PATH}}` |
| `init_prompt.txt` | PRD creation guidance | `{{CONTEXT}}` |
| `edit_prompt.txt` | PRD editing guidance | None |
| `convert_prompt.txt` | prd.md to prd.json conversion | None |

## References

- [Go embed documentation](https://pkg.go.dev/embed)
- `embed/` directory - All embedded files
- `embed/embed.go` - Embed accessors
