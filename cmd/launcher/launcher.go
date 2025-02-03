package launcher

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"

	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
)

func ExecutePowerShellProcess(encodedCommand string, shellPath string) error {
	l.Logger.Info("Executing PowerShell process", "ShellPath", shellPath)

	// Create a temporary file
	tmpFile, err := ioutil.TempFile("", "ps_script_*.ps1")
	if err != nil {
		l.Logger.Error("Failed to create temporary file", "Error", err)
		return err
	}
	defer os.Remove(tmpFile.Name())

	// Write the encoded command to the temporary file
	scriptContent := fmt.Sprintf(
		"$command = [System.Text.Encoding]::Unicode.GetString([System.Convert]::FromBase64String('%s'))\nInvoke-Expression $command",
		encodedCommand,
	)
	if _, err := tmpFile.Write([]byte(scriptContent)); err != nil {
		l.Logger.Error("Failed to write to temporary file", "Error", err)
		return err
	}
	if err := tmpFile.Close(); err != nil {
		l.Logger.Error("Failed to close temporary file", "Error", err)
		return err
	}

	// Execute the PowerShell script from the temporary file
	command := fmt.Sprintf(
		"Start-Process -FilePath \"%s\" -ArgumentList \"-NoProfile -NoExit -File %s\"",
		shellPath, tmpFile.Name(),
	)
	l.Logger.Info("PowerShell command", "Command", command)
	cmd := exec.Command("powershell", "-Command", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}

	err = cmd.Run()
	if err != nil {
		l.Logger.Error("Failed to start PowerShell process", "Error", err)
		return err
	}
	l.Logger.Debug("PowerShell process started", "PID", cmd.Process.Pid)
	l.Logger.Info("PowerShell process started successfully")
	return nil
}
