package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"test_task/domain/models"
)

type taskManager interface {
	AddTask([]byte) error
	GetTask([]byte) ([]byte, error)
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

		//передаю менеджеру полученный json файл
		err = management.AddTask(data)
		if err != nil {
			http.Error(w, "Error processing task", http.StatusBadRequest)
			log.Printf("Error processing task: %v", err)
			return
		}

		//отправляем успешный ответ
		response := models.Response{
			Status:  "success",
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
		w.Header().Set("Content-Type", "application/json")
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

		//передаю менеджеру полученный json файл
		req, err := management.GetTask(data)
		if err != nil {
			http.Error(w, "Error processing task", http.StatusBadRequest)
			log.Println("Error processing task.")
			return
		}

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
