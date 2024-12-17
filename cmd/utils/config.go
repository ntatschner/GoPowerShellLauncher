package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

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
}

var config *Config

func LoadConfig() (*Config, error) {
	if config != nil {
		return config, nil
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME\\AppData\\Local\\GoPowerShellLauncher")
	viper.AddConfigPath("C:\\ProgramData\\GoPowerShellLauncher")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	config = &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return config, nil
}
