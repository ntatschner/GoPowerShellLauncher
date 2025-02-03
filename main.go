//go:generate goversioninfo
/*
Copyright © 2024 Nigel Tatschner
*/
package main

import (
	"os"

	"github.com/ntatschner/GoPowerShellLauncher/cmd"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/launcher"
	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/utils"
)

func cleanupTempFiles() {
	for _, file := range launcher.TempFiles {
		os.Remove(file)
	}
}

var MousetrapHelpText = ""

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
	defer cleanupTempFiles()
	l.Logger.Info("Starting..")
	cmd.Execute()
}
