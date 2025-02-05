package utils

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"

	shortcut "github.com/nyaosorg/go-windows-shortcut"
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
	log.Printf("Loading config file")
	if config != nil {
		return config, nil
	}
	usr, err := user.Current()
	if err != nil {
		log.Printf("Error getting current user: %v", err)
		return nil, fmt.Errorf("error getting current user: %w", err)
	}
	UserConfigDir := filepath.Join(usr.HomeDir, "Documents", "GoPowerShellLauncher")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(UserConfigDir)

	exe, exeerr := os.Executable()
	log.Printf("Executable: %s", exe)
	var exeDir string
	if exeerr != nil {
		log.Printf("Error getting executable: %v", exeerr)
		return nil, fmt.Errorf("error getting executable: %w", exeerr)
	}
	if filepath.Ext(exe) == ".lnk" {
		log.Printf("Shortcut detected: %s", exe)
		exePath, _, direrr := shortcut.Read(exe)
		if direrr != nil {
			log.Printf("Error reading shortcut: %v", direrr)
			return nil, fmt.Errorf("error reading shortcut: %w", direrr)
		}
		if exePath != "" {
			log.Printf("Shortcut path: %s", exePath)
			exeDir = filepath.Dir(exePath)
		} else {
			log.Printf("Shortcut path is empty: %s", exePath)
			exeDir = ""
		}
	} else {
		exeDir = filepath.Dir(exe)
	}

	configPaths := []string{
		".",
		UserConfigDir,
		exeDir,
	}

	// Log the configuration files found
	for _, path := range configPaths {
		log.Printf("Checking config directory: %s", path)
		if path == "." {
			log.Println("Getting current working directory")
			path, _ = os.Getwd()
		}
		configFile := filepath.Join(path, "config.yaml")
		log.Printf("Checking config file: %s", configFile)
		if _, err := os.Stat(configFile); err == nil {
			log.Printf("Config file exists: %s", configFile)
			ConfigStoreData = append(ConfigStoreData, ConfigStore{Path: configFile, Exists: true})
		} else {
			log.Printf("Config file doesn't exist: %s", configFile)
			ConfigStoreData = append(ConfigStoreData, ConfigStore{Path: configFile, Exists: false})
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %v", err)
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	config = &Config{}
	if err := viper.Unmarshal(config); err != nil {
		log.Printf("Unable to decode into struct: %v", err)
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}
	log.Printf("Config file loaded: %+v", config)
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
