package utils

import (
	"encoding/json"
	"os"
)

type Config struct {
	CsvPath string `json:"csv_path"`
	Logging struct {
		LogPath  string `json:"log_path"`
		LogFile  string `json:"log_file"`
		LogLevel string `json:"log_level"`
	} `json:"logging"`
}

func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
