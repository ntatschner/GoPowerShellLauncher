package utils

import (
	"encoding/json"
	"io"
	"os"

	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
)

type Config struct {
	CsvPath string `json:"csv_path"`
}

func LoadConfig(filePath string) (*Config, error) {
	l.Logger.Info("Loading configuration file", "Path", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
