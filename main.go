/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/ntatschner/GoPowerShellLauncher/cmd"
	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
)

func main() {
	// Initialize the logger
	l.InitLogger()
	l.Logger.Info("Starting..")
	cmd.Execute()
}
