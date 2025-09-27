package repository

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"test_task/domain/models"
	"test_task/pkg/parser"
	"test_task/pkg/tools"
)

type Repository struct {
	pathDir string
	mu      sync.Mutex
}

// конструктор
func NewRepository(pathDir string) (*Repository, error) {
	repo := &Repository{
		pathDir: pathDir,
	}

	//проверяем директорию в случае отсутствия создаём её
	err := tools.CreateDir(repo.pathDir)
	if err != nil {
		log.Printf("Error create directory, path: %s\n", repo.pathDir)
		return nil, err
	}

	return repo, nil
}

// сохраняет/создаёт задачу
func (r *Repository) Save(task *models.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	//парсим структуру данных в json
	jsonFile, err := parser.ParsTaskToJson(task)
	if err != nil {
		log.Printf("Error create json file, task id: %d\n", task.ID)
		return err
	}

	//создаём путь к файлу используя id задания
	filePath := filepath.Join(r.pathDir, strconv.Itoa(task.ID)+".json")
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error open file for rewriting task, id: %d\n", task.ID)
		return err
	}
	defer file.Close()

	//записываем в файл данные json
	_, err = file.Write(jsonFile)
	if err != nil {
		log.Printf("Error writing task in file, id: %d\n", task.ID)
		return err
	}

	return nil

}

// возвращает задачу по id
func (r *Repository) Get(id int) ([]byte, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	//проверяем существование файла
	filePath := filepath.Join(r.pathDir, strconv.Itoa(id)+".json")
	if tools.FileExists(filePath) {
		//в случае успеха открываем файл на чтение
		file, err := os.Open(filePath)
		if err != nil {
			log.Printf("Error open file task id: %d, for reading.\n", id)
			return nil, err
		}
		defer file.Close()

		//читаем файл
		jsonFile, err := io.ReadAll(file)
		if err != nil {
			log.Printf("Error read file task, id: %d\n", id)
			return nil, err
		}
		//в случае успеха возвращаем считанный файл
		return jsonFile, nil
	} else {
		log.Printf("Error, task with id: %d does not exist.\n", id)
		return nil, fmt.Errorf("task with id: %d does not exist.", id)
	}
}
