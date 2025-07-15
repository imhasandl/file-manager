package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
)

func DownloadFileToZip(zipWriter *zip.Writer, fileURL, filename string) error {
	resp, err := http.Get(fileURL)
	if err != nil {
		return fmt.Errorf("не удалось скачать файл: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("файл недоступен, код ответа: %d", resp.StatusCode)
	}

	fileWriter, err := zipWriter.Create(filename)
	if err != nil {
		return fmt.Errorf("не удалось создать файл в архиве: %v", err)
	}

	_, err = io.Copy(fileWriter, resp.Body)
	if err != nil {
		return fmt.Errorf("не удалось записать файл в архив: %v", err)
	}

	return nil
}
