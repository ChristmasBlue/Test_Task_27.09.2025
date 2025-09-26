package tools

import "test_task/domain/models"

func NewTask() *models.Task {
	return &models.Task{
		URLs:    make([]string, 0),
		Results: make(map[string]models.FileResult),
	}
}
