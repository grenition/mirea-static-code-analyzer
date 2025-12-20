package services

import (
	"regexp"
	"strconv"
	"strings"
)

// LinterParser parses output from various linters
type LinterParser struct{}

func NewLinterParser() *LinterParser {
	return &LinterParser{}
}

// ParseESLint parses ESLint output (JSON format)
func (p *LinterParser) ParseESLint(output string, filePath string) []Issue {
	var issues []Issue
	// ESLint JSON format: [{"filePath":"file.js","messages":[...]}]
	// For simplicity, we'll parse text output if available, or use regex for JSON
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, filePath) {
			// Try to extract line number and message
			re := regexp.MustCompile(`(\d+):(\d+)\s+(error|warning|info)\s+(.+?)\s+\((.+?)\)`)
			matches := re.FindStringSubmatch(line)
			if len(matches) >= 6 {
				lineNum, _ := strconv.Atoi(matches[1])
				severity := strings.ToLower(matches[3])
				message := matches[4]
				ruleCode := matches[5]

				// Map ESLint severity to our severity
				sev := "info"
				if severity == "error" {
					sev = "error"
				} else if severity == "warning" {
					sev = "warn"
				}

				issues = append(issues, Issue{
					Severity: sev,
					RuleCode: "ESLINT_" + ruleCode,
					Message:  message,
					FilePath: filePath,
					Line:     &lineNum,
				})
			}
		}
	}
	return issues
}

// ParsePylint parses pylint output
func (p *LinterParser) ParsePylint(output string, filePath string) []Issue {
	var issues []Issue
	lines := strings.Split(output, "\n")
	re := regexp.MustCompile(`^([^:]+):(\d+):(\d+):\s*([^:]+):\s*(.+?)\s*\((.+?)\)`)

	for _, line := range lines {
		if !strings.Contains(line, filePath) {
			continue
		}
		matches := re.FindStringSubmatch(line)
		if len(matches) >= 7 {
			lineNum, _ := strconv.Atoi(matches[2])
			severity := strings.ToLower(matches[4])
			message := matches[5]
			ruleCode := matches[6]

			// Map pylint severity
			sev := "info"
			if severity == "error" || severity == "fatal" {
				sev = "error"
			} else if severity == "warning" {
				sev = "warn"
			}

			issues = append(issues, Issue{
				Severity: sev,
				RuleCode: "PYLINT_" + ruleCode,
				Message:  message,
				FilePath: filePath,
				Line:     &lineNum,
			})
		}
	}
	return issues
}

// ParseGolangciLint parses golangci-lint output
func (p *LinterParser) ParseGolangciLint(output string, filePath string) []Issue {
	var issues []Issue
	lines := strings.Split(output, "\n")
	re := regexp.MustCompile(`^([^:]+):(\d+):(\d+):\s*(.+?)\s*\((.+?)\)`)

	for _, line := range lines {
		if !strings.Contains(line, filePath) {
			continue
		}
		matches := re.FindStringSubmatch(line)
		if len(matches) >= 6 {
			lineNum, _ := strconv.Atoi(matches[2])
			message := matches[4]
			ruleCode := matches[5]

			// Determine severity from rule code (golangci-lint doesn't always include severity)
			sev := "warn"
			if strings.Contains(ruleCode, "errcheck") || strings.Contains(ruleCode, "gosec") {
				sev = "error"
			}

			issues = append(issues, Issue{
				Severity: sev,
				RuleCode: "GOLANGCI_" + ruleCode,
				Message:  message,
				FilePath: filePath,
				Line:     &lineNum,
			})
		}
	}
	return issues
}

// ParseCppcheck parses cppcheck output
func (p *LinterParser) ParseCppcheck(output string, filePath string) []Issue {
	var issues []Issue
	lines := strings.Split(output, "\n")
	re := regexp.MustCompile(`\[([^\]]+):(\d+)\]:\s*\(([^)]+)\)\s*(.+)`)

	for _, line := range lines {
		if !strings.Contains(line, filePath) {
			continue
		}
		matches := re.FindStringSubmatch(line)
		if len(matches) >= 5 {
			lineNum, _ := strconv.Atoi(matches[2])
			severity := strings.ToLower(matches[3])
			message := matches[4]

			sev := "info"
			if severity == "error" {
				sev = "error"
			} else if severity == "warning" {
				sev = "warn"
			}

			issues = append(issues, Issue{
				Severity: sev,
				RuleCode: "CPPCHECK_" + matches[1],
				Message:  message,
				FilePath: filePath,
				Line:     &lineNum,
			})
		}
	}
	return issues
}

// ParseRubocop parses rubocop output
func (p *LinterParser) ParseRubocop(output string, filePath string) []Issue {
	var issues []Issue
	lines := strings.Split(output, "\n")
	re := regexp.MustCompile(`^([^:]+):(\d+):(\d+):\s*([^:]+):\s*(.+?)\s*\((.+?)\)`)

	for _, line := range lines {
		if !strings.Contains(line, filePath) {
			continue
		}
		matches := re.FindStringSubmatch(line)
		if len(matches) >= 7 {
			lineNum, _ := strconv.Atoi(matches[2])
			severity := strings.ToLower(matches[4])
			message := matches[5]
			ruleCode := matches[6]

			sev := "info"
			if severity == "error" {
				sev = "error"
			} else if severity == "warning" {
				sev = "warn"
			}

			issues = append(issues, Issue{
				Severity: sev,
				RuleCode: "RUBOCOP_" + ruleCode,
				Message:  message,
				FilePath: filePath,
				Line:     &lineNum,
			})
		}
	}
	return issues
}

// ParsePHPCS parses PHP CodeSniffer output
func (p *LinterParser) ParsePHPCS(output string, filePath string) []Issue {
	var issues []Issue
	lines := strings.Split(output, "\n")
	re := regexp.MustCompile(`(\d+)\s+\|\s+([^|]+)\s+\|\s+(.+)`)

	for _, line := range lines {
		if !strings.Contains(line, filePath) {
			continue
		}
		matches := re.FindStringSubmatch(line)
		if len(matches) >= 4 {
			lineNum, _ := strconv.Atoi(matches[1])
			severity := strings.ToLower(strings.TrimSpace(matches[2]))
			message := strings.TrimSpace(matches[3])

			sev := "info"
			if severity == "error" {
				sev = "error"
			} else if severity == "warning" {
				sev = "warn"
			}

			issues = append(issues, Issue{
				Severity: sev,
				RuleCode: "PHPCS",
				Message:  message,
				FilePath: filePath,
				Line:     &lineNum,
			})
		}
	}
	return issues
}

