package parser

import (
	"encoding/json"
	"fmt"
	"log"
	"test_task/domain/models"
	"test_task/pkg/tools"
)

// парсим json
func ParsJsonToTask(jsonFile []byte) (*models.Task, error) {
	task := tools.NewTask()
	err := json.Unmarshal(jsonFile, task)
	if err != nil {
		log.Println("Error create new Task model, during parsing.")
		return nil, err
	}
	return task, nil
}

// создаём json
func ParsTaskToJson(task *models.Task) ([]byte, error) {
	jsonFile, err := json.MarshalIndent(task, "", "	")
	if err != nil {
		log.Println("Error create Json file, during parsing.")
		return nil, err
	}
	if len(jsonFile) == 0 {
		log.Println("Error, Json file is empty.")
		return nil, fmt.Errorf("file is empty")
	}
	return jsonFile, nil
}
