package utils

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"

	"github.com/spf13/viper"
)

type Profile struct {
	Name string `mapstructure:"name"`
	Path string `mapstructure:"path"`
}

type Shortcut struct {
	ID          string    `mapstructure:"id"`
	Name        string    `mapstructure:"name"`
	Destination string    `mapstructure:"destination"`
	Profiles    []Profile `mapstructure:"profiles"`
}

type Config struct {
	Profile struct {
		Path      string `mapstructure:"path"`
		Recursive bool   `mapstructure:"recursive"`
	} `mapstructure:"profile"`
	Logging struct {
		Path  string `mapstructure:"path"`
		File  string `mapstructure:"file"`
		Level string `mapstructure:"level"`
	} `mapstructure:"logging"`
	Shortcuts []Shortcut `mapstructure:"shortcuts"`
}

var UserConfigDir string

type ConfigStore struct {
	Path   string
	Exists bool
}

var ConfigStoreData []ConfigStore

var config *Config

func LoadConfig() (*Config, error) {
	if config != nil {
		return config, nil
	}
	usr, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("error getting current user: %w", err)
	}
	UserConfigDir := filepath.Join(usr.HomeDir, "Documents", "GoPowerShellLauncher")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(UserConfigDir)

	execPath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	execDir := filepath.Dir(execPath)

	configPaths := []string{
		".",
		UserConfigDir,
		execDir,
	}

	// Log the configuration files found
	for _, path := range configPaths {
		if path == "." {
			path, _ = os.Getwd()
		}
		configFile := filepath.Join(path, "config.yaml")
		if _, err := os.Stat(configFile); err == nil {
			ConfigStoreData = append(ConfigStoreData, ConfigStore{Path: configFile, Exists: true})
		} else {
			ConfigStoreData = append(ConfigStoreData, ConfigStore{Path: configFile, Exists: false})
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	config = &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return config, nil
}

func GenerateUniqueID() string {
	config, err := LoadConfig()
	if err != nil {
		return "1"
	}

	maxID := 0
	for _, shortcut := range config.Shortcuts {
		id, err := strconv.Atoi(shortcut.ID)
		if err == nil && id > maxID {
			maxID = id
		}
	}

	return strconv.Itoa(maxID + 1)
}
