package launcher

import (
	"fmt"
	"os"
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
	cmd := exec.Command("cmd", "/C", "start", "/b", "/wait", "powershell", "-Command", command)
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

func ExecuteInsideShell(encodedCmd string) error {
	l.Logger.Debug("Executing command inside shell")
	// Get caller shell path

	validShells := []string{"powershell", "pwsh"}
	var shell string
	for _, s := range validShells {
		if _, err := exec.LookPath(s); err == nil {
			l.Logger.Debug("Shell found", "Shell", s)
			shell = s
			break
		}
	}
	var shellerr error
	executable, shellerr := os.Executable()
	if shellerr != nil || executable != shell {
		l.Logger.Error("No valid shell found")
		return fmt.Errorf("no valid shell found")
	}
	l.Logger.Debug("Shell executable found", "Executable", executable)
	cmd := exec.Command(fmt.Sprintf("%s -EncodedCommand %s", shell, encodedCmd))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		fmt.Println("Error running command:", err)
	}
	return nil
}
