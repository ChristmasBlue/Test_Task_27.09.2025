package manager

import (
	"fmt"
	"log"
	"sync"
	"test_task/internal/constants"
	"test_task/internal/domain"
	"test_task/internal/models"
	"test_task/pkg/tools"
	"time"
)

type taskQueue interface {
	Enqueue(*domain.Task) error
}

type taskRepository interface {
	Save(*domain.Task) error
	Get(int) (*domain.Task, error)
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
func (m *TaskManager) GetTask(taskId int) (*domain.Task, error) {
	// Ищем задачу в репозитории
	task, err := m.repo.Get(taskId)
	if err != nil {
		log.Printf("Error searching in repository task id %d: %v", taskId, err)
		return nil, err
	}

	return task, nil
}

// получение задачи, запрос POST
func (m *TaskManager) AddTask(addTask *models.AddTaskDto) (int, error) {
	//получаем из запроса задачу
	//парсим в структуру
	task := tools.NewTask()
	task.URLs = addTask.URLs

	//проверка URL
	if len(task.URLs) == 0 {
		log.Println("Error urls is empty.")
		return 0, fmt.Errorf("error urls is empty")
	}

	//добавляем необходимую информацию о задании
	task.CreatedAt = time.Now()
	task.Status = constants.StatusPending
	task.ID = int(time.Now().UnixMilli())

	//добавляем задание в очередь
	err := m.queue.Enqueue(task)
	if err != nil {
		log.Println("Error adding task in queue.")
		return 0, err
	}

	return task.ID, nil
}
