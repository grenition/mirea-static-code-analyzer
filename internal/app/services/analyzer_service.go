package services

import (
	"regexp"
	"strings"
	"github.com/google/uuid"
	"webapp/internal/app/models"
)

type AnalyzerService struct {
	rules []Rule
}

type Rule struct {
	Code     string
	Severity string
	Check    func(file FileInput) []Issue
}

type Issue struct {
	Severity string
	RuleCode string
	Message  string
	FilePath string
	Line     *int
}

func NewAnalyzerService() *AnalyzerService {
	service := &AnalyzerService{}
	service.initRules()
	return service
}

func (s *AnalyzerService) initRules() {
	s.rules = []Rule{
		{
			Code:     "LINE_TOO_LONG",
			Severity: "warn",
			Check: func(file FileInput) []Issue {
				var issues []Issue
				lines := strings.Split(string(file.Content), "\n")
				for i, line := range lines {
					if len(line) > 120 {
						lineNum := i + 1
						issues = append(issues, Issue{
							Severity: "warn",
							RuleCode: "LINE_TOO_LONG",
							Message:  "Line exceeds 120 characters",
							FilePath: file.Path,
							Line:     &lineNum,
						})
					}
				}
				return issues
			},
		},
		{
			Code:     "TRAILING_WHITESPACE",
			Severity: "info",
			Check: func(file FileInput) []Issue {
				var issues []Issue
				lines := strings.Split(string(file.Content), "\n")
				for i, line := range lines {
					if len(line) > 0 && (line[len(line)-1] == ' ' || line[len(line)-1] == '\t') {
						lineNum := i + 1
						issues = append(issues, Issue{
							Severity: "info",
							RuleCode: "TRAILING_WHITESPACE",
							Message:  "Trailing whitespace detected",
							FilePath: file.Path,
							Line:     &lineNum,
						})
					}
				}
				return issues
			},
		},
		{
			Code:     "TAB_INDENT",
			Severity: "info",
			Check: func(file FileInput) []Issue {
				var issues []Issue
				content := string(file.Content)
				if strings.Contains(content, "\t") {
					issues = append(issues, Issue{
						Severity: "info",
						RuleCode: "TAB_INDENT",
						Message:  "File contains tab characters",
						FilePath: file.Path,
						Line:     nil,
					})
				}
				return issues
			},
		},
		{
			Code:     "TODO_FOUND",
			Severity: "warn",
			Check: func(file FileInput) []Issue {
				var issues []Issue
				lines := strings.Split(string(file.Content), "\n")
				for i, line := range lines {
					upperLine := strings.ToUpper(line)
					if strings.Contains(upperLine, "TODO") || strings.Contains(upperLine, "FIXME") {
						lineNum := i + 1
						issues = append(issues, Issue{
							Severity: "warn",
							RuleCode: "TODO_FOUND",
							Message:  "TODO or FIXME comment found",
							FilePath: file.Path,
							Line:     &lineNum,
						})
					}
				}
				return issues
			},
		},
		{
			Code:     "DEBUG_PRINT",
			Severity: "info",
			Check: func(file FileInput) []Issue {
				var issues []Issue
				content := string(file.Content)
				lines := strings.Split(content, "\n")
				for i, line := range lines {
					upperLine := strings.ToUpper(line)
					if strings.Contains(upperLine, "CONSOLE.LOG(") ||
						strings.Contains(upperLine, "FMT.PRINTLN(") ||
						strings.Contains(upperLine, "PRINT(") {
						lineNum := i + 1
						issues = append(issues, Issue{
							Severity: "info",
							RuleCode: "DEBUG_PRINT",
							Message:  "Debug print statement found",
							FilePath: file.Path,
							Line:     &lineNum,
						})
					}
				}
				return issues
			},
		},
		{
			Code:     "HARDCODED_SECRET",
			Severity: "error",
			Check: func(file FileInput) []Issue {
				var issues []Issue
				content := string(file.Content)
				lines := strings.Split(content, "\n")
				pattern := regexp.MustCompile(`(?i)(password|apikey|secret|token)\s*=\s*["'][^"']+["']`)
				for i, line := range lines {
					if pattern.MatchString(line) {
						lineNum := i + 1
						issues = append(issues, Issue{
							Severity: "error",
							RuleCode: "HARDCODED_SECRET",
							Message:  "Potential hardcoded secret detected",
							FilePath: file.Path,
							Line:     &lineNum,
						})
					}
				}
				return issues
			},
		},
		{
			Code:     "INSECURE_HTTP_URL",
			Severity: "warn",
			Check: func(file FileInput) []Issue {
				var issues []Issue
				content := string(file.Content)
				if strings.Contains(content, "http://") {
					lines := strings.Split(content, "\n")
					for i, line := range lines {
						if strings.Contains(line, "http://") {
							lineNum := i + 1
							issues = append(issues, Issue{
								Severity: "warn",
								RuleCode: "INSECURE_HTTP_URL",
								Message:  "Insecure HTTP URL detected",
								FilePath: file.Path,
								Line:     &lineNum,
							})
						}
					}
				}
				return issues
			},
		},
		{
			Code:     "FILE_TOO_LARGE",
			Severity: "warn",
			Check: func(file FileInput) []Issue {
				var issues []Issue
				lines := strings.Split(string(file.Content), "\n")
				if len(lines) > 800 {
					issues = append(issues, Issue{
						Severity: "warn",
						RuleCode: "FILE_TOO_LARGE",
						Message:  "File exceeds 800 lines",
						FilePath: file.Path,
						Line:     nil,
					})
				}
				return issues
			},
		},
	}
}

func (s *AnalyzerService) Analyze(files []FileInput) ([]Issue, map[string]int) {
	var allIssues []Issue

	for _, file := range files {
		for _, rule := range s.rules {
			issues := rule.Check(file)
			allIssues = append(allIssues, issues...)
		}
	}

	summary := map[string]int{
		"error_count": 0,
		"warn_count":  0,
		"info_count":  0,
	}

	for _, issue := range allIssues {
		switch issue.Severity {
		case "error":
			summary["error_count"]++
		case "warn":
			summary["warn_count"]++
		case "info":
			summary["info_count"]++
		}
	}

	return allIssues, summary
}

func (s *AnalyzerService) ConvertToModels(issues []Issue, analysisID uuid.UUID) []*models.Issue {
	issueModels := make([]*models.Issue, len(issues))
	for i, issue := range issues {
		issueModels[i] = &models.Issue{
			ID:         uuid.New(),
			AnalysisID: analysisID,
			Severity:   issue.Severity,
			RuleCode:   issue.RuleCode,
			Message:    issue.Message,
			FilePath:   issue.FilePath,
			Line:       issue.Line,
		}
	}
	return issueModels
}

