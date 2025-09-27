package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"

	//"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"test_task/pkg/tools"
	"time"
)

type Handler struct {
	pathDir string //путь для скачивания файлов
}

// конструктор
func NewHandler(pathDir string) (*Handler, error) {
	h := &Handler{
		pathDir: pathDir,
	}

	//проверяем существование директории для загрузки файлов, в случае отсутствия создаём
	err := tools.CreateDir(h.pathDir)
	if err != nil {
		log.Printf("Error open/create directory for download. Path: %s\n", h.pathDir)
		return nil, err
	}

	return h, nil
}

// получает ссылку для скачивания и сохранения
func (h *Handler) Download(downloadURL string) (string, error) {

	if !strings.HasPrefix(downloadURL, "http") {
		log.Printf("Error, invalid URL: %s\n", downloadURL)
		return "", fmt.Errorf("invalid URL: %s", downloadURL)
	}

	// простой GET запрос
	resp, err := http.Get(downloadURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	//проверка статуса кода
	if resp.StatusCode != http.StatusOK {
		log.Printf("Error status code: %d,\n url: %s\n", resp.StatusCode, downloadURL)
		return "", fmt.Errorf("Error status code: %d,\n url: %s\n", resp.StatusCode, downloadURL)
	}

	//созаём имя файла из URL
	filename := filepath.Base(downloadURL)
	if filename == "" || filename == "." || filename == "/" {
		filename = "download_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	}

	//filePath содержит путь скачиваемого файла с его названием
	filePath := filepath.Join(h.pathDir, filename)

	//создаём файл
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file for: %s\n", downloadURL)
		return "", err
	}
	defer file.Close()

	//копируем данные
	_, err = io.Copy(file, resp.Body)
	file.Close()
	if err != nil {
		log.Printf("Error saving file for %s\n", downloadURL)
		//удаляем частично скаченный файл
		os.Remove(filePath)
		return "", err
	}

	log.Printf("Successfuly downloaded: %s -> %s", downloadURL, filePath)
	return filePath, nil

}
