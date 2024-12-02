package launcher

import (
	"os"
	"os/exec"
	"syscall"

	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/utils"
)

func MergeSelectedProfiles(selected []string) string {
	l.Logger.Info("Merging selected profiles", "Selected", selected)
	var merged string
	for i := range selected {
		content, err := utils.GetProfileContent(selected[i])
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

func ExecutePowerShellProcess(encodedCommand string, shellPath string) error {
	l.Logger.Info("Executing PowerShell process", "ShellPath", shellPath)
	command := "Start-Process -FilePath " + shellPath + " -ArgumentList '-NoExit -EncodedCommand '" + encodedCommand + "''" + " -Wait" + " -NoNewWindow"
	l.Logger.Info("PowerShell command", "Command", command)
	cmd := exec.Command("cmd", "/C", "start", "/wait", shellPath, "-NoExit", "-Command", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}

	err := cmd.Run()
	l.Logger.Info("PowerShell process started", "PID", cmd.Process.Pid)
	if err != nil {
		l.Logger.Error("Failed to start PowerShell process", "Error", err)
		return err
	}
	l.Logger.Info("PowerShell process started successfully")
	return nil
}
