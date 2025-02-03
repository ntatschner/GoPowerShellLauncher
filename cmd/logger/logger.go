package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/charmbracelet/log"
)

var Logger *log.Logger
var logFile *os.File

const maxLogSize = 10 * 1024 * 1024 // 10 MB
const maxLogFiles = 2

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

func CheckLogSize(logPath, logFileName string) error {
	fileInfo, err := logFile.Stat()
	if err != nil {
		return err
	}
	if fileInfo.Size() >= maxLogSize {
		// Rotate the log file
		backupName := fmt.Sprintf("%s.%s", logFileName, time.Now().Format("20060102T150405"))
		if err := os.Rename(filepath.Join(logPath, logFileName), filepath.Join(logPath, backupName)); err != nil {
			return err
		}
		// Open a new log file
		logFile, err = os.OpenFile(filepath.Join(logPath, logFileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		// Update the logger's output
		writer := io.Writer(logFile)
		Logger.SetOutput(writer)

		// Remove old log files if there are more than maxLogFiles
		if err := RemoveOldLogFiles(logPath, logFileName); err != nil {
			return err
		}
	}
	return nil
}

func RemoveOldLogFiles(logPath, logFileName string) error {
	files, err := os.ReadDir(logPath)
	if err != nil {
		return err
	}

	var logFiles []os.FileInfo
	for _, file := range files {
		if file.Name() == logFileName || filepath.Ext(file.Name()) == ".log" {
			fileInfo, err := file.Info()
			if err != nil {
				return err
			}
			logFiles = append(logFiles, fileInfo)
		}
	}

	if len(logFiles) <= maxLogFiles {
		return nil
	}

	sort.Slice(logFiles, func(i, j int) bool {
		return logFiles[i].ModTime().After(logFiles[j].ModTime())
	})

	for _, file := range logFiles[maxLogFiles:] {
		if err := os.Remove(filepath.Join(logPath, file.Name())); err != nil {
			return err
		}
	}

	return nil
}
