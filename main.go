/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/ntatschner/GoPowerShellLauncher/cmd"
	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/utils"
)

func main() {
	config, err := utils.LoadConfig()
	if err != nil {
		panic(err)
	}
	// Initialize the logger
	err = l.InitLogger(config.Logging.Path, config.Logging.File, config.Logging.Level)
	if err != nil {
		panic(err)
	}
	defer l.CloseLogger()
	l.Logger.Info("Starting..")
	cmd.Execute()
}
