package context

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadContextFiles_NoDirsExist(t *testing.T) {
	tmpDir := t.TempDir()
	result, err := loadContextFilesWithHome(tmpDir, tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestLoadContextFiles_ProjectDirOnly(t *testing.T) {
	tmpHome := t.TempDir()
	tmpProject := t.TempDir()
	contextDir := filepath.Join(tmpProject, ".chief", "context")
	if err := os.MkdirAll(contextDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(contextDir, "platform.md"), []byte("Platform info"), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := loadContextFilesWithHome(tmpProject, tmpHome)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Platform info" {
		t.Errorf("expected 'Platform info', got %q", result)
	}
}

func TestLoadContextFiles_GlobalDirOnly(t *testing.T) {
	tmpHome := t.TempDir()
	tmpProject := t.TempDir()
	globalDir := filepath.Join(tmpHome, ".claude", "context")
	if err := os.MkdirAll(globalDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(globalDir, "repos.md"), []byte("Repo info"), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := loadContextFilesWithHome(tmpProject, tmpHome)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Repo info" {
		t.Errorf("expected 'Repo info', got %q", result)
	}
}

func TestLoadContextFiles_BothDirs(t *testing.T) {
	tmpHome := t.TempDir()
	tmpProject := t.TempDir()

	globalDir := filepath.Join(tmpHome, ".claude", "context")
	os.MkdirAll(globalDir, 0755)
	os.WriteFile(filepath.Join(globalDir, "global.md"), []byte("Global"), 0644)

	projectDir := filepath.Join(tmpProject, ".chief", "context")
	os.MkdirAll(projectDir, 0755)
	os.WriteFile(filepath.Join(projectDir, "project.md"), []byte("Project"), 0644)

	result, err := loadContextFilesWithHome(tmpProject, tmpHome)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Global") {
		t.Error("expected result to contain 'Global'")
	}
	if !strings.Contains(result, "Project") {
		t.Error("expected result to contain 'Project'")
	}
	// Global should come before Project
	globalIdx := strings.Index(result, "Global")
	projectIdx := strings.Index(result, "Project")
	if globalIdx > projectIdx {
		t.Error("expected global context before project context")
	}
}

func TestLoadContextFiles_EmptyHomeDir(t *testing.T) {
	tmpProject := t.TempDir()
	contextDir := filepath.Join(tmpProject, ".chief", "context")
	os.MkdirAll(contextDir, 0755)
	os.WriteFile(filepath.Join(contextDir, "info.md"), []byte("Info"), 0644)

	result, err := loadContextFilesWithHome(tmpProject, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Info" {
		t.Errorf("expected 'Info', got %q", result)
	}
}

func TestLoadMarkdownFiles_SortedOrder(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "02-schema.md"), []byte("Schema"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "01-platform.md"), []byte("Platform"), 0644)

	result, err := loadMarkdownFiles(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	expected := "Platform\n\n---\n\nSchema"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestLoadMarkdownFiles_IgnoresNonMdFiles(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "notes.txt"), []byte("TXT"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "context.md"), []byte("MD"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "data.json"), []byte("{}"), 0644)

	result, err := loadMarkdownFiles(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	if result != "MD" {
		t.Errorf("expected 'MD', got %q", result)
	}
}

func TestLoadMarkdownFiles_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()
	result, err := loadMarkdownFiles(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestLoadMarkdownFiles_NonExistentDir(t *testing.T) {
	result, err := loadMarkdownFiles("/nonexistent/path")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestLoadMarkdownFiles_IgnoresSubdirectories(t *testing.T) {
	tmpDir := t.TempDir()
	os.MkdirAll(filepath.Join(tmpDir, "subdir"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "subdir", "nested.md"), []byte("Nested"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "top.md"), []byte("Top"), 0644)

	result, err := loadMarkdownFiles(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	if result != "Top" {
		t.Errorf("expected 'Top', got %q", result)
	}
}
