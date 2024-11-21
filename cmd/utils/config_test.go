package utils

import (
	"os"
	"testing"
)

func createTempConfigFile(content string) (string, error) {
	file, err := os.CreateTemp("", "config_test_*.json")
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := file.Write([]byte(content)); err != nil {
		return "", err
	}

	return file.Name(), nil
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		expected    *Config
		expectError bool
	}{
		{
			name:        "Valid config",
			fileContent: `{"csv_path": "/path/to/csv"}`,
			expected:    &Config{CsvPath: "/path/to/csv"},
			expectError: false,
		},
		{
			name:        "Invalid JSON",
			fileContent: `{"csv_path": "/path/to/csv"`,
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Empty config",
			fileContent: `{}`,
			expected:    &Config{CsvPath: ""},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath, err := createTempConfigFile(tt.fileContent)
			if err != nil {
				t.Fatalf("Failed to create temp config file: %v", err)
			}
			defer os.Remove(filePath)

			config, err := loadConfig(filePath)
			if (err != nil) != tt.expectError {
				t.Errorf("loadConfig() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if !tt.expectError && *config != *tt.expected {
				t.Errorf("loadConfig() = %v, expected %v", config, tt.expected)
			}
		})
	}
}
