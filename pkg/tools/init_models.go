package tools

import "test_task/internal/domain"

func NewTask() *domain.Task {
	return &domain.Task{
		URLs:    make([]string, 0),
		Results: make(map[string]domain.FileResult),
	}
}
