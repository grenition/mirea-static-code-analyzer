package services

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FileInput struct {
	Path    string
	Content []byte
}

type StorageService struct {
	storageDir string
	allowedExts map[string]bool
}

func NewStorageService(storageDir string) *StorageService {
	allowedExts := map[string]bool{
		".go": true, ".cs": true, ".js": true, ".ts": true,
		".py": true, ".cpp": true, ".h": true,
		".hpp": true, ".php": true, ".rb": true, ".kt": true,
		".swift": true, ".sql": true, ".html": true, ".css": true,
		".md": true,
	}

	return &StorageService{
		storageDir:  storageDir,
		allowedExts: allowedExts,
	}
}

func (s *StorageService) ExtractZip(zipPath, destDir string) ([]FileInput, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var files []FileInput

	for _, f := range r.File {
		// Prevent path traversal
		cleanPath := filepath.Clean(f.Name)
		if strings.Contains(cleanPath, "..") || filepath.IsAbs(cleanPath) {
			continue
		}

		// Check extension
		ext := strings.ToLower(filepath.Ext(cleanPath))
		if !s.allowedExts[ext] {
			continue
		}

		// Skip directories and symlinks
		if f.FileInfo().IsDir() || f.FileInfo().Mode()&os.ModeSymlink != 0 {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			continue
		}

		content, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			continue
		}

		files = append(files, FileInput{
			Path:    cleanPath,
			Content: content,
		})
	}

	return files, nil
}

func (s *StorageService) SaveZipFile(uploadPath, analysisID string) (string, error) {
	destDir := filepath.Join(s.storageDir, analysisID)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", err
	}

	destPath := filepath.Join(destDir, "upload.zip")
	destFile, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer destFile.Close()

	srcFile, err := os.Open(uploadPath)
	if err != nil {
		return "", err
	}
	defer srcFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return destPath, err
}

