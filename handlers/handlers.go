package handlers

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/imhasandl/file-manager/models"
	"github.com/imhasandl/file-manager/utils"
)

type Config struct {
	tasks    map[uuid.UUID]*models.Task
	mu       sync.Mutex
	maxTasks int
	maxFiles int
}

func NewAPIConfig() *Config {

	return &Config{
		tasks:    make(map[uuid.UUID]*models.Task),
		maxTasks: 3,
		maxFiles: 3,
	}
}

func (cfg *Config) CreateTask(w http.ResponseWriter, r *http.Request) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	activeTasks := 0
	for _, task := range cfg.tasks {
		if task.ArchiveURL == "" {
			activeTasks++
		}
	}

	if activeTasks >= cfg.maxTasks {
		utils.RespondWithError(w, http.StatusTooManyRequests, "Сервер занят. Задачи превышают лимит: 3", nil)
		return
	}

	task := &models.Task{
		ID:        uuid.New(),
		Files:     make([]models.FileInfo, 0),
		CreatedAt: time.Now(),
	}

	cfg.tasks[task.ID] = task
	utils.RespondWithJSON(w, http.StatusCreated, task)
}

func (cfg *Config) AddFile(w http.ResponseWriter, r *http.Request) {
	taskID, err := utils.GetTaskIDFromPath(r.URL.Path)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var requestData struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "неверный JSON", err)
		return
	}

	if err := utils.ValidateURL(requestData.URL); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	task, exists := cfg.tasks[taskID]
	if !exists {
		utils.RespondWithError(w, http.StatusNotFound, "задача не найдена", nil)
		return
	}

	if task.ArchiveURL != "" {
		utils.RespondWithError(w, http.StatusBadRequest, "задача уже в архиве", nil)
		return
	}

	if len(task.Files) >= cfg.maxFiles {
		utils.RespondWithError(w, http.StatusBadRequest, "достигнуто максимальное количество файлов", nil)
		return
	}

	if !utils.IsValidFileType(requestData.URL) {
		utils.RespondWithError(w, http.StatusBadRequest, "файл не является .pdf или .jpeg типом", nil)
		return
	}

	filename := utils.GetFilenameFromURL(requestData.URL)
	fileInfo := models.FileInfo{
		URL:      requestData.URL,
		Filename: filename,
	}

	task.Files = append(task.Files, fileInfo)

	if len(task.Files) == cfg.maxFiles {
		go cfg.processTask(task)
	}

	utils.RespondWithJSON(w, http.StatusOK, task)
}

func (cfg *Config) processTask(task *models.Task) {
	archivePath := utils.CreateArchivePath(task.ID)

	zipFile, err := os.Create(archivePath)
	if err != nil {
		return
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, fileInfo := range task.Files {
		err := utils.DownloadFileToZip(zipWriter, fileInfo.URL, fileInfo.Filename)
		if err != nil {
			continue
		}
	}

	cfg.mu.Lock()
	task.ArchiveURL = utils.CreateArchiveURL(task.ID)
	cfg.mu.Unlock()
}

func (cfg *Config) GetTaskStatus(w http.ResponseWriter, r *http.Request) {
	taskID, err := utils.GetTaskIDFromPath(r.URL.Path)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	task, exists := cfg.tasks[taskID]
	if !exists {
		utils.RespondWithError(w, http.StatusNotFound, "Задача не найдена", nil)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, task)
}

func (cfg *Config) DownloadArchive(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Base(r.URL.Path)
	filePath := filepath.Join("archives", filename)

	file, err := os.Open(filePath)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Файл не найден", err)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	io.Copy(w, file)
}
