package logger

import (
	"io"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

var Logger *log.Logger
var logFile *os.File

func InitLogger(logPath, logFileName, logLevel string) error {
	var err error
	if logFileName == "" {
		logFileName = "GoPowerShellLauncher.log"
	}
	if logPath != "" {
		if err = os.MkdirAll(logPath, os.ModePerm); err != nil {
			return err
		}
	} else {
		logPath, err = os.Getwd()
		if err != nil {
			return err
		}
	}
	// Open a file for writing logs
	logFile, err = os.OpenFile(logPath+string(os.PathSeparator)+logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	// Create a multi-writer to write logs to both the file and the standard output
	//writer := io.MultiWriter(os.Stdout, logFile)
	writer := io.Writer(logFile)

	// Create a new logger and set its output to the multi-writer
	Logger = log.New(writer)
	level, logerr := log.ParseLevel(logLevel)
	if logerr != nil {
		Logger.Errorf("Failed to parse log level: %v", logerr)
		level = log.InfoLevel
	}
	Logger.SetLevel(level)
	Logger.SetOutput(writer)
	Logger.SetPrefix("GoPowerShellLauncher 🤖:")
	Logger.SetTimeFormat(time.Kitchen)
	Logger.SetReportTimestamp(true)
	Logger.SetReportCaller(true)
	Logger.Info("Logger initialized")
	return nil
}

func CloseLogger() {
	if logFile != nil {
		logFile.Close()
	}
}
