package patterns

import (
	"fmt"
	"strings"

	"github.com/codemeapixel/cadence/internal/git"
	"github.com/codemeapixel/cadence/internal/metrics"
)

// CommitMessageStrategy detects AI-generated commit message patterns
type CommitMessageStrategy struct {
	enabled bool
}

func NewCommitMessageStrategy() *CommitMessageStrategy {
	return &CommitMessageStrategy{enabled: true}
}

func (s *CommitMessageStrategy) Name() string {
	return "commit_message_analysis"
}

func (s *CommitMessageStrategy) Detect(pair *git.CommitPair, repoStats *metrics.RepositoryStats) (bool, string) {
	if !s.enabled {
		return false, ""
	}

	msg := strings.ToLower(pair.Current.Message)

	// AI-generated commit message patterns
	aiPatterns := []string{
		"implement",
		"add functionality",
		"update code",
		"refactor code",
		"improve implementation",
		"enhance functionality",
		"optimize performance",
		"fix issues",
		"improve code quality",
		"add new features",
		"update implementation",
		"add support for",
	}

	// Generic/template commit messages
	genericPatterns := []string{
		"initial commit",
		"update readme",
		"update dependencies",
		"minor fixes",
		"code cleanup",
		"bug fixes",
		"improvements",
		"updates",
		"changes",
		"modifications",
	}

	aiScore := 0
	genericScore := 0

	for _, pattern := range aiPatterns {
		if strings.Contains(msg, pattern) {
			aiScore++
		}
	}

	for _, pattern := range genericPatterns {
		if strings.Contains(msg, pattern) {
			genericScore++
		}
	}

	// Flag if message is very generic or uses multiple AI patterns
	if aiScore >= 2 || genericScore >= 1 {
		return true, fmt.Sprintf(
			"Suspicious commit message patterns - generic/AI-like phrasing (AI patterns: %d, generic: %d)",
			aiScore, genericScore,
		)
	}

	// Check for overly descriptive but generic messages
	words := strings.Fields(msg)
	if len(words) > 8 && (strings.Contains(msg, "implement") || strings.Contains(msg, "functionality")) {
		return true, "Overly verbose yet generic commit message - typical of AI generation"
	}

	return false, ""
}

// NamingPatternStrategy detects suspicious naming conventions
type NamingPatternStrategy struct {
	enabled bool
}

func NewNamingPatternStrategy() *NamingPatternStrategy {
	return &NamingPatternStrategy{enabled: true}
}

func (s *NamingPatternStrategy) Name() string {
	return "naming_pattern_analysis"
}

func (s *NamingPatternStrategy) Detect(pair *git.CommitPair, repoStats *metrics.RepositoryStats) (bool, string) {
	if !s.enabled {
		return false, ""
	}

	// Analyze actual diff content if available
	if pair.DiffContent != "" {
		return s.analyzeCodeContent(pair.DiffContent)
	}

	// Fallback to commit message analysis
	msg := strings.ToLower(pair.Current.Message)
	genericNames := []string{
		"variable", "function", "method", "class", "object", "instance",
		"data", "result", "value", "item", "element", "component",
		"helper", "utility", "manager", "handler", "service",
	}

	nameCount := 0
	for _, name := range genericNames {
		if strings.Contains(msg, name) {
			nameCount++
		}
	}

	if nameCount >= 2 {
		return true, fmt.Sprintf(
			"Commit message contains multiple generic naming terms (%d) - may indicate AI-generated variable names",
			nameCount,
		)
	}

	return false, ""
}

func (s *NamingPatternStrategy) analyzeCodeContent(diffContent string) (bool, string) {
	lines := strings.Split(diffContent, "\n")
	addedLines := make([]string, 0)

	// Extract only added lines (starting with +)
	for _, line := range lines {
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			addedLines = append(addedLines, strings.TrimPrefix(line, "+"))
		}
	}

	if len(addedLines) == 0 {
		return false, ""
	}

	// AI slop patterns in actual code
	suspiciousPatterns := 0
	totalPatterns := 0

	codeContent := strings.Join(addedLines, "\n")

	// Generic variable names
	genericVarPatterns := []string{
		"var1", "var2", "temp", "data", "result", "value", "item",
		"element", "obj", "instance", "helper", "utility", "manager",
	}
	for _, pattern := range genericVarPatterns {
		if strings.Contains(strings.ToLower(codeContent), pattern) {
			suspiciousPatterns++
		}
		totalPatterns++
	}

	// TODO comments (AI often generates these)
	todoCount := strings.Count(strings.ToLower(codeContent), "todo")
	fixmeCount := strings.Count(strings.ToLower(codeContent), "fixme")
	if todoCount > 2 || fixmeCount > 1 {
		suspiciousPatterns++
	}

	// Overly consistent naming (camelCase perfection)
	words := strings.Fields(codeContent)
	perfectCamelCaseCount := 0
	for _, word := range words {
		if len(word) > 4 && isPerfectCamelCase(word) {
			perfectCamelCaseCount++
		}
	}
	if len(words) > 10 && float64(perfectCamelCaseCount)/float64(len(words)) > 0.3 {
		suspiciousPatterns++
	}

	// Check for suspiciously perfect error handling
	if strings.Contains(codeContent, "catch") || strings.Contains(codeContent, "except") {
		errorHandlingCount := strings.Count(codeContent, "catch") + strings.Count(codeContent, "except") + strings.Count(codeContent, "try")
		if errorHandlingCount > len(addedLines)/20 { // Too many try-catch blocks
			suspiciousPatterns++
		}
	}

	if suspiciousPatterns >= 2 {
		return true, fmt.Sprintf(
			"Code contains multiple AI-slop patterns (%d detected) - generic names, TODO comments, perfect patterns",
			suspiciousPatterns,
		)
	}

	return false, ""
}

func isPerfectCamelCase(word string) bool {
	if len(word) < 2 {
		return false
	}

	// Check if it starts with lowercase and has exactly one uppercase letter
	if !(word[0] >= 'a' && word[0] <= 'z') {
		return false
	}

	uppercaseCount := 0
	for _, r := range word {
		if r >= 'A' && r <= 'Z' {
			uppercaseCount++
		}
	}

	return uppercaseCount == 1
}

// StructuralConsistencyStrategy detects overly consistent code structure patterns
type StructuralConsistencyStrategy struct {
	enabled bool
}

func NewStructuralConsistencyStrategy() *StructuralConsistencyStrategy {
	return &StructuralConsistencyStrategy{enabled: true}
}

func (s *StructuralConsistencyStrategy) Name() string {
	return "structural_consistency_analysis"
}

func (s *StructuralConsistencyStrategy) Detect(pair *git.CommitPair, repoStats *metrics.RepositoryStats) (bool, string) {
	if !s.enabled {
		return false, ""
	}

	// Flag commits where additions and deletions are suspiciously proportional
	if pair.Stats.Additions > 100 && pair.Stats.Deletions > 100 {
		ratio := float64(pair.Stats.Additions) / float64(pair.Stats.Deletions)

		// AI often generates very balanced refactoring (close to 1:1 ratio)
		if ratio >= 0.9 && ratio <= 1.1 {
			return true, fmt.Sprintf(
				"Suspiciously balanced addition/deletion ratio: %.2f - may indicate automated refactoring",
				ratio,
			)
		}

		// Or very consistent ratios (e.g., exactly 2:1, 3:1)
		if isNearInteger(ratio, 0.05) || isNearInteger(1.0/ratio, 0.05) {
			return true, fmt.Sprintf(
				"Suspiciously consistent addition/deletion ratio: %.2f - may indicate template-based generation",
				ratio,
			)
		}
	}

	return false, ""
}

// BurstPatternStrategy detects suspicious commit timing bursts
type BurstPatternStrategy struct {
	maxCommitsPerHour int
	enabled           bool
}

func NewBurstPatternStrategy(maxPerHour int) *BurstPatternStrategy {
	return &BurstPatternStrategy{
		maxCommitsPerHour: maxPerHour,
		enabled:           true,
	}
}

func (s *BurstPatternStrategy) Name() string {
	return "burst_pattern_analysis"
}

func (s *BurstPatternStrategy) Detect(pair *git.CommitPair, repoStats *metrics.RepositoryStats) (bool, string) {
	if !s.enabled || s.maxCommitsPerHour <= 0 {
		return false, ""
	}

	// This would need to be implemented with access to broader commit history
	// For now, check if commits are extremely close together (< 5 minutes)
	if pair.TimeDelta.Seconds() < 300 && pair.Stats.Additions > 50 {
		return true, fmt.Sprintf(
			"Rapid commit pattern: %.1f seconds between substantial commits - may indicate batch processing",
			pair.TimeDelta.Seconds(),
		)
	}

	return false, ""
}

// ErrorHandlingPatternStrategy detects lack of error handling (AI often omits this)
type ErrorHandlingPatternStrategy struct {
	enabled bool
}

func NewErrorHandlingPatternStrategy() *ErrorHandlingPatternStrategy {
	return &ErrorHandlingPatternStrategy{enabled: true}
}

func (s *ErrorHandlingPatternStrategy) Name() string {
	return "error_handling_analysis"
}

func (s *ErrorHandlingPatternStrategy) Detect(pair *git.CommitPair, repoStats *metrics.RepositoryStats) (bool, string) {
	if !s.enabled {
		return false, ""
	}

	// Analyze actual code content if available
	if pair.DiffContent != "" && pair.Stats.Additions > 50 {
		return s.analyzeErrorHandling(pair.DiffContent, pair.Stats.Additions)
	}

	// Fallback to commit message analysis for large additions
	if pair.Stats.Additions > 200 {
		msg := strings.ToLower(pair.Current.Message)
		hasErrorPatterns := strings.Contains(msg, "error") ||
			strings.Contains(msg, "exception") ||
			strings.Contains(msg, "try") ||
			strings.Contains(msg, "catch") ||
			strings.Contains(msg, "handle")

		if !hasErrorPatterns && pair.Stats.Additions > 300 {
			return true, fmt.Sprintf(
				"Large code addition (%d lines) with no error handling mentions - AI often omits error handling",
				pair.Stats.Additions,
			)
		}
	}

	return false, ""
}

func (s *ErrorHandlingPatternStrategy) analyzeErrorHandling(diffContent string, additions int64) (bool, string) {
	lines := strings.Split(diffContent, "\n")
	addedLines := make([]string, 0)

	// Extract only added lines
	for _, line := range lines {
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			addedLines = append(addedLines, strings.TrimPrefix(line, "+"))
		}
	}

	if len(addedLines) < 20 { // Too small to analyze
		return false, ""
	}

	codeContent := strings.ToLower(strings.Join(addedLines, "\n"))

	// Count actual error handling patterns
	errorHandlingPatterns := 0

	// Basic error handling
	errorHandlingPatterns += strings.Count(codeContent, "try")
	errorHandlingPatterns += strings.Count(codeContent, "catch")
	errorHandlingPatterns += strings.Count(codeContent, "except")
	errorHandlingPatterns += strings.Count(codeContent, "throw")
	errorHandlingPatterns += strings.Count(codeContent, "throws")

	// Advanced error handling
	errorHandlingPatterns += strings.Count(codeContent, "error")
	errorHandlingPatterns += strings.Count(codeContent, "exception")
	errorHandlingPatterns += strings.Count(codeContent, "handle")

	// Language-specific patterns
	errorHandlingPatterns += strings.Count(codeContent, "if err != nil") // Go
	errorHandlingPatterns += strings.Count(codeContent, ".catch(")       // JavaScript promises
	errorHandlingPatterns += strings.Count(codeContent, "rescue")        // Ruby

	// Calculate expected error handling density
	expectedErrorHandling := len(addedLines) / 30 // Roughly 1 error check per 30 lines

	// Flag if there's insufficient error handling for the amount of code
	if additions > 100 && errorHandlingPatterns < expectedErrorHandling {
		return true, fmt.Sprintf(
			"Large code addition (%d lines) with insufficient error handling (%d patterns, expected ~%d) - typical AI omission",
			additions, errorHandlingPatterns, expectedErrorHandling,
		)
	}

	// Also flag overly perfect error handling (every line has error checking)
	if errorHandlingPatterns > len(addedLines)/5 {
		return true, fmt.Sprintf(
			"Excessive error handling patterns (%d in %d lines) - may indicate AI over-compensation",
			errorHandlingPatterns, len(addedLines),
		)
	}

	return false, ""
}

// TemplatePatternStrategy detects template-like code patterns
type TemplatePatternStrategy struct {
	enabled bool
}

func NewTemplatePatternStrategy() *TemplatePatternStrategy {
	return &TemplatePatternStrategy{enabled: true}
}

func (s *TemplatePatternStrategy) Name() string {
	return "template_pattern_analysis"
}

func (s *TemplatePatternStrategy) Detect(pair *git.CommitPair, repoStats *metrics.RepositoryStats) (bool, string) {
	if !s.enabled {
		return false, ""
	}

	// Analyze actual code content if available
	if pair.DiffContent != "" && pair.Stats.Additions > 50 {
		detected, reason := s.analyzeTemplatePatterns(pair.DiffContent)
		if detected {
			return detected, reason
		}
	}

	// Fallback to commit message analysis
	msg := strings.ToLower(pair.Current.Message)
	templatePatterns := []string{
		"boilerplate", "template", "skeleton", "scaffold", "stub", "placeholder", "todo", "fixme",
	}

	patternCount := 0
	for _, pattern := range templatePatterns {
		if strings.Contains(msg, pattern) {
			patternCount++
		}
	}

	if patternCount > 0 && pair.Stats.Additions > 100 {
		return true, fmt.Sprintf(
			"Template/boilerplate patterns detected in large commit (%d lines) - may be AI-generated scaffold",
			pair.Stats.Additions,
		)
	}

	return false, ""
}

func (s *TemplatePatternStrategy) analyzeTemplatePatterns(diffContent string) (bool, string) {
	lines := strings.Split(diffContent, "\n")
	addedLines := make([]string, 0)

	// Extract only added lines
	for _, line := range lines {
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			addedLines = append(addedLines, strings.TrimPrefix(line, "+"))
		}
	}

	if len(addedLines) < 10 {
		return false, ""
	}

	codeContent := strings.Join(addedLines, "\n")
	lowerContent := strings.ToLower(codeContent)

	suspiciousPatterns := 0

	// Template comments
	templateComments := []string{
		"todo", "fixme", "placeholder", "implement", "add code here",
		"your code here", "example", "sample", "template", "boilerplate",
	}
	for _, pattern := range templateComments {
		if strings.Contains(lowerContent, pattern) {
			suspiciousPatterns++
		}
	}

	// Repetitive patterns (AI loves patterns)
	repetitiveCount := 0
	for i := 0; i < len(addedLines)-1; i++ {
		line1 := strings.TrimSpace(addedLines[i])
		line2 := strings.TrimSpace(addedLines[i+1])

		if len(line1) > 10 && len(line2) > 10 {
			// Check for very similar lines (only 1-2 words different)
			words1 := strings.Fields(line1)
			words2 := strings.Fields(line2)

			if len(words1) == len(words2) && len(words1) > 3 {
				differences := 0
				for j := 0; j < len(words1); j++ {
					if words1[j] != words2[j] {
						differences++
					}
				}
				if differences <= 2 { // Very similar lines
					repetitiveCount++
				}
			}
		}
	}

	if repetitiveCount > len(addedLines)/8 { // More than 12.5% repetitive
		suspiciousPatterns++
	}

	// Perfect indentation consistency (humans are messier)
	indentationLevels := make(map[int]int)
	for _, line := range addedLines {
		if strings.TrimSpace(line) != "" {
			indent := 0
			for _, r := range line {
				if r == ' ' {
					indent++
				} else if r == '\t' {
					indent += 4
				} else {
					break
				}
			}
			indentationLevels[indent]++
		}
	}

	// If 90%+ of lines have consistent indentation increments
	totalLines := len(addedLines)
	if totalLines > 20 {
		consistentIndentation := 0
		for _, count := range indentationLevels {
			if count > totalLines/10 { // If any indentation level is used >10% of the time
				consistentIndentation += count
			}
		}
		if float64(consistentIndentation)/float64(totalLines) > 0.9 {
			suspiciousPatterns++
		}
	}

	// Unused imports/variables (AI often generates these)
	if strings.Contains(lowerContent, "import") {
		importLines := make([]string, 0)
		for _, line := range addedLines {
			if strings.Contains(strings.ToLower(line), "import") {
				importLines = append(importLines, line)
			}
		}

		// Rough heuristic: if there are many imports but little code usage
		if len(importLines) > 5 && len(addedLines) < len(importLines)*10 {
			suspiciousPatterns++
		}
	}

	if suspiciousPatterns >= 2 {
		return true, fmt.Sprintf(
			"Template/generated code patterns detected (%d indicators) - repetitive structure, perfect formatting, template comments",
			suspiciousPatterns,
		)
	}

	return false, ""
}

// Helper functions

// isNearInteger checks if a float64 is close to an integer value
func isNearInteger(value, tolerance float64) bool {
	remainder := value - float64(int(value))
	return remainder < tolerance || remainder > (1.0-tolerance)
}

// FileExtensionPattern detects suspicious file creation patterns
type FileExtensionPatternStrategy struct {
	enabled bool
}

func NewFileExtensionPatternStrategy() *FileExtensionPatternStrategy {
	return &FileExtensionPatternStrategy{enabled: true}
}

func (s *FileExtensionPatternStrategy) Name() string {
	return "file_extension_analysis"
}

func (s *FileExtensionPatternStrategy) Detect(pair *git.CommitPair, repoStats *metrics.RepositoryStats) (bool, string) {
	if !s.enabled {
		return false, ""
	}

	// Check for suspicious creation of many files at once (typical of generators)
	if pair.Stats.FilesChanged > 10 && pair.Stats.Additions > 1000 {
		avgLinesPerFile := float64(pair.Stats.Additions) / float64(pair.Stats.FilesChanged)

		// AI generators often create files with very consistent sizes
		if avgLinesPerFile > 50 && avgLinesPerFile < 200 {
			consistency := 1.0 - (avgLinesPerFile - float64(int(avgLinesPerFile)))
			if consistency > 0.8 {
				return true, fmt.Sprintf(
					"Suspicious file creation pattern: %d files with consistent size (~%.0f lines each) - may be generated",
					pair.Stats.FilesChanged, avgLinesPerFile,
				)
			}
		}
	}

	return false, ""
}
