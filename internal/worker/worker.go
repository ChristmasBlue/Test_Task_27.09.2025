package worker

import (
	"context"
	"log"
	"sync"
	"test_task/internal/constants"
	"test_task/internal/domain"
	"time"
)

type taskRepository interface {
	Save(*domain.Task) error
}

type taskHandler interface {
	Download(string) (string, error)
}

type taskQueue interface {
	Dequeue() (*domain.Task, error)
	TaskCompleated(int) error
	IsEmpty() bool
}

type Worker struct {
	repo    taskRepository
	handler taskHandler
	queue   taskQueue
	wg      *sync.WaitGroup
	ctx     context.Context
}

// конструктор
func NewWorker(handler taskHandler, queue taskQueue, repo taskRepository, wg *sync.WaitGroup, ctx context.Context) *Worker {
	return &Worker{
		repo:    repo,
		handler: handler,
		queue:   queue,
		wg:      wg,
		ctx:     ctx,
	}
}

func (w *Worker) ProcessTask(idWorker int) {
	defer w.wg.Done()

	for {

		select {
		case <-w.ctx.Done():
			log.Printf("Воркер %d закончил работу\n", idWorker)
			return
		default:

			//проверям очередь на пустоту
			if !w.queue.IsEmpty() {

				//берём из очереди задачу
				task, err := w.queue.Dequeue()
				if err != nil {
					log.Println("Error pop task from queue.")
					continue
				}

				//выставляем статус задачи "Выполняется"
				task.Status = constants.StatusRunning
				//сразу сохраняем состояние, в случае ошибки просто не обрабатываем задачу
				err = w.repo.Save(task)
				if err != nil {
					log.Printf("Error saving task, id: %d\n", task.ID)
					continue
				}

				//перед обработкой задачи проверяем контекст
				if w.isContextDone() {
					log.Printf("Воркер %d закончил работу\n", idWorker)
					return
				}

				//проходим циклом по всем ссылкам, по которым нужно скачать файл
				for _, url := range task.URLs {

					//проверяем обработана ли ссылка или нет
					//если обработана пропускаем
					if _, exist := task.Results[url]; exist {
						continue
					}

					//перед обработкой URL проверяем контекст на остановку
					if w.isContextDone() {
						log.Printf("Воркер %d закончил работу\n", idWorker)
						return
					}
					//результат выполнения по ссылке
					res := domain.FileResult{
						URL:   url,
						Path:  "",
						Error: nil,
						Ok:    true,
					}

					//скачиваем файл
					filePath, err := w.handler.Download(url)
					if err != nil {
						//в случае недудачи, сохраняем ошибку
						log.Printf("Error download file, url: %s\n", url)
						res.Error = err
					} else {
						//в случае успеха, сохраняем путь к файлу
						res.Path = filePath
					}

					//добавляем результат обработки ссылки в статус задачи
					task.Results[url] = res

					//обновляем время изменения
					task.UpdatedAt = time.Now()

					//сохраняем данные о статусе выполнения задания
					err = w.repo.Save(task)
					if err != nil {
						log.Printf("Error saving task, id: %d\n", task.ID)
						continue
					}

					log.Printf("Processing compleated, url: %s\n", url)
				}

				//выставляем статус задачи "Выполнено"
				task.Status = constants.StatusCompleted
				err = w.repo.Save(task)
				if err != nil {
					log.Printf("Error saving task, id: %d\n", task.ID)
					continue
				}

				//отправляем id задания которое выполнилось
				err = w.queue.TaskCompleated(task.ID)
				if err != nil {
					log.Printf("Error deleted from queue task, id: %d\n", task.ID)
				}
			}

			//таймаут воркера
			time.Sleep(5 * time.Second)
		}
	}
}

// проверка контекста на команду остановки
func (w *Worker) isContextDone() bool {
	select {
	case <-w.ctx.Done():
		return true
	default:
		return false
	}
}
