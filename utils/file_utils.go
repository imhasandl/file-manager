package utils

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func IsValidFileType(fileURL string) bool {
	lower := strings.ToLower(fileURL)

	if strings.HasSuffix(lower, ".pdf") || strings.HasSuffix(lower, ".jpeg") {
		return true
	}

	return false
}

func GetFilenameFromURL(fileURL string) string {
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return "unknown_file"
	}

	filename := filepath.Base(parsedURL.Path)
	if filename == "." || filename == "/" || filename == "" {
		return "unknown_file"
	}

	return filename
}

func CreateArchivePath(taskID uuid.UUID) string {
	return fmt.Sprintf("archives/archive_%s.zip", taskID.String())
}

func CreateArchiveURL(taskID uuid.UUID) string {
	return fmt.Sprintf("/download/archive_%s.zip", taskID.String())
}
