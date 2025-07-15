package models

import (
	"time"

	"github.com/google/uuid"
)

type FileInfo struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
}

type Task struct {
	ID         uuid.UUID  `json:"id"`
	Files      []FileInfo `json:"files"`
	CreatedAt  time.Time  `json:"created_at"`
	ArchiveURL string     `json:"archive_url,omitempty"`
}
