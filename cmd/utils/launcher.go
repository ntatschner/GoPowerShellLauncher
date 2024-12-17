package utils

import (
	"os"

	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
)

func MergeSelectedProfiles(selected []string) string {
	l.Logger.Info("Merging selected profiles", "Selected", selected)
	var merged string
	for i := range selected {
		content, err := GetProfileContent(selected[i])
		if err != nil {
			l.Logger.Warn("Error reading profile content", "Error", err)
			continue
		}
		merged += content + "\n"
	}
	return merged
}

// create temp file with merged profiles
func CreateTempFile(merged string) (string, error) {
	l.Logger.Info("Creating temp file", "Merged", merged)
	tempDir := os.TempDir()
	tempFile, err := os.CreateTemp(tempDir, "merged-*.ps1")
	if err != nil {
		l.Logger.Error("Failed to create temp file", "Error", err)
		return "", err
	}
	defer tempFile.Close()

	_, err = tempFile.WriteString(merged)
	if err != nil {
		l.Logger.Error("Failed to write to temp file", "Error", err)
		return "", err
	}

	l.Logger.Info("Temp file created successfully", "TempFile", tempFile.Name())
	return tempFile.Name(), nil
}
