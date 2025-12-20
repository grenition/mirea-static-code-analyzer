package services

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"webapp/internal/app/models"
)

type AnalyzerService struct {
	rules         []Rule
	dockerService *DockerService
	parser        *LinterParser
	storageDir    string
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
	storageDir := os.Getenv("STORAGE_DIR")
	if storageDir == "" {
		storageDir = "./storage"
	}
	service := &AnalyzerService{
		dockerService: NewDockerService(),
		parser:        NewLinterParser(),
		storageDir:    storageDir,
	}
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

func (s *AnalyzerService) Analyze(files []FileInput, analysisID string) ([]Issue, map[string]int) {
	var allIssues []Issue

	// Group files by language for Docker analysis
	filesByLang := s.groupFilesByLanguage(files)

	// Run Docker-based linters
	dockerIssues := s.runDockerLinters(filesByLang, analysisID)
	allIssues = append(allIssues, dockerIssues...)

	// Run simple rules as fallback/complement
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

func (s *AnalyzerService) groupFilesByLanguage(files []FileInput) map[string][]FileInput {
	groups := make(map[string][]FileInput)
	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file.Path))
		lang := s.getLanguageFromExt(ext)
		if lang != "" {
			groups[lang] = append(groups[lang], file)
		}
	}
	return groups
}

func (s *AnalyzerService) getLanguageFromExt(ext string) string {
	langMap := map[string]string{
		".js":   "javascript",
		".jsx":  "javascript",
		".ts":   "typescript",
		".tsx":  "typescript",
		".py":   "python",
		".go":   "go",
		".java": "java",
		".cpp":  "cpp",
		".c":    "cpp",
		".h":    "cpp",
		".hpp":  "cpp",
		".php":  "php",
		".rb":   "ruby",
		".kt":   "kotlin",
		".swift": "swift",
	}
	return langMap[ext]
}

func (s *AnalyzerService) runDockerLinters(filesByLang map[string][]FileInput, analysisID string) []Issue {
	var allIssues []Issue

	// Check if Docker container is available
	if err := s.dockerService.EnsureContainerRunning(); err != nil {
		// Container not running, skip Docker linters
		return allIssues
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for lang, files := range filesByLang {
		// Create temporary directory for this analysis
		tmpDir, err := s.createTempAnalysisDir(files, analysisID)
		if err != nil {
			continue
		}
		defer os.RemoveAll(tmpDir)

		// Get Docker workdir path (storage is mounted at /workspace)
		// tmpDir is relative to storageDir, so we need the relative path
		relPath, err := filepath.Rel(s.storageDir, tmpDir)
		if err != nil {
			relPath = filepath.Base(tmpDir)
		}
		dockerWorkDir := filepath.Join("/workspace", relPath)

		// Run appropriate linter based on language
		var issues []Issue
		switch lang {
		case "javascript", "typescript":
			issues = s.runESLint(ctx, files, dockerWorkDir)
		case "python":
			issues = s.runPylint(ctx, files, dockerWorkDir)
		case "go":
			issues = s.runGolangciLint(ctx, files, dockerWorkDir)
		case "cpp":
			issues = s.runCppcheck(ctx, files, dockerWorkDir)
		case "ruby":
			issues = s.runRubocop(ctx, files, dockerWorkDir)
		case "php":
			issues = s.runPHPCS(ctx, files, dockerWorkDir)
		}

		allIssues = append(allIssues, issues...)
	}

	return allIssues
}

func (s *AnalyzerService) createTempAnalysisDir(files []FileInput, analysisID string) (string, error) {
	// Use analysisID to create a consistent directory
	analysisDir := filepath.Join(s.storageDir, analysisID, "linter_analysis")
	if err := os.MkdirAll(analysisDir, 0755); err != nil {
		return "", err
	}
	
	// Create a unique subdirectory for this batch
	tmpDir, err := os.MkdirTemp(analysisDir, "batch-*")
	if err != nil {
		return "", err
	}

	for _, file := range files {
		filePath := filepath.Join(tmpDir, file.Path)
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return tmpDir, err
		}
		if err := os.WriteFile(filePath, file.Content, 0644); err != nil {
			return tmpDir, err
		}
	}

	return tmpDir, nil
}

func (s *AnalyzerService) runESLint(ctx context.Context, files []FileInput, workDir string) []Issue {
	var allIssues []Issue
	for _, file := range files {
		// Run ESLint on single file (use --no-eslintrc to avoid config issues)
		cmd := []string{"npx", "--yes", "eslint", "--no-eslintrc", "--format", "compact", file.Path}
		stdout, stderr, _ := s.dockerService.RunLinter(ctx, cmd, workDir)
		// ESLint may return non-zero exit code with issues, check both stdout and stderr
		output := stdout + stderr
		if output != "" {
			issues := s.parser.ParseESLint(output, file.Path)
			allIssues = append(allIssues, issues...)
		}
	}
	return allIssues
}

func (s *AnalyzerService) runPylint(ctx context.Context, files []FileInput, workDir string) []Issue {
	var allIssues []Issue
	for _, file := range files {
		cmd := []string{"pylint", "--output-format=text", file.Path}
		stdout, _, _ := s.dockerService.RunLinter(ctx, cmd, workDir)
		issues := s.parser.ParsePylint(stdout, file.Path)
		allIssues = append(allIssues, issues...)
	}
	return allIssues
}

func (s *AnalyzerService) runGolangciLint(ctx context.Context, files []FileInput, workDir string) []Issue {
	var allIssues []Issue
	// golangci-lint works on packages/directories
	// Run on the workDir to analyze all Go files
	cmd := []string{"golangci-lint", "run", "--no-config", "--disable-all", "--enable=errcheck,govet,staticcheck", "."}
	stdout, _, _ := s.dockerService.RunLinter(ctx, cmd, workDir)
	// Parse output for all files
	for _, file := range files {
		issues := s.parser.ParseGolangciLint(stdout, file.Path)
		allIssues = append(allIssues, issues...)
	}
	return allIssues
}

func (s *AnalyzerService) runCppcheck(ctx context.Context, files []FileInput, workDir string) []Issue {
	var allIssues []Issue
	for _, file := range files {
		cmd := []string{"cppcheck", "--enable=all", "--output-file=-", file.Path}
		stdout, _, _ := s.dockerService.RunLinter(ctx, cmd, workDir)
		issues := s.parser.ParseCppcheck(stdout, file.Path)
		allIssues = append(allIssues, issues...)
	}
	return allIssues
}

func (s *AnalyzerService) runRubocop(ctx context.Context, files []FileInput, workDir string) []Issue {
	var allIssues []Issue
	for _, file := range files {
		cmd := []string{"rubocop", "--format", "simple", file.Path}
		stdout, _, _ := s.dockerService.RunLinter(ctx, cmd, workDir)
		issues := s.parser.ParseRubocop(stdout, file.Path)
		allIssues = append(allIssues, issues...)
	}
	return allIssues
}

func (s *AnalyzerService) runPHPCS(ctx context.Context, files []FileInput, workDir string) []Issue {
	var allIssues []Issue
	for _, file := range files {
		// Use php -l for basic syntax checking
		cmd := []string{"php", "-l", file.Path}
		stdout, stderr, _ := s.dockerService.RunLinter(ctx, cmd, workDir)
		output := stdout + stderr
		// PHP -l outputs errors to stderr
		if strings.Contains(output, "Parse error") || strings.Contains(output, "Fatal error") {
			// Extract line number from PHP error
			re := regexp.MustCompile(`on line (\d+)`)
			matches := re.FindStringSubmatch(output)
			var lineNum *int
			if len(matches) >= 2 {
				if num, err := strconv.Atoi(matches[1]); err == nil {
					lineNum = &num
				}
			}
			allIssues = append(allIssues, Issue{
				Severity: "error",
				RuleCode: "PHP_SYNTAX",
				Message:  strings.TrimSpace(output),
				FilePath: file.Path,
				Line:     lineNum,
			})
		}
	}
	return allIssues
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

