//go:build windows
// +build windows

package utils

import (
	"fmt"
	"os/exec"
	"syscall"

	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
)

func executeCommandWithPowershell(encodedCmd string) error {
	l.Logger.Debug("Executing command with PowerShell")
	// Get path of powershell.exe
	powershellPath, err := exec.LookPath("powershell")
	if err != nil {
		l.Logger.Error("Failed to find PowerShell executable", "Error", err)
		return err
	}
	l.Logger.Debug("PowerShell executable found", "Path", powershellPath)
	command := fmt.Sprintf(
		"Start-Process -FilePath \"%s\" -ArgumentList \"-NoProfile -NonInteractive -WindowStyle Hidden -EncodedCommand %s\"",
		powershellPath, encodedCmd,
	)
	l.Logger.Debug("PowerShell command", "Command", command)
	cmd := exec.Command("cmd", "/C", "start", "/b", "/wait", "powershell", "-Command", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}
	exerr := cmd.Run()
	if exerr != nil {
		l.Logger.Error("Failed to start PowerShell process", "Error", err)
		return err
	}
	l.Logger.Debug("PowerShell process started", "PID", cmd.Process.Pid)
	l.Logger.Info("PowerShell process started successfully")
	return nil
}
