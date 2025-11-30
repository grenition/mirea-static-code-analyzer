package models

type AnalyzeRequest struct {
	Files []FileInput `json:"files"`
}

type FileInput struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type AnalyzeResponse struct {
	Files []FileResult `json:"files"`
}

type FileResult struct {
	Path         string       `json:"path"`
	Comment      string       `json:"comment"`
	LineComments []LineComment `json:"line_comments"`
}

type LineComment struct {
	Line    int    `json:"line"`
	Comment string `json:"comment"`
}

