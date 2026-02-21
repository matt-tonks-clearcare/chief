// Package context provides automatic context file loading for PRD generation.
// It scans global (~/.claude/context/) and project-level (~/.chief/projects/<project>/context/)
// directories for .md files and returns their concatenated content.
package context

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/minicodemonkey/chief/internal/paths"
)

// LoadContextFiles scans global (~/.claude/context/) and project-level
// (~/.chief/projects/<project>/context/) directories for .md files, reads them, and returns
// their concatenated content. Returns empty string if neither directory
// exists or contains .md files.
func LoadContextFiles(baseDir string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Can't resolve home dir; skip global context
		homeDir = ""
	}
	return loadContextFilesWithHome(baseDir, homeDir)
}

// loadContextFilesWithHome is the testable core that accepts an explicit home directory.
func loadContextFilesWithHome(baseDir, homeDir string) (string, error) {
	var parts []string

	// Global directory: ~/.claude/context/
	if homeDir != "" {
		globalDir := filepath.Join(homeDir, ".claude", "context")
		content, err := loadMarkdownFiles(globalDir)
		if err != nil {
			return "", err
		}
		if content != "" {
			parts = append(parts, content)
		}
	}

	// Project-level directory: ~/.chief/projects/<project>/context/
	projectDir := paths.ContextDir(baseDir)
	content, err := loadMarkdownFiles(projectDir)
	if err != nil {
		return "", err
	}
	if content != "" {
		parts = append(parts, content)
	}

	return strings.Join(parts, "\n\n"), nil
}

// loadMarkdownFiles reads all .md files from dir, sorted by filename, and
// returns their concatenated content. Returns ("", nil) if dir doesn't exist.
func loadMarkdownFiles(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	var mdFiles []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(strings.ToLower(e.Name()), ".md") {
			mdFiles = append(mdFiles, e.Name())
		}
	}
	sort.Strings(mdFiles)

	var sections []string
	for _, name := range mdFiles {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return "", err
		}
		sections = append(sections, string(data))
	}

	return strings.Join(sections, "\n\n---\n\n"), nil
}
