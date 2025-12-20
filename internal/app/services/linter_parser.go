package services

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// LinterParser parses output from various linters
type LinterParser struct{}

func NewLinterParser() *LinterParser {
	return &LinterParser{}
}

// ParseESLint parses ESLint output (compact format)
func (p *LinterParser) ParseESLint(output string, filePath string) []Issue {
	var issues []Issue
	lines := strings.Split(output, "\n")
	
	// ESLint compact format has two patterns:
	// 1. filePath: line X, col Y, Severity - message. (rule)
	//    Example: /tmp/test.js: line 4, col 10, Error - 'calculateSum' is defined but never used. (no-unused-vars)
	// 2. filePath:line:column: severity message (rule)
	//    Example: /path/to/file.js:10:5: error 'unused' is assigned a value but never used (no-unused-vars)
	
	re1 := regexp.MustCompile(`([^:]+):\s+line\s+(\d+),\s+col\s+(\d+),\s+(Error|Warning|Info)\s+-\s+(.+?)\s+\((.+?)\)`)
	re2 := regexp.MustCompile(`([^:]+):(\d+):(\d+):\s+(error|warning|info)\s+(.+?)\s+\((.+?)\)`)
	
	for _, line := range lines {
		if strings.Contains(line, filePath) || strings.Contains(line, "snippet") {
			// Try first pattern (with "line X, col Y")
			matches := re1.FindStringSubmatch(line)
			if len(matches) >= 7 {
				lineNum, _ := strconv.Atoi(matches[2])
				severity := strings.ToLower(matches[4])
				message := strings.TrimSpace(matches[5])
				ruleCode := matches[6]

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
				continue
			}
			
			// Try second pattern (standard format)
			matches = re2.FindStringSubmatch(line)
			if len(matches) >= 7 {
				lineNum, _ := strconv.Atoi(matches[2])
				severity := strings.ToLower(matches[4])
				message := strings.TrimSpace(matches[5])
				ruleCode := matches[6]

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
				continue
			}
			
			// Try simpler pattern for cases without rule code
			reSimple := regexp.MustCompile(`([^:]+):\s+line\s+(\d+),\s+col\s+(\d+),\s+(Error|Warning|Info)\s+-\s+(.+)`)
			matches = reSimple.FindStringSubmatch(line)
			if len(matches) >= 6 {
				lineNum, _ := strconv.Atoi(matches[2])
				severity := strings.ToLower(matches[4])
				message := strings.TrimSpace(matches[5])

				sev := "info"
				if severity == "error" {
					sev = "error"
				} else if severity == "warning" {
					sev = "warn"
				}

				issues = append(issues, Issue{
					Severity: sev,
					RuleCode: "ESLINT",
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
func (p *LinterParser) ParsePylint(output string, filePath string, dockerPath string) []Issue {
	var issues []Issue
	lines := strings.Split(output, "\n")
	// Pylint format: /path/to/file.py:line:column: severity: message (rule-code)
	// Example: /tmp/test.py:8:4: W0612: Unused variable 'unused' (unused-variable)
	re := regexp.MustCompile(`^([^:]+):(\d+):(\d+):\s*([^:]+):\s*(.+?)\s*\((.+?)\)`)

	// Extract base filename for matching
	baseFile := filepath.Base(filePath)
	baseDockerFile := filepath.Base(dockerPath)

	for _, line := range lines {
		// Skip lines that don't look like pylint output
		lineTrimmed := strings.TrimSpace(line)
		if lineTrimmed == "" || 
		   strings.HasPrefix(lineTrimmed, "***") ||
		   strings.HasPrefix(lineTrimmed, "---") ||
		   strings.HasPrefix(lineTrimmed, "Your code has been rated") ||
		   !strings.Contains(lineTrimmed, ":") {
			continue
		}
		
		matches := re.FindStringSubmatch(lineTrimmed)
		if len(matches) >= 7 {
			matchedPath := matches[1]
			matchedBase := filepath.Base(matchedPath)
			
			// Match by base filename (most reliable) or if path contains our file path
			matchesFile := matchedBase == baseFile || matchedBase == baseDockerFile ||
			              strings.Contains(matchedPath, baseFile) || 
			              strings.Contains(matchedPath, baseDockerFile) ||
			              strings.Contains(matchedPath, filePath) || 
			              strings.Contains(matchedPath, dockerPath) ||
			              strings.Contains(matchedPath, "snippet")
			
			if matchesFile {
				lineNum, _ := strconv.Atoi(matches[2])
				severity := strings.ToLower(strings.TrimSpace(matches[4]))
				message := strings.TrimSpace(matches[5])
				ruleCode := strings.TrimSpace(matches[6])

				// Map pylint severity
				sev := "info"
				if severity == "error" || severity == "fatal" || severity == "e" || severity == "f" {
					sev = "error"
				} else if severity == "warning" || severity == "w" {
					sev = "warn"
				}

				issues = append(issues, Issue{
					Severity: sev,
					RuleCode: "PYLINT_" + ruleCode,
					Message:  message,
					FilePath: filePath, // Use original file path, not Docker path
					Line:     &lineNum,
				})
			}
		}
	}
	return issues
}

// ParseGolangciLint parses golangci-lint output
func (p *LinterParser) ParseGolangciLint(output string, filePath string) []Issue {
	var issues []Issue
	lines := strings.Split(output, "\n")
	
	// golangci-lint format: file.go:line:column: message (linter)
	// Example: test.go:27:6: unused declared and not used (typecheck)
	// Also: test.go:6:2: "strings" imported and not used (typecheck)
	re := regexp.MustCompile(`([^:]+):(\d+):(\d+):\s*(.+)$`)

	// Extract base filename for matching
	baseFile := filepath.Base(filePath)

	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) >= 5 {
			matchedFile := matches[1]
			matchedBase := filepath.Base(matchedFile)
			
			// Check if this line matches our file by base name
			if matchedBase == baseFile || strings.Contains(matchedFile, baseFile) || 
			   strings.Contains(line, filePath) || strings.Contains(line, "snippet") {
				lineNum, _ := strconv.Atoi(matches[2])
				// Extract full message - everything after column number
				fullMessage := strings.TrimSpace(matches[4])
				
				// Extract rule code from end of message if present: "message (linter)"
				ruleCode := ""
				ruleCodeRe := regexp.MustCompile(`\s+\(([^)]+)\)$`)
				if ruleMatches := ruleCodeRe.FindStringSubmatch(fullMessage); len(ruleMatches) >= 2 {
					ruleCode = ruleMatches[1]
					// Remove rule code from message
					fullMessage = ruleCodeRe.ReplaceAllString(fullMessage, "")
				}
				message := strings.TrimSpace(fullMessage)

				// Determine severity from rule code
				sev := "warn"
				if strings.Contains(ruleCode, "errcheck") || strings.Contains(ruleCode, "gosec") || 
				   strings.Contains(ruleCode, "staticcheck") || strings.Contains(ruleCode, "govet") ||
				   strings.Contains(ruleCode, "typecheck") {
					sev = "error"
				}

				if ruleCode == "" {
					ruleCode = "unknown"
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
	}
	return issues
}

// ParseCppcheck parses cppcheck output
func (p *LinterParser) ParseCppcheck(output string, filePath string, dockerPath string) []Issue {
	var issues []Issue
	lines := strings.Split(output, "\n")
	
	// cppcheck has multiple output formats:
	// 1. [filePath:line]: (severity) message [ruleCode]
	//    Example: [/workspace/test/snippet.cpp:25]: (error) Array 'arr[10]' accessed at index 15, which is out of bounds. [arrayIndexOutOfBounds]
	// 2. filePath:line:column: severity: message [ruleCode]
	//    Example: test/snippet.cpp:1:33: error: Null pointer dereference: p [nullPointer]
	
	// Pattern 1: [filePath:line]: (severity) message [ruleCode]
	// Match the whole line and extract parts
	re1 := regexp.MustCompile(`\[([^\]]+):(\d+)\]:\s*\(([^)]+)\)\s*(.+)$`)
	// Pattern 2: filePath:line:column: severity: message [ruleCode]
	// Match the whole line and extract parts
	re2 := regexp.MustCompile(`([^:]+):(\d+):(\d+):\s+(error|warning|style|performance|portability|information):\s*(.+)$`)

	// Extract base filename for matching
	baseFile := filepath.Base(filePath)
	baseDockerFile := filepath.Base(dockerPath)

	for _, line := range lines {
		// Skip lines that don't look like cppcheck output or are notes
		lineTrimmed := strings.TrimSpace(line)
		if lineTrimmed == "" || 
		   strings.HasPrefix(lineTrimmed, "Checking") || 
		   strings.HasPrefix(lineTrimmed, "note:") || 
		   strings.HasPrefix(lineTrimmed, "^") ||
		   !strings.Contains(lineTrimmed, ":") {
			continue
		}
		
		// Try first format: [filePath:line]: (severity) message
		matches := re1.FindStringSubmatch(line)
		if len(matches) >= 5 {
			matchedPath := matches[1]
			matchedBase := filepath.Base(matchedPath)
			
			// Match by base filename (most reliable) or if path contains our file path
			matchesFile := matchedBase == baseFile || matchedBase == baseDockerFile ||
			              strings.Contains(matchedPath, baseFile) || 
			              strings.Contains(matchedPath, baseDockerFile) ||
			              strings.Contains(matchedPath, filePath) || 
			              strings.Contains(matchedPath, dockerPath) ||
			              strings.Contains(matchedPath, "snippet")
			
			if matchesFile {
				lineNum, _ := strconv.Atoi(matches[2])
				severity := strings.ToLower(matches[3])
				// Extract full message - everything after severity
				fullMessage := strings.TrimSpace(matches[4])
				
				// Extract rule code from end of message if present: "message [ruleCode]"
				ruleCode := ""
				ruleCodeRe := regexp.MustCompile(`\s+\[([^\]]+)\]$`)
				if ruleMatches := ruleCodeRe.FindStringSubmatch(fullMessage); len(ruleMatches) >= 2 {
					ruleCode = ruleMatches[1]
					// Remove rule code from message
					fullMessage = ruleCodeRe.ReplaceAllString(fullMessage, "")
				}
				message := strings.TrimSpace(fullMessage)

				sev := "info"
				if severity == "error" {
					sev = "error"
				} else if severity == "warning" {
					sev = "warn"
				}

				// Extract rule code from message if not provided
				if ruleCode == "" {
					if strings.Contains(message, "uninitialized") || strings.Contains(message, "Uninitialized") {
						ruleCode = "uninitvar"
					} else if strings.Contains(message, "out of bounds") {
						ruleCode = "arrayIndexOutOfBounds"
					} else if strings.Contains(message, "null pointer") || strings.Contains(message, "Null pointer") {
						ruleCode = "nullPointer"
					} else if strings.Contains(message, "Memory leak") || strings.Contains(message, "memory leak") {
						ruleCode = "memleak"
					} else if strings.Contains(message, "deallocated") || strings.Contains(message, "Dereferencing") {
						ruleCode = "deallocuse"
					} else if strings.Contains(message, "not initialized") {
						ruleCode = "uninitdata"
					} else {
						ruleCode = severity
					}
				}

				issues = append(issues, Issue{
					Severity: sev,
					RuleCode: "CPPCHECK_" + ruleCode,
					Message:  message,
					FilePath: filePath, // Use original file path, not Docker path
					Line:     &lineNum,
				})
				continue
			}
		}
		
		// Try second format: filePath:line:column: severity: message
		matches = re2.FindStringSubmatch(line)
		if len(matches) >= 6 {
			matchedPath := matches[1]
			matchedBase := filepath.Base(matchedPath)
			
			// Match by base filename
			matchesFile := matchedBase == baseFile || matchedBase == baseDockerFile ||
			              strings.Contains(matchedPath, baseFile) || 
			              strings.Contains(matchedPath, baseDockerFile) ||
			              strings.Contains(matchedPath, filePath) || 
			              strings.Contains(matchedPath, dockerPath) ||
			              strings.Contains(matchedPath, "snippet")
			
			if matchesFile {
				lineNum, _ := strconv.Atoi(matches[2])
				severity := strings.ToLower(matches[4])
				// Extract full message - everything after severity
				fullMessage := strings.TrimSpace(matches[5])
				
				// Extract rule code from end of message if present: "message [ruleCode]"
				ruleCode := ""
				ruleCodeRe := regexp.MustCompile(`\s+\[([^\]]+)\]$`)
				if ruleMatches := ruleCodeRe.FindStringSubmatch(fullMessage); len(ruleMatches) >= 2 {
					ruleCode = ruleMatches[1]
					// Remove rule code from message
					fullMessage = ruleCodeRe.ReplaceAllString(fullMessage, "")
				}
				message := strings.TrimSpace(fullMessage)

				sev := "info"
				if severity == "error" {
					sev = "error"
				} else if severity == "warning" {
					sev = "warn"
				}

				// Extract rule code from message if not provided
				if ruleCode == "" {
					if strings.Contains(message, "uninitialized") || strings.Contains(message, "Uninitialized") {
						ruleCode = "uninitvar"
					} else if strings.Contains(message, "out of bounds") {
						ruleCode = "arrayIndexOutOfBounds"
					} else if strings.Contains(message, "null pointer") || strings.Contains(message, "Null pointer") {
						ruleCode = "nullPointer"
					} else if strings.Contains(message, "Memory leak") || strings.Contains(message, "memory leak") {
						ruleCode = "memleak"
					} else if strings.Contains(message, "deallocated") || strings.Contains(message, "Dereferencing") {
						ruleCode = "deallocuse"
					} else if strings.Contains(message, "not initialized") {
						ruleCode = "uninitdata"
					} else {
						ruleCode = severity
					}
				}

				issues = append(issues, Issue{
					Severity: sev,
					RuleCode: "CPPCHECK_" + ruleCode,
					Message:  message,
					FilePath: filePath,
					Line:     &lineNum,
				})
			}
		}
	}
	return issues
}

// ParsePHPCS parses PHP syntax checker output
func (p *LinterParser) ParsePHPCS(output string, filePath string) []Issue {
	var issues []Issue
	lines := strings.Split(output, "\n")
	
	// PHP -l output format: Parse error: syntax error, unexpected ';' in file.php on line 5
	re := regexp.MustCompile(`(Parse error|Fatal error|Warning):\s*(.+?)\s+in\s+[^\s]+\s+on\s+line\s+(\d+)`)

	for _, line := range lines {
		if strings.Contains(line, "Parse error") || strings.Contains(line, "Fatal error") || strings.Contains(line, "Warning") {
			matches := re.FindStringSubmatch(line)
			if len(matches) >= 4 {
				errorType := matches[1]
				message := strings.TrimSpace(matches[2])
				lineNum, _ := strconv.Atoi(matches[3])

				sev := "error"
				if strings.Contains(errorType, "Warning") {
					sev = "warn"
				}

				issues = append(issues, Issue{
					Severity: sev,
					RuleCode: "PHP_SYNTAX",
					Message:  message,
					FilePath: filePath,
					Line:     &lineNum,
				})
			} else {
				// Try simpler pattern
				reSimple := regexp.MustCompile(`on\s+line\s+(\d+)`)
				matches := reSimple.FindStringSubmatch(line)
				if len(matches) >= 2 {
					lineNum, _ := strconv.Atoi(matches[1])
					sev := "error"
					if strings.Contains(line, "Warning") {
						sev = "warn"
					}
					issues = append(issues, Issue{
						Severity: sev,
						RuleCode: "PHP_SYNTAX",
						Message:  strings.TrimSpace(line),
						FilePath: filePath,
						Line:     &lineNum,
					})
				}
			}
		}
	}
	return issues
}


