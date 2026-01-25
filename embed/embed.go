package embed

import (
	_ "embed"
	"strings"
)

//go:embed prompt.txt
var promptTemplate string

// GetPrompt returns the agent prompt with the PRD path substituted.
func GetPrompt(prdPath string) string {
	return strings.ReplaceAll(promptTemplate, "{{PRD_PATH}}", prdPath)
}
