package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf16"

	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
)

type HashValidator interface {
	ValidateHash(expectedHash, filePath string) (bool, error)
}

type DefaultHashValidator struct{}

func (d DefaultHashValidator) ValidateHash(expectedHash, filePath string) (bool, error) {
	return validateHash(expectedHash, filePath)
}

func validatePath(path string) (bool, error) {
	l.Logger.Info("Validating path", "Path", path)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, fmt.Errorf("path does not exist: %s", path)
	}
	if err != nil {
		return false, fmt.Errorf("error accessing path: %s", err)
	}
	return true, nil
}

func compareHashes(hash1, hash2 []byte) (bool, error) {
	l.Logger.Info("Comparing hashes", "Hash1", hex.EncodeToString(hash1), "Hash2", hex.EncodeToString(hash2))
	if len(hash1) != len(hash2) {
		return false, fmt.Errorf("hashes have different lengths")
	}
	for i := range hash1 {
		if hash1[i] != hash2[i] {
			return false, fmt.Errorf("hashes do not match")
		}
	}
	l.Logger.Info("Hashes match")
	return true, nil
}

func validateHash(expectedHash, filePath string) (bool, error) {
	l.Logger.Info("Validating hash", "ExpectedHash", expectedHash, "FilePath", filePath)

	// Decode the expected hash from hex
	expectedHashBytes, err := hex.DecodeString(expectedHash)
	if err != nil {
		return false, fmt.Errorf("invalid expected hash: %v", err)
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return false, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Compute the SHA-256 hash of the file
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return false, fmt.Errorf("error hashing file: %v", err)
	}
	computedHash := hasher.Sum(nil)

	// Compare the computed hash with the expected hash
	_, err = compareHashes(expectedHashBytes, computedHash)
	if err != nil {
		return false, fmt.Errorf("hash mismatch: expected %x, got %x", expectedHashBytes, computedHash)
	}

	return true, nil
}

func validateShellVersion(shellVersion string) (bool, error) {
	l.Logger.Info("Validating shell version", "ShellVersion", shellVersion)
	shellVersion = strings.ToLower(shellVersion)
	switch shellVersion {
	case "pwsh", "powershell", "all":
		l.Logger.Info("Shell version is valid")
		return true, nil
	}
	return false, fmt.Errorf("invalid shell version: %s", shellVersion)
}

func validateDescription(description string) (bool, error) {
	l.Logger.Info("Validating description", "Description", description)
	if len(description) > 100 {
		return false, fmt.Errorf("description is too long (max 100 characters)")
	}

	l.Logger.Info("Description is valid")
	return true, nil
}

func GetProfileContent(path string) (string, error) {
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

func LoadProfileContent(profilePath string) (string, error) {
	l.Logger.Info("Loading profile content", "ProfilePath", profilePath)
	file, err := os.Open(profilePath)
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

func NormalizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func EncodeCommand(command string) (string, error) {
	// Convert the string to UTF-16LE
	utf16LE := utf16.Encode([]rune(command))
	buf := new(bytes.Buffer)
	for _, r := range utf16LE {
		if err := binary.Write(buf, binary.LittleEndian, r); err != nil {
			return "", fmt.Errorf("failed to encode command to UTF-16LE: %v", err)
		}
	}

	// Base64 encode the UTF-16LE bytes
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	if encoded == "" {
		return "", fmt.Errorf("failed to encode command")
	}
	return encoded, nil
}
