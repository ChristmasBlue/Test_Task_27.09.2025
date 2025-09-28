package parser

import (
	"encoding/json"
	"fmt"
	"log"
)

// создаём json
func ParsTaskToJson(model interface{}) ([]byte, error) {
	jsonFile, err := json.MarshalIndent(model, "", "	")
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
