package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"test_task/internal/config"
	"test_task/internal/constants"
	"test_task/internal/handler"
	"test_task/internal/manager"
	"test_task/internal/queue"
	"test_task/internal/repository"
	"test_task/internal/services"
	"test_task/internal/worker"
	"test_task/pkg/tools"
	"time"
)

func Run() {
	var wg sync.WaitGroup

	//подключаем среды окружения
	conf, err := config.GetAddiction(constants.FileYaml)
	if err != nil {
		log.Printf("Error connecting environment variable: %v\n", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	//запускаем в отдельной горутине ожидание команды на остановку сервиса
	go tools.StopContext(&ctx, cancel)
	defer cancel()

	handlerTask, err := handler.NewHandler(conf.DirectoryForDownload)
	if err != nil {
		log.Printf("Error create/open directory for download: %v\n", err)
		return
	}

	repo, err := repository.NewRepository(conf.DirectoryTaskStatus)
	if err != nil {
		log.Printf("Error create/open directory for repository: %v\n", err)
		return
	}

	queueTask, err := queue.NewQueue(conf.DirectoryQueue, conf.FileQueue, repo)
	if err != nil {
		log.Printf("Error create/open directory for queue: %v\n", err)
		return
	}

	man := manager.NewManager(queueTask, repo)

	workerTask := worker.NewWorker(handlerTask, queueTask, repo, &wg, ctx)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go workerTask.ProcessTask(i + 1)
	}

	http.HandleFunc("/add", services.HandlerAddTask(man))
	http.HandleFunc("/info", services.HandlerGetInfoTask(man))

	// Создаем HTTP сервер с поддержкой graceful shutdown
	server := &http.Server{
		Addr:    conf.AddressServer,
		Handler: nil, // использует DefaultServeMux
	}

	// Запускаем сервер в отдельной горутине
	serverErr := make(chan error, 1)
	go func() {
		fmt.Printf("Сервер запущен на http://%s\n", conf.AddressServer)
		fmt.Printf("Endpoint: POST http://%s/add\n", conf.AddressServer)
		fmt.Printf("Endpoint: GET http://%s/info\n", conf.AddressServer)
		fmt.Println("Остановите сервер сочетанием клавиш Ctrl+C")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// Ожидаем либо ошибку сервера, либо отмену контекста
	select {
	case err := <-serverErr:
		log.Printf("Ошибка сервера: %v\n", err)
		cancel() // Отменяем контекст при ошибке сервера
	case <-ctx.Done():
		log.Println("Получен сигнал остановки...")
	}

	// Graceful shutdown сервера
	log.Println("Останавливаем HTTP сервер...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Ошибка при остановке сервера: %v\n", err)
	}

	wg.Wait()
	log.Println("Сервер остановлен корректно")
}
