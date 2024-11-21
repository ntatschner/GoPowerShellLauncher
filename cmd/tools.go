package cmd

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func validatePath(path string) error {
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
	switch shellVersion {
	case "pwsh", "powershell":
		return nil
	default:
		return fmt.Errorf("invalid shell version: %s", shellVersion)
	}
}
