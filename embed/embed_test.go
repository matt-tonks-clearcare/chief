package embed

import (
	"strings"
	"testing"
)

func TestGetPrompt(t *testing.T) {
	prdPath := "/path/to/prd.json"
	prompt := GetPrompt(prdPath)

	// Verify the PRD path placeholder was substituted
	if strings.Contains(prompt, "{{PRD_PATH}}") {
		t.Error("Expected {{PRD_PATH}} to be substituted")
	}

	// Verify the PRD path appears in the prompt
	if !strings.Contains(prompt, prdPath) {
		t.Errorf("Expected prompt to contain PRD path %q", prdPath)
	}

	// Verify the prompt contains key instructions
	if !strings.Contains(prompt, "chief-complete") {
		t.Error("Expected prompt to contain chief-complete instruction")
	}

	if !strings.Contains(prompt, "ralph-status") {
		t.Error("Expected prompt to contain ralph-status instruction")
	}

	if !strings.Contains(prompt, "passes: true") {
		t.Error("Expected prompt to contain passes: true instruction")
	}
}

func TestPromptTemplateNotEmpty(t *testing.T) {
	if promptTemplate == "" {
		t.Error("Expected promptTemplate to be embedded and non-empty")
	}
}
