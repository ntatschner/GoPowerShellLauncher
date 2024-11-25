package logger

import (
	"io"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

var Logger *log.Logger
var logFile *os.File

func InitLogger() {
	var err error
	// Open a file for writing logs
	logFile, err = os.OpenFile("GoPowerShellLauncher.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	// Create a multi-writer to write logs to both the file and the standard output
	//multiWriter := io.MultiWriter(os.Stdout, logFile)
	writer := io.Writer(logFile)

	// Create a new logger and set its output to the multi-writer
	Logger = log.New(writer)
	Logger.SetOutput(writer)
	Logger.SetPrefix("GoPowerShellLauncher 🤖:")
	Logger.SetTimeFormat(time.Kitchen)
	Logger.SetReportTimestamp(true)
	Logger.SetReportCaller(true)
	Logger.Info("Logger initialized")
}

// CloseLogger closes the log file
func CloseLogger() {
	if logFile != nil {
		logFile.Close()
	}
}
