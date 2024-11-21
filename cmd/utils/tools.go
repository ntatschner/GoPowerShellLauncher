package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"

	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
)

func validatePath(path string) error {
	l.Logger.Info("Validating path", "Path", path)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", path)
	}
	if err != nil {
		return fmt.Errorf("error accessing path: %s", err)
	}
	return nil
}

func validateHash(hash string, filePath string) error {
	l.Logger.Info("Validating hash", "Hash", hash, "Path", filePath)
	// Check if the path is valid before validating the hash
	if err := validatePath(filePath); err != nil {
		return err
	}

	// Decode the base64-encoded hash
	decodedHash, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return fmt.Errorf("invalid base64 hash: %v", err)
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Compute the SHA-256 hash of the file
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return fmt.Errorf("error hashing file: %v", err)
	}
	computedHash := hasher.Sum(nil)

	// Compare the computed hash with the provided hash
	if !compareHashes(decodedHash, computedHash) {
		return fmt.Errorf("hash mismatch: expected %x, got %x", decodedHash, computedHash)
	}

	return nil
}

func compareHashes(hash1, hash2 []byte) bool {
	l.Logger.Info("Comparing hashes", "Hash1", hash1, "Hash2", hash2)
	if len(hash1) != len(hash2) {
		return false
	}
	for i := range hash1 {
		if hash1[i] != hash2[i] {
			return false
		}
	}
	return true
}

func validateShellVersion(shellVersion string) error {
	l.Logger.Info("Validating shell version", "ShellVersion", shellVersion)
	switch shellVersion {
	case "pwsh", "powershell":
		return nil
	default:
		return fmt.Errorf("invalid shell version: %s", shellVersion)
	}
}

func validateDescription(description string) error {
	l.Logger.Info("Validating description", "Description", description)
	if len(description) > 100 {
		return fmt.Errorf("description is too long (max 100 characters)")
	}
	return nil
}

func getProfileContent(path string) (string, error) {
	l.Logger.Info("Getting profile content", "Path", path)
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func mergeSelectedProfiles(selected map[int]struct{}) string {
	l.Logger.Info("Merging selected profiles", "Selected", selected)
	var merged string
	for i := range selected {
		content, err := getProfileContent(profiles[i].path)
		if err != nil {
			l.Logger.Warn("Error reading profile content", "Error", err)
			continue
		}
		merged += content + "\n"
	}
	return merged
}

func launchPowerShell(m model) tea.Cmd {
	l.Logger.Info("Launching PowerShell")
	return func() tea.Msg {
		merged := mergeSelectedProfiles(m.selected)
		// create temporary file to store the merged profile
		tmpFile, err := os.CreateTemp("", "merged_profile_*.ps1")
		if err != nil {
			logger.Errorf("Error creating temporary file", "Error", err)
			return shellClosed{}
		}
		_, err = tmpFile.Write([]byte(merged))
		if err != nil {
			logger.Errorf("Error writing to temporary file", "Error", err)
			tmpFile.Close()
			os.Remove(tmpFile.Name())
			return shellClosed{}
		}
		tmpFile.Close() // Close the file before running the PowerShell command

		cmd := exec.Command("cmd", "/C", "start", "/wait", "powershell", "-NoExit", "-File", tmpFile.Name())
		cmd.SysProcAttr = &syscall.SysProcAttr{
			CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}

		err = cmd.Start()
		if err != nil {
			logger.Errorf("Error starting PowerShell command", "Error", err)
			os.Remove(tmpFile.Name())
			return shellClosed{}
		}
		cmd.Wait()
		// Remove the temporary file after starting the command
		os.Remove(tmpFile.Name())
		return shellClosed{}
	}
}
