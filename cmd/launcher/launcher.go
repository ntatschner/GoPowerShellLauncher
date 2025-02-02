package launcher

import (
	"fmt"
	"os/exec"
	"syscall"

	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
)

func ExecutePowerShellProcess(encodedCommand string, shellPath string) error {
	l.Logger.Info("Executing PowerShell process", "ShellPath", shellPath)
	command := fmt.Sprintf(
		"Start-Process -FilePath \"%s\" -ArgumentList \"-NoProfile -NoExit -EncodedCommand %s\"",
		shellPath, encodedCommand,
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
