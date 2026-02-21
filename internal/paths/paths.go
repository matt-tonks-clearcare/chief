package paths

import (
	"os"
	"path/filepath"
)

// homeDir returns the user's home directory, panicking if it can't be resolved.
var homeDir = func() string {
	h, err := os.UserHomeDir()
	if err != nil {
		panic("cannot resolve home directory: " + err.Error())
	}
	return h
}

// SetHomeDir overrides the home directory used by all path functions.
// Intended for testing. Returns a restore function.
func SetHomeDir(dir string) func() {
	old := homeDir
	homeDir = func() string { return dir }
	return func() { homeDir = old }
}

// projectID returns the directory name used to identify a project.
func projectID(projectDir string) string {
	return filepath.Base(projectDir)
}

// ChiefDir returns ~/.chief/projects/<project-dir-name>/
func ChiefDir(projectDir string) string {
	return filepath.Join(homeDir(), ".chief", "projects", projectID(projectDir))
}

// PRDsDir returns ~/.chief/projects/<project-dir-name>/prds/
func PRDsDir(projectDir string) string {
	return filepath.Join(ChiefDir(projectDir), "prds")
}

// PRDDir returns ~/.chief/projects/<project-dir-name>/prds/<name>/
func PRDDir(projectDir string, name string) string {
	return filepath.Join(PRDsDir(projectDir), name)
}

// PRDPath returns ~/.chief/projects/<project-dir-name>/prds/<name>/prd.json
func PRDPath(projectDir string, name string) string {
	return filepath.Join(PRDDir(projectDir, name), "prd.json")
}

// ConfigPath returns ~/.chief/projects/<project-dir-name>/config.yaml
func ConfigPath(projectDir string) string {
	return filepath.Join(ChiefDir(projectDir), "config.yaml")
}

// WorktreeDir returns ~/.chief/projects/<project-dir-name>/worktrees/<name>/
func WorktreeDir(projectDir string, name string) string {
	return filepath.Join(ChiefDir(projectDir), "worktrees", name)
}

// WorktreesDir returns ~/.chief/projects/<project-dir-name>/worktrees/
func WorktreesDir(projectDir string) string {
	return filepath.Join(ChiefDir(projectDir), "worktrees")
}

// ContextDir returns ~/.chief/projects/<project-dir-name>/context/
func ContextDir(projectDir string) string {
	return filepath.Join(ChiefDir(projectDir), "context")
}
