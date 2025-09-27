package manager

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"test_task/domain/constants"
	"test_task/domain/models"
	"test_task/pkg/tools"
	"time"
)

type taskQueue interface {
	Enqueue(*models.Task) error
}

type taskRepository interface {
	Save(*models.Task) error
	Get(int) ([]byte, error)
}

type TaskManager struct {
	queue taskQueue
	repo  taskRepository
	wg    *sync.WaitGroup
}

// конструктор
func NewManager(queue taskQueue, repo taskRepository) *TaskManager {
	return &TaskManager{
		queue: queue,
		repo:  repo,
	}
}

// получение задачи, запрос GET
func (m *TaskManager) GetTask(taskId []byte) ([]byte, error) {
	//получаем из запроса ID
	id := models.RequestIdTask{}
	err := json.Unmarshal(taskId, &id)
	if err != nil {
		log.Printf("Error unmarshaling json task: %s", string(taskId))
		return nil, err
	}

	// Ищем задачу в репозитории
	task, err := m.repo.Get(id.ID)
	if err != nil {
		log.Printf("Error searching in repository task id %d: %v", id.ID, err)
		return nil, err
	}

	return task, nil
}

// получение задачи, запрос POST
func (m *TaskManager) AddTask(taskJson []byte) error {
	//получаем из запроса задачу
	//парсим в структуру
	task := tools.NewTask()
	err := json.Unmarshal(taskJson, task)
	if err != nil {
		log.Printf("Error unmarshaling json task: %s", string(taskJson))
		return err
	}

	// Простая валидация
	if task.ID == 0 {
		log.Println("Error id is empty.")
		return fmt.Errorf("error id is empty")
	}

	//проверка URL
	if len(task.URLs) == 0 {
		log.Println("Error urls is empty.")
		return fmt.Errorf("error urls is empty")
	}

	//добавляем необходимую информацию о задании
	task.CreatedAt = time.Now()
	task.Status = constants.StatusPending

	//добавляем задание в очередь
	err = m.queue.Enqueue(task)
	if err != nil {
		log.Println("Error adding task in queue.")
		return err
	}

	return nil
}
