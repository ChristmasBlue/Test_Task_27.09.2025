package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DirectoryForDownload string `yaml:"directory_for_download"`
	DirectoryTaskStatus  string `yaml:"directory_task_status"`
	DirectoryTasks       string `yaml:"directory_tasks"`
	DirectoryQueue       string `yaml:"directory_queue"`
	FileQueue            string `yaml:"file_queue"`
	AddressServer        string `yaml:"address_server"`
}

func GetAddiction(pathYaml string) (*Config, error) {
	if pathYaml == "" {
		// Относительный путь от cmd/main.go к config/config.yaml
		pathYaml = "../config/config.yaml"
	}
	data, err := os.ReadFile(pathYaml)
	if err != nil {
		log.Printf("Error reading config.yaml: %v", err)
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		log.Printf("Error parse YAML: %v", err)
		return nil, err
	}

	return config, nil
}
