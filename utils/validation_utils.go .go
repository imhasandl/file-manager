package utils

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

func ValidateURL(url string) error {
	if strings.TrimSpace(url) == "" {
		return errors.New("URL не может быть пустым")
	}
	return nil
}

func GetTaskIDFromPath(path string) (uuid.UUID, error) {
	pathSplit := strings.Split(path, "/")
	if len(pathSplit) < 3 {
		return uuid.Nil, errors.New("неверный формат URL")
	}

	taskID, err := uuid.Parse(pathSplit[2])
	if err != nil {
		return uuid.Nil, errors.New("неверный ID задачи")
	}

	return taskID, nil
}
