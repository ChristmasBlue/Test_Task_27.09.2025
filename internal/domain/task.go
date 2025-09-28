package domain

import (
	"test_task/internal/constants"
	"time"
)

type Task struct {
	ID        int                   `json:"id"`         //Id заказа
	URLs      []string              `json:"urls"`       //список ссылок для скачивания файлов
	Status    constants.TaskStatus  `json:"status"`     //статус выполнения задачи
	CreatedAt time.Time             `json:"created_at"` //время создания задачи
	UpdatedAt time.Time             `json:"updated_at"` //время изменения задачи
	Results   map[string]FileResult `json:"results"`    //детальная информация ссылок
}
