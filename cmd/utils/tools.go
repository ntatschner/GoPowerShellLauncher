package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
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

func compareHashes(hash1, hash2 []byte) bool {
	l.Logger.Info("Comparing hashes", "Hash1", hex.EncodeToString(hash1), "Hash2", hex.EncodeToString(hash2))
	if len(hash1) != len(hash2) {
		return false
	}
	for i := range hash1 {
		if hash1[i] != hash2[i] {
			return false
		}
	}
	l.Logger.Info("Hashes match")
	return true
}

func validateHash(expectedHash, filePath string) error {
	l.Logger.Info("Validating hash", "ExpectedHash", expectedHash, "FilePath", filePath)

	// Decode the expected hash from hex
	expectedHashBytes, err := hex.DecodeString(expectedHash)
	if err != nil {
		return fmt.Errorf("invalid expected hash: %v", err)
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

	// Compare the computed hash with the expected hash
	if !compareHashes(expectedHashBytes, computedHash) {
		return fmt.Errorf("hash mismatch: expected %x, got %x", expectedHashBytes, computedHash)
	}

	return nil
}

func validateShellVersion(shellVersion string) error {
	l.Logger.Info("Validating shell version", "ShellVersion", shellVersion)
	shellVersion = strings.ToLower(shellVersion)
	switch shellVersion {
	case "pwsh", "powershell", "all":
		l.Logger.Info("Shell version is valid")
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

	l.Logger.Info("Description is valid")
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

func MergeSelectedProfiles(selected map[int]struct{}) string {
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

func loadAvailableShells() []list.Item {
	l.Logger.Info("Loading available shells")
	var items []list.Item
	for _, shell := range shells {
		items = append(items, shell)
	}
	return items
}
