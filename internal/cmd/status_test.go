package cmd

import (
	"os"
	"testing"

	"github.com/minicodemonkey/chief/internal/paths"
)

func TestRunStatusWithValidPRD(t *testing.T) {
	tmpHome := t.TempDir()
	restore := paths.SetHomeDir(tmpHome)
	defer restore()

	tmpDir := t.TempDir()

	prdDir := paths.PRDDir(tmpDir, "test")
	if err := os.MkdirAll(prdDir, 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	prdJSON := `{
  "project": "Test Project",
  "description": "Test description",
  "userStories": [
    {"id": "US-001", "title": "Story 1", "passes": true, "priority": 1},
    {"id": "US-002", "title": "Story 2", "passes": false, "priority": 2},
    {"id": "US-003", "title": "Story 3", "passes": false, "inProgress": true, "priority": 3}
  ]
}`
	if err := os.WriteFile(paths.PRDPath(tmpDir, "test"), []byte(prdJSON), 0644); err != nil {
		t.Fatalf("Failed to create prd.json: %v", err)
	}

	opts := StatusOptions{
		Name:    "test",
		BaseDir: tmpDir,
	}

	err := RunStatus(opts)
	if err != nil {
		t.Errorf("RunStatus() returned error: %v", err)
	}
}

func TestRunStatusWithDefaultName(t *testing.T) {
	tmpHome := t.TempDir()
	restore := paths.SetHomeDir(tmpHome)
	defer restore()

	tmpDir := t.TempDir()

	prdDir := paths.PRDDir(tmpDir, "main")
	if err := os.MkdirAll(prdDir, 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	prdJSON := `{
  "project": "Main Project",
  "userStories": []
}`
	if err := os.WriteFile(paths.PRDPath(tmpDir, "main"), []byte(prdJSON), 0644); err != nil {
		t.Fatalf("Failed to create prd.json: %v", err)
	}

	opts := StatusOptions{
		Name:    "",
		BaseDir: tmpDir,
	}

	err := RunStatus(opts)
	if err != nil {
		t.Errorf("RunStatus() with default name returned error: %v", err)
	}
}

func TestRunStatusWithMissingPRD(t *testing.T) {
	tmpHome := t.TempDir()
	restore := paths.SetHomeDir(tmpHome)
	defer restore()

	tmpDir := t.TempDir()

	opts := StatusOptions{
		Name:    "nonexistent",
		BaseDir: tmpDir,
	}

	err := RunStatus(opts)
	if err == nil {
		t.Error("Expected error for missing PRD")
	}
}

func TestRunListWithNoPRDs(t *testing.T) {
	tmpHome := t.TempDir()
	restore := paths.SetHomeDir(tmpHome)
	defer restore()

	tmpDir := t.TempDir()

	opts := ListOptions{
		BaseDir: tmpDir,
	}

	err := RunList(opts)
	if err != nil {
		t.Errorf("RunList() returned error: %v", err)
	}
}

func TestRunListWithPRDs(t *testing.T) {
	tmpHome := t.TempDir()
	restore := paths.SetHomeDir(tmpHome)
	defer restore()

	tmpDir := t.TempDir()

	prds := []struct {
		name    string
		project string
		stories string
	}{
		{
			"auth",
			"Authentication",
			`[{"id": "US-001", "title": "Login", "passes": true, "priority": 1},
			 {"id": "US-002", "title": "Logout", "passes": false, "priority": 2}]`,
		},
		{
			"api",
			"API Service",
			`[{"id": "US-001", "title": "Endpoints", "passes": true, "priority": 1},
			 {"id": "US-002", "title": "Auth", "passes": true, "priority": 2},
			 {"id": "US-003", "title": "Rate limiting", "passes": true, "priority": 3}]`,
		},
	}

	for _, p := range prds {
		prdDir := paths.PRDDir(tmpDir, p.name)
		if err := os.MkdirAll(prdDir, 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}

		prdJSON := `{"project": "` + p.project + `", "userStories": ` + p.stories + `}`
		if err := os.WriteFile(paths.PRDPath(tmpDir, p.name), []byte(prdJSON), 0644); err != nil {
			t.Fatalf("Failed to create prd.json: %v", err)
		}
	}

	opts := ListOptions{
		BaseDir: tmpDir,
	}

	err := RunList(opts)
	if err != nil {
		t.Errorf("RunList() returned error: %v", err)
	}
}

func TestRunListSkipsInvalidPRDs(t *testing.T) {
	tmpHome := t.TempDir()
	restore := paths.SetHomeDir(tmpHome)
	defer restore()

	tmpDir := t.TempDir()

	validDir := paths.PRDDir(tmpDir, "valid")
	if err := os.MkdirAll(validDir, 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}
	validJSON := `{"project": "Valid", "userStories": []}`
	if err := os.WriteFile(paths.PRDPath(tmpDir, "valid"), []byte(validJSON), 0644); err != nil {
		t.Fatalf("Failed to create prd.json: %v", err)
	}

	invalidDir := paths.PRDDir(tmpDir, "invalid")
	if err := os.MkdirAll(invalidDir, 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	badJsonDir := paths.PRDDir(tmpDir, "badjson")
	if err := os.MkdirAll(badJsonDir, 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}
	if err := os.WriteFile(paths.PRDPath(tmpDir, "badjson"), []byte("not json"), 0644); err != nil {
		t.Fatalf("Failed to create prd.json: %v", err)
	}

	opts := ListOptions{
		BaseDir: tmpDir,
	}

	err := RunList(opts)
	if err != nil {
		t.Errorf("RunList() returned error: %v", err)
	}
}

func TestRunStatusAllComplete(t *testing.T) {
	tmpHome := t.TempDir()
	restore := paths.SetHomeDir(tmpHome)
	defer restore()

	tmpDir := t.TempDir()

	prdDir := paths.PRDDir(tmpDir, "done")
	if err := os.MkdirAll(prdDir, 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	prdJSON := `{
  "project": "Complete Project",
  "userStories": [
    {"id": "US-001", "title": "Story 1", "passes": true, "priority": 1},
    {"id": "US-002", "title": "Story 2", "passes": true, "priority": 2}
  ]
}`
	if err := os.WriteFile(paths.PRDPath(tmpDir, "done"), []byte(prdJSON), 0644); err != nil {
		t.Fatalf("Failed to create prd.json: %v", err)
	}

	opts := StatusOptions{
		Name:    "done",
		BaseDir: tmpDir,
	}

	err := RunStatus(opts)
	if err != nil {
		t.Errorf("RunStatus() returned error: %v", err)
	}
}

func TestRunStatusEmptyPRD(t *testing.T) {
	tmpHome := t.TempDir()
	restore := paths.SetHomeDir(tmpHome)
	defer restore()

	tmpDir := t.TempDir()

	prdDir := paths.PRDDir(tmpDir, "empty")
	if err := os.MkdirAll(prdDir, 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	prdJSON := `{"project": "Empty Project", "userStories": []}`
	if err := os.WriteFile(paths.PRDPath(tmpDir, "empty"), []byte(prdJSON), 0644); err != nil {
		t.Fatalf("Failed to create prd.json: %v", err)
	}

	opts := StatusOptions{
		Name:    "empty",
		BaseDir: tmpDir,
	}

	err := RunStatus(opts)
	if err != nil {
		t.Errorf("RunStatus() returned error: %v", err)
	}
}
