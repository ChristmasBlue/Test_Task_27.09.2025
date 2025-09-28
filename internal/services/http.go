package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"test_task/internal/domain"
	"test_task/internal/models"
	"test_task/pkg/parser"
)

type taskManager interface {
	AddTask(*models.AddTaskDto) (int, error)
	GetTask(int) (*domain.Task, error)
}

// добавление задания в обработку, запрос POST
func HandlerAddTask(management taskManager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Настраиваем заголовки CORS
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Обрабатываем OPTIONS запрос (для CORS)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Проверяем метод
		if r.Method != "POST" {
			sendError(w, "Только POST запросы разрешены", http.StatusMethodNotAllowed)
			log.Println("Error request is method not POST.")
			return
		}

		//читаем запрос
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			log.Println("Error reading request body.")
			return
		}
		defer r.Body.Close()

		//проверяем что тело запроса не пустое
		if len(data) == 0 {
			http.Error(w, "Empty request body", http.StatusBadRequest)
			log.Println("Empty request body.")
			return
		}

		//создаю модель Dto
		taskDto := &models.AddTaskDto{URLs: make([]string, 0)}
		err = json.Unmarshal(data, taskDto)
		if err != nil {
			http.Error(w, "Bad request body", http.StatusBadRequest)
			log.Printf("Bad request body: %s\n", data)
			return
		}

		//передаю менеджеру полученную модель
		id, err := management.AddTask(taskDto)
		if err != nil {
			http.Error(w, "Error processing task", http.StatusBadRequest)
			log.Printf("Error processing task: %v", err)
			return
		}

		//отправляем успешный ответ
		response := models.Response{
			Status:  "success",
			ID:      strconv.Itoa(id),
			Message: fmt.Sprintf("Задача успешно получена."),
			Details: "Сервис работает в режиме получения запросов",
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// получение информации о задаче, запрос GET
func HandlerGetInfoTask(management taskManager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Настраиваем заголовки CORS
		//w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Обрабатываем OPTIONS запрос (для CORS)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Проверяем метод
		if r.Method != "GET" {
			sendError(w, "Только GET запросы разрешены", http.StatusMethodNotAllowed)
			return
		}

		//получаем id из запроса
		queryParam := r.URL.Query()

		idQuery := queryParam.Get("id")

		//проверяем что тело запроса не пустое
		if idQuery == "" {
			http.Error(w, "Empty request body", http.StatusBadRequest)
			log.Println("Empty request body.")
			return
		}

		id, err := strconv.Atoi(idQuery)
		if err != nil {
			http.Error(w, "Invalid Id.", http.StatusBadRequest)
			log.Printf("Invalid Id: %s\n", idQuery)
			return
		}

		//передаю менеджеру полученный json файл
		//TODO get id from data
		task, err := management.GetTask(id)
		if err != nil {
			http.Error(w, "Error processing task", http.StatusBadRequest)
			log.Println("Error processing task.")
			return
		}

		//парсим задачу
		req, err := parser.ParsTaskToJson(task)
		if err != nil {
			http.Error(w, "Error processing task", http.StatusBadRequest)
			log.Printf("Error parsing task: %v\n", req)
			return
		}
		//TODO convert task to json
		//отправляем успешный ответ
		response := models.Response{
			Status:  "success",
			Message: fmt.Sprintf("Информация о задаче"),
			Data:    req,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// Вспомогательная функция для ошибок
func sendError(w http.ResponseWriter, message string, statusCode int) {
	response := models.Response{
		Status:  "error",
		Message: message,
	}
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
