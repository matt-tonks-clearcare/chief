// Package embed provides embedded prompt templates used by Chief.
// All prompts are embedded at compile time using Go's embed directive.
package embed

import (
	_ "embed"
	"strings"
)

//go:embed prompt.txt
var promptTemplate string

//go:embed init_prompt.txt
var initPromptTemplate string

//go:embed edit_prompt.txt
var editPromptTemplate string

//go:embed convert_prompt.txt
var convertPromptTemplate string

//go:embed detect_setup_prompt.txt
var detectSetupPromptTemplate string

// GetPrompt returns the agent prompt with the PRD path and ticket prefix substituted.
// If ticketPrefix is empty, the placeholder is replaced with "[Story ID]" so the
// agent falls back to using the story ID in the commit message.
func GetPrompt(prdPath, ticketPrefix string) string {
	result := strings.ReplaceAll(promptTemplate, "{{PRD_PATH}}", prdPath)
	if ticketPrefix == "" {
		ticketPrefix = "[Story ID]"
	}
	return strings.ReplaceAll(result, "{{TICKET_PREFIX}}", ticketPrefix)
}

// GetInitPrompt returns the PRD generator prompt with the PRD directory and optional context substituted.
func GetInitPrompt(prdDir, context string) string {
	if context == "" {
		context = "No additional context provided. Ask the user what they want to build."
	}
	result := strings.ReplaceAll(initPromptTemplate, "{{PRD_DIR}}", prdDir)
	return strings.ReplaceAll(result, "{{CONTEXT}}", context)
}

// GetEditPrompt returns the PRD editor prompt with the PRD directory substituted.
func GetEditPrompt(prdDir string) string {
	return strings.ReplaceAll(editPromptTemplate, "{{PRD_DIR}}", prdDir)
}

// GetConvertPrompt returns the PRD converter prompt with the PRD content inlined.
func GetConvertPrompt(prdContent string) string {
	return strings.ReplaceAll(convertPromptTemplate, "{{PRD_CONTENT}}", prdContent)
}

// GetDetectSetupPrompt returns the prompt for detecting project setup commands.
func GetDetectSetupPrompt() string {
	return detectSetupPromptTemplate
}
