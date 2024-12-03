/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"os"

	"github.com/ntatschner/GoPowerShellLauncher/cmd"
	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/utils"
)

func main() {
	cwd, _ := os.Getwd()
	configPath := fmt.Sprintf("%s%c%s", cwd, os.PathSeparator, "config.json")
	config, err := utils.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}
	// Initialize the logger
	err = l.InitLogger(config.Logging.LogPath, config.Logging.LogFile, config.Logging.LogLevel)
	if err != nil {
		panic(err)
	}
	defer l.CloseLogger()
	l.Logger.Info("Starting..")
	cmd.Execute()
}
