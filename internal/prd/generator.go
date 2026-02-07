package prd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/minicodemonkey/chief/embed"
)

// spinner frames for the loading indicator
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// ConvertOptions contains configuration for PRD conversion.
type ConvertOptions struct {
	PRDDir string // Directory containing prd.md
	Merge  bool   // Auto-merge progress on conversion conflicts
	Force  bool   // Auto-overwrite on conversion conflicts
}

// ProgressConflictChoice represents the user's choice when a progress conflict is detected.
type ProgressConflictChoice int

const (
	ChoiceMerge     ProgressConflictChoice = iota // Keep status for matching story IDs
	ChoiceOverwrite                               // Discard all progress
	ChoiceCancel                                  // Cancel conversion
)

// Convert converts prd.md to prd.json using Claude one-shot mode.
// Claude is responsible for writing the prd.json file directly.
// This function is called:
// - After chief new (new PRD creation)
// - After chief edit (PRD modification)
// - Before chief run if prd.md is newer than prd.json
//
// Progress protection:
// - If prd.json has progress (passes: true or inProgress: true) and prd.md changed:
//   - opts.Merge: auto-merge, preserving status for matching story IDs
//   - opts.Force: auto-overwrite, discarding all progress
//   - Neither: prompt the user with Merge/Overwrite/Cancel options
func Convert(opts ConvertOptions) error {
	prdMdPath := filepath.Join(opts.PRDDir, "prd.md")
	prdJsonPath := filepath.Join(opts.PRDDir, "prd.json")

	// Check if prd.md exists
	if _, err := os.Stat(prdMdPath); os.IsNotExist(err) {
		return fmt.Errorf("prd.md not found in %s", opts.PRDDir)
	}

	// Resolve absolute path so the prompt can specify exact file locations
	absPRDDir, err := filepath.Abs(opts.PRDDir)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	// Check for existing progress before conversion
	var existingPRD *PRD
	hasProgress := false
	if existing, err := LoadPRD(prdJsonPath); err == nil {
		existingPRD = existing
		hasProgress = HasProgress(existing)
	}

	// Run Claude to convert prd.md and write prd.json directly
	if err := runClaudeConversion(absPRDDir); err != nil {
		return err
	}

	// Validate that Claude wrote a valid prd.json
	newPRD, err := loadAndValidateConvertedPRD(prdJsonPath)
	if err != nil {
		// Retry once: ask Claude to fix the invalid JSON
		fmt.Println("Conversion produced invalid JSON, retrying...")
		if retryErr := runClaudeJSONFix(absPRDDir, err); retryErr != nil {
			return fmt.Errorf("conversion retry failed: %w", retryErr)
		}

		newPRD, err = loadAndValidateConvertedPRD(prdJsonPath)
		if err != nil {
			return fmt.Errorf("conversion produced invalid JSON after retry: %w", err)
		}
	}

	// Re-save through Go's JSON encoder to guarantee proper escaping and formatting
	normalizedContent, err := json.MarshalIndent(newPRD, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal PRD: %w", err)
	}

	// Handle progress protection if existing prd.json has progress
	if hasProgress && existingPRD != nil {
		choice := ChoiceOverwrite // Default to overwrite if no progress

		if opts.Merge {
			choice = ChoiceMerge
		} else if opts.Force {
			choice = ChoiceOverwrite
		} else {
			// Prompt user for choice
			var promptErr error
			choice, promptErr = promptProgressConflict(existingPRD, newPRD)
			if promptErr != nil {
				return fmt.Errorf("failed to prompt for choice: %w", promptErr)
			}
		}

		switch choice {
		case ChoiceCancel:
			return fmt.Errorf("conversion cancelled by user")
		case ChoiceMerge:
			// Merge progress from existing PRD into new PRD
			MergeProgress(existingPRD, newPRD)
			// Re-marshal with merged progress
			mergedContent, err := json.MarshalIndent(newPRD, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal merged PRD: %w", err)
			}
			normalizedContent = mergedContent
		case ChoiceOverwrite:
			// Use the new PRD as-is (no progress)
		}
	}

	// Write the final normalized prd.json
	if err := os.WriteFile(prdJsonPath, append(normalizedContent, '\n'), 0644); err != nil {
		return fmt.Errorf("failed to write prd.json: %w", err)
	}

	return nil
}

// runClaudeConversion runs Claude one-shot to convert prd.md and write prd.json.
func runClaudeConversion(absPRDDir string) error {
	prompt := embed.GetConvertPrompt(absPRDDir)

	cmd := exec.Command("claude",
		"--dangerously-skip-permissions",
		"-p", prompt,
	)
	cmd.Dir = absPRDDir

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start Claude: %w", err)
	}

	return waitWithSpinner(cmd, "Converting prd.md to prd.json...", &stderr)
}

// runClaudeJSONFix asks Claude to fix an invalid prd.json file.
func runClaudeJSONFix(absPRDDir string, validationErr error) error {
	fixPrompt := fmt.Sprintf(
		"The file at %s/prd.json contains invalid JSON. The error is: %s\n\n"+
			"Read the file, fix the JSON (pay special attention to escaping double quotes inside string values with backslashes), "+
			"and write the corrected JSON back to %s/prd.json.",
		absPRDDir, validationErr.Error(), absPRDDir,
	)

	cmd := exec.Command("claude",
		"--dangerously-skip-permissions",
		"-p", fixPrompt,
	)
	cmd.Dir = absPRDDir

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start Claude: %w", err)
	}

	return waitWithSpinner(cmd, "Fixing prd.json...", &stderr)
}

// loadAndValidateConvertedPRD loads prd.json and validates it can be parsed as a PRD.
func loadAndValidateConvertedPRD(prdJsonPath string) (*PRD, error) {
	prd, err := LoadPRD(prdJsonPath)
	if err != nil {
		return nil, err
	}
	if prd.Project == "" {
		return nil, fmt.Errorf("prd.json missing required 'project' field")
	}
	if len(prd.UserStories) == 0 {
		return nil, fmt.Errorf("prd.json has no user stories")
	}
	return prd, nil
}

// waitWithSpinner runs a spinner while waiting for a command to finish.
func waitWithSpinner(cmd *exec.Cmd, message string, stderr *bytes.Buffer) error {
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	frame := 0
	ticker := time.NewTicker(80 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case err := <-done:
			fmt.Print("\r\033[K")
			if err != nil {
				return fmt.Errorf("Claude failed: %s", stderr.String())
			}
			return nil
		case <-ticker.C:
			fmt.Printf("\r%s %s", spinnerFrames[frame%len(spinnerFrames)], message)
			frame++
		}
	}
}

// NeedsConversion checks if prd.md is newer than prd.json, indicating conversion is needed.
// Returns true if:
// - prd.md exists and prd.json does not exist
// - prd.md exists and is newer than prd.json
// Returns false if:
// - prd.md does not exist
// - prd.json is newer than or same age as prd.md
func NeedsConversion(prdDir string) (bool, error) {
	prdMdPath := filepath.Join(prdDir, "prd.md")
	prdJsonPath := filepath.Join(prdDir, "prd.json")

	// Check if prd.md exists
	mdInfo, err := os.Stat(prdMdPath)
	if os.IsNotExist(err) {
		// No prd.md, no conversion needed
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to stat prd.md: %w", err)
	}

	// Check if prd.json exists
	jsonInfo, err := os.Stat(prdJsonPath)
	if os.IsNotExist(err) {
		// prd.md exists but prd.json doesn't - needs conversion
		return true, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to stat prd.json: %w", err)
	}

	// Both exist - compare modification times
	return mdInfo.ModTime().After(jsonInfo.ModTime()), nil
}

// cleanJSONOutput removes markdown code blocks and trims whitespace from Claude's output.
func cleanJSONOutput(output string) string {
	output = strings.TrimSpace(output)

	// Remove markdown code blocks if present
	if strings.HasPrefix(output, "```json") {
		output = strings.TrimPrefix(output, "```json")
	} else if strings.HasPrefix(output, "```") {
		output = strings.TrimPrefix(output, "```")
	}

	if strings.HasSuffix(output, "```") {
		output = strings.TrimSuffix(output, "```")
	}

	return strings.TrimSpace(output)
}

// validateJSON checks if the given string is valid JSON.
func validateJSON(content string) error {
	var js json.RawMessage
	if err := json.Unmarshal([]byte(content), &js); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	return nil
}

// HasProgress checks if the PRD has any progress (passes: true or inProgress: true).
func HasProgress(prd *PRD) bool {
	if prd == nil {
		return false
	}
	for _, story := range prd.UserStories {
		if story.Passes || story.InProgress {
			return true
		}
	}
	return false
}

// MergeProgress merges progress from the old PRD into the new PRD.
// For stories with matching IDs, it preserves the Passes and InProgress status.
// New stories (in newPRD but not in oldPRD) are added without progress.
// Removed stories (in oldPRD but not in newPRD) are dropped.
func MergeProgress(oldPRD, newPRD *PRD) {
	if oldPRD == nil || newPRD == nil {
		return
	}

	// Create a map of old story statuses by ID
	oldStatus := make(map[string]struct {
		passes     bool
		inProgress bool
	})
	for _, story := range oldPRD.UserStories {
		oldStatus[story.ID] = struct {
			passes     bool
			inProgress bool
		}{
			passes:     story.Passes,
			inProgress: story.InProgress,
		}
	}

	// Apply old status to matching stories in new PRD
	for i := range newPRD.UserStories {
		if status, exists := oldStatus[newPRD.UserStories[i].ID]; exists {
			newPRD.UserStories[i].Passes = status.passes
			newPRD.UserStories[i].InProgress = status.inProgress
		}
	}
}

// promptProgressConflict prompts the user to choose how to handle a progress conflict.
func promptProgressConflict(oldPRD, newPRD *PRD) (ProgressConflictChoice, error) {
	// Count stories with progress
	progressCount := 0
	for _, story := range oldPRD.UserStories {
		if story.Passes || story.InProgress {
			progressCount++
		}
	}

	// Show warning
	fmt.Println()
	fmt.Printf("⚠️  Warning: prd.json has progress (%d stories with status)\n", progressCount)
	fmt.Println()
	fmt.Println("How would you like to proceed?")
	fmt.Println()
	fmt.Println("  [m] Merge  - Keep status for matching story IDs, add new stories, drop removed stories")
	fmt.Println("  [o] Overwrite - Discard all progress and use the new PRD")
	fmt.Println("  [c] Cancel - Cancel conversion and keep existing prd.json")
	fmt.Println()
	fmt.Print("Choice [m/o/c]: ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return ChoiceCancel, fmt.Errorf("failed to read input: %w", err)
	}

	input = strings.TrimSpace(strings.ToLower(input))
	switch input {
	case "m", "merge":
		return ChoiceMerge, nil
	case "o", "overwrite":
		return ChoiceOverwrite, nil
	case "c", "cancel", "":
		return ChoiceCancel, nil
	default:
		fmt.Printf("Invalid choice %q, cancelling conversion.\n", input)
		return ChoiceCancel, nil
	}
}
