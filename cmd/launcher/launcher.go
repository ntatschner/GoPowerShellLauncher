package launcher

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
)

var TempFiles []string

func ExecutePowerShellProcess(finalProfile string, shellPath string) error {
	l.Logger.Info("Executing PowerShell process", "ShellPath", shellPath)
	tmpFile, tmperr := os.CreateTemp("", "encoded_command_*.ps1")
	if tmperr != nil {
		l.Logger.Error("Failed to create temporary file", "Error", tmperr)
		return tmperr
	}
	TempFiles = append(TempFiles, tmpFile.Name())

	tmpFile.WriteString(finalProfile)
	tmpFile.Close()
	command := fmt.Sprintf(
		"Start-Process -FilePath \"%s\" -ArgumentList \"-NoProfile -NoExit -File %s\"",
		shellPath, tmpFile.Name(),
	)
	l.Logger.Info("PowerShell command", "Command", command)
	cmd := exec.Command("powershell", "-Command", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}

	err := cmd.Run()
	if err != nil {
		l.Logger.Error("Failed to start PowerShell process", "Error", err)
		return err
	}
	l.Logger.Debug("PowerShell process started", "PID", cmd.Process.Pid)
	l.Logger.Info("PowerShell process started successfully")
	return nil
}
