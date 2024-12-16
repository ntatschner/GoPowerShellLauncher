package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	ProfilePath string `json:"profile_path"`
	Recursive   bool   `json:"recursive"`
	Logging     struct {
		LogPath  string `json:"log_path"`
		LogFile  string `json:"log_file"`
		LogLevel string `json:"log_level"`
	} `json:"logging"`
}

func LoadConfig() (*Config, error) {
	cwd, _ := os.Getwd()
	filePath := filepath.Join(cwd, "config.json")
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
