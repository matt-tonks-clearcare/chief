package tui

import (
	"fmt"
	"testing"
	"unicode/utf8"

	"github.com/minicodemonkey/chief/internal/loop"
)

func TestRenderEntryWithBranchAndWorktree(t *testing.T) {
	p := &PRDPicker{
		basePath:   "/project",
		currentPRD: "",
		entries: []PRDEntry{
			{
				Name:        "auth",
				Completed:   8,
				Total:       8,
				LoopState:   loop.LoopStateComplete,
				Branch:      "chief/auth",
				WorktreeDir: "/project/.chief/worktrees/auth",
			},
		},
	}

	result := p.renderEntry(p.entries[0], false, 80)
	if result == "" {
		t.Fatal("expected non-empty render result")
	}
	// Should contain branch name
	if !containsText(result, "chief/auth") {
		t.Errorf("expected branch 'chief/auth' in output, got: %s", result)
	}
	// Should contain worktree path
	if !containsText(result, ".chief/worktrees/auth/") {
		t.Errorf("expected worktree path in output, got: %s", result)
	}
}

func TestRenderEntryNoBranch(t *testing.T) {
	p := &PRDPicker{
		basePath:   "/project",
		currentPRD: "",
		entries: []PRDEntry{
			{
				Name:      "auth",
				Completed: 3,
				Total:     8,
				LoopState: loop.LoopStateRunning,
				Iteration: 2,
				Branch:    "",
			},
		},
	}

	result := p.renderEntry(p.entries[0], false, 80)
	if result == "" {
		t.Fatal("expected non-empty render result")
	}
	// Should NOT contain branch brackets
	if containsText(result, "chief/") {
		t.Errorf("expected no branch in output when branch is empty, got: %s", result)
	}
}

func TestRenderEntryNoBranchShowsCurrentDirectoryWhenOthersHaveBranch(t *testing.T) {
	p := &PRDPicker{
		basePath:   "/project",
		currentPRD: "",
		entries: []PRDEntry{
			{
				Name:      "legacy",
				Completed: 5,
				Total:     5,
				LoopState: loop.LoopStateComplete,
				Branch:    "",
			},
			{
				Name:        "auth",
				Completed:   8,
				Total:       8,
				LoopState:   loop.LoopStateComplete,
				Branch:      "chief/auth",
				WorktreeDir: "/project/.chief/worktrees/auth",
			},
		},
	}

	// Render the entry without a branch — should show "(current directory)" since another has a branch
	result := p.renderEntry(p.entries[0], false, 80)
	if !containsText(result, "(current directory)") {
		t.Errorf("expected '(current directory)' for branchless entry when others have branches, got: %s", result)
	}
}

func TestRenderEntryNarrowTerminalOmitsBranchPath(t *testing.T) {
	p := &PRDPicker{
		basePath:   "/project",
		currentPRD: "",
		entries: []PRDEntry{
			{
				Name:        "auth",
				Completed:   8,
				Total:       8,
				LoopState:   loop.LoopStateComplete,
				Branch:      "chief/auth",
				WorktreeDir: "/project/.chief/worktrees/auth",
			},
		},
	}

	// Very narrow width — should not crash and should omit branch/path info
	result := p.renderEntry(p.entries[0], false, 35)
	if result == "" {
		t.Fatal("expected non-empty render result even at narrow width")
	}
	// At 35 chars wide, remaining space (35-32=3) is too small for branch info
	if containsText(result, "chief/auth") {
		t.Errorf("expected branch to be omitted at narrow width, got: %s", result)
	}
}

func TestFormatBranchPathFull(t *testing.T) {
	p := &PRDPicker{}

	result := p.formatBranchPath("chief/auth", ".chief/worktrees/auth/", 50)
	expected := "  chief/auth  .chief/worktrees/auth/"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestFormatBranchPathTruncatesPath(t *testing.T) {
	p := &PRDPicker{}

	result := p.formatBranchPath("chief/auth", ".chief/worktrees/auth/", 30)
	// Should contain branch but path should be truncated with …
	if !containsSubstring(result, "chief/auth") {
		t.Errorf("expected branch in truncated output, got: %s", result)
	}
	runeCount := utf8.RuneCountInString(result)
	if runeCount > 30 {
		t.Errorf("expected result to fit within 30 display chars, got %d: %s", runeCount, result)
	}
}

func TestFormatBranchPathTruncatesBranch(t *testing.T) {
	p := &PRDPicker{}

	// Very small width — only room for branch (truncated)
	result := p.formatBranchPath("chief/very-long-branch-name", ".chief/worktrees/auth/", 15)
	runeCount := utf8.RuneCountInString(result)
	if runeCount > 15 {
		t.Errorf("expected result to fit within 15 display chars, got %d: %s", runeCount, result)
	}
}

func TestWorktreeDisplayPathWithWorktree(t *testing.T) {
	p := &PRDPicker{basePath: "/project"}

	entry := PRDEntry{WorktreeDir: "/project/.chief/worktrees/auth"}
	result := p.worktreeDisplayPath(entry)
	if result != ".chief/worktrees/auth/" {
		t.Errorf("expected '.chief/worktrees/auth/', got %q", result)
	}
}

func TestWorktreeDisplayPathWithoutWorktree(t *testing.T) {
	p := &PRDPicker{basePath: "/project"}

	entry := PRDEntry{WorktreeDir: ""}
	result := p.worktreeDisplayPath(entry)
	if result != "(current directory)" {
		t.Errorf("expected '(current directory)', got %q", result)
	}
}

func TestHasAnyBranch(t *testing.T) {
	p := &PRDPicker{
		entries: []PRDEntry{
			{Name: "a", Branch: ""},
			{Name: "b", Branch: "chief/b"},
		},
	}
	if !p.hasAnyBranch() {
		t.Error("expected hasAnyBranch() to return true when one entry has a branch")
	}

	p2 := &PRDPicker{
		entries: []PRDEntry{
			{Name: "a", Branch: ""},
			{Name: "c", Branch: ""},
		},
	}
	if p2.hasAnyBranch() {
		t.Error("expected hasAnyBranch() to return false when no entries have branches")
	}
}

func TestRenderEntryWithLoadError(t *testing.T) {
	p := &PRDPicker{
		basePath: "/project",
		entries: []PRDEntry{
			{
				Name:        "broken",
				LoadError:   fmt.Errorf("parse error"),
				Branch:      "chief/broken",
				WorktreeDir: "/project/.chief/worktrees/broken",
			},
		},
	}

	result := p.renderEntry(p.entries[0], false, 80)
	// With load error, should show [error] but not branch/worktree info
	if !containsText(result, "error") {
		t.Errorf("expected [error] in output, got: %s", result)
	}
}

// containsText checks if rendered output contains a substring (ignoring ANSI codes).
func containsText(rendered, substr string) bool {
	// Strip ANSI escape sequences for comparison
	return containsSubstring(stripAnsi(rendered), substr)
}

// containsSubstring is a simple substring check.
func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && indexOf(s, substr) >= 0
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// stripAnsi removes ANSI escape codes from a string.
func stripAnsi(s string) string {
	var result []byte
	i := 0
	for i < len(s) {
		if s[i] == '\x1b' && i+1 < len(s) && s[i+1] == '[' {
			// Skip to end of escape sequence
			j := i + 2
			for j < len(s) && !((s[j] >= 'A' && s[j] <= 'Z') || (s[j] >= 'a' && s[j] <= 'z')) {
				j++
			}
			if j < len(s) {
				j++ // skip the final letter
			}
			i = j
		} else {
			result = append(result, s[i])
			i++
		}
	}
	return string(result)
}
