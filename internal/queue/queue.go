package queue

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"test_task/internal/domain"
	"test_task/pkg/tools"
)

type taskRepository interface {
	Get(int) (*domain.Task, error) //возвращает состояние задачи по id
	Save(*domain.Task) error       //созадёт/сохраняет задачу
}

type Queue struct {
	filePath string
	queue    []*domain.Task
	queueID  []int
	mu       sync.Mutex
	repo     taskRepository
}

// конструктор
func NewQueue(pathDir, filename string, repo taskRepository) (*Queue, error) {
	filePathName := filepath.Join(pathDir, filename)

	q := &Queue{
		filePath: filePathName,
		queue:    make([]*domain.Task, 0),
		queueID:  make([]int, 0),
		repo:     repo,
	}
	//создаём/проверяем директорию
	err := tools.CreateDir(pathDir)
	if err != nil {
		log.Println("Error, can't create directory for queue.")
		return nil, err
	}
	//создаём/проверяем и восстанавливаем очередь файл
	if tools.FileExists(filePathName) {
		//прочитали очерёдность id задач
		ids, err := restorationIdTaskInQueue(filePathName)
		if err != nil {
			log.Println("Error open file to restoration.")

			//если не удалось открыть файл и прочитать из него данные, просто пересоздаём его
			file, err := os.Create(filePathName)
			if err != nil {
				log.Println("Error create file to queue.")
				return nil, err
			}
			file.Close()
		} else {
			//проходим циклом по срезу id для добавления в очередь
			for _, id := range ids {
				//получаем статус задания из репозитория по id
				task, err := q.repo.Get(id)
				if err != nil {
					log.Println("Error read file from repository.")
					return nil, err
				}

				//добавляем задание в очередь
				q.queue = append(q.queue, task)
				q.queueID = append(q.queueID, int(id))
			}
		}

	} else {
		file, err := os.Create(filePathName)
		if err != nil {
			log.Println("Error create file to queue.")
			return nil, err
		}
		file.Close()
	}

	return q, nil
}

// добавление в очередь
func (q *Queue) Enqueue(task *domain.Task) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	//открытие/создание файла для записи
	file, err := os.OpenFile(q.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error, can't create or open file: %s\n", q.filePath)
		return err
	}
	defer file.Close()

	//дозапись в файл
	_, err = file.WriteString(strconv.Itoa(task.ID) + "\n")
	if err != nil {
		log.Printf("Error writing in file queue: %s\n", q.filePath)
		return err
	}

	//создаём новую задачу в репозитории
	err = q.repo.Save(task)
	if err != nil {
		log.Printf("Error create new task in repository, id: %d\n", task.ID)
		return err
	}

	//добавление элемента в конец очереди
	q.queue = append(q.queue, task)
	q.queueID = append(q.queueID, task.ID)

	return nil
}

// получение из очереди с дальнейшим очищением из очереди
func (q *Queue) Dequeue() (*domain.Task, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	//проверяем очередь на пустоту
	if len(q.queue) == 0 {
		log.Println("Error, queue is empty.")
		return nil, fmt.Errorf("queue is empty")
	}

	//получаем задачу из очреди
	task := q.queue[0]

	//удаляем задачу из очереди
	q.queue = q.queue[1:]

	return task, nil
}

// проверка на пустоту
func (q *Queue) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.queue) == 0

}

// TaskCompleated вызывается после выполнения задачи, очищает из файла очереди выполненную задачу
func (q *Queue) TaskCompleated(id int) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	//проходим циклом и ищем задачу (в основном id задачи лежит в начале, поэтому цикл будет выполняться быстро)
	for i, queueId := range q.queueID {
		if id == queueId {
			//удаляем элемент из очереди
			q.queueID = append(q.queueID[:i], q.queueID[i+1:]...)
			//пересохраняем очередь в файл
			err := q.saveFullQueue()
			if err != nil {
				log.Printf("Error rewriting queue in file: %s\n", q.filePath)
				return err
			}
			break
		}
	}
	return nil
}

// saveFullQueue пересохраняет всю очередь в файл
func (q *Queue) saveFullQueue() error {
	//пересоздаём файл очереди
	file, err := os.Create(q.filePath)
	if err != nil {
		log.Println("Error recreation file queue.")
		return err
	}
	defer file.Close()

	//перезаписываем все элементы очереди
	for _, id := range q.queueID {
		//err = binary.Write(file, binary.LittleEndian, int64(id))
		_, err = file.WriteString(strconv.Itoa(id) + "\n")
		if err != nil {
			log.Printf("Error rewriting element queue: %d\n", id)
			return err
		}
	}

	return nil
}

// восстановление очереди id из файла
func restorationIdTaskInQueue(filename string) ([]int, error) {
	//открываем файл для чтения
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error open file to restore queue, path file: %s\n", filename)
		return nil, err
	}
	defer file.Close()

	//создаём буфер для чтения из файла
	scanner := bufio.NewScanner(file)
	ids := make([]int, 0)

	//читаем построчно из файла
	for scanner.Scan() {
		//конвертируем строку в int
		id, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Printf("Error reading queue. error: %v\n", err)
			return nil, err
		}
		//добавляем в очередь
		ids = append(ids, id)
	}

	//проверяем чтение на ошибки
	if err = scanner.Err(); err != nil {
		log.Printf("Error scanner queue error: %v\n", err)
		return nil, err
	}

	return ids, nil
}
