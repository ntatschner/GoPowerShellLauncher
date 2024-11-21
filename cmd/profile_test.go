package cmd

import (
	"crypto/sha256"
	"encoding/base64"
	"os"
	"testing"
)

func createTempFileWithContent(content string) (string, error) {
	file, err := os.CreateTemp("", "testfile_*.ps1")
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := file.Write([]byte(content)); err != nil {
		return "", err
	}

	return file.Name(), nil
}

func generateBase64Hash(content string) string {
	hasher := sha256.New()
	hasher.Write([]byte(content))
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}

func TestLoadProfile(t *testing.T) {
	// Create temporary files with known content
	file1Content := "content1"
	file2Content := "content2"
	file1Path, err := createTempFileWithContent(file1Content)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(file1Path)

	file2Path, err := createTempFileWithContent(file2Content)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(file2Path)

	// Generate base64-encoded hashes for the files
	hash1 := generateBase64Hash(file1Content)
	hash2 := generateBase64Hash(file2Content)

	tests := []struct {
		name     string
		line     []string
		expected profile
	}{
		{
			name: "Valid profile",
			line: []string{file1Path, hash1, "pwsh"},
			expected: profile{
				path:                file1Path,
				hash:                hash1,
				shellVersion:        "pwsh",
				isValidPath:         true,
				isValidHash:         true,
				isValidShellVersion: true,
			},
		},
		{
			name: "Invalid path",
			line: []string{"", hash1, "pwsh"},
			expected: profile{
				path:                "",
				hash:                hash1,
				shellVersion:        "pwsh",
				isValidPath:         false,
				isValidHash:         false,
				isValidShellVersion: true,
			},
		},
		{
			name: "Invalid hash",
			line: []string{file1Path, "", "pwsh"},
			expected: profile{
				path:                file1Path,
				hash:                "",
				shellVersion:        "pwsh",
				isValidPath:         true,
				isValidHash:         false,
				isValidShellVersion: true,
			},
		},
		{
			name: "Invalid shell version",
			line: []string{file2Path, hash2, "invalidShell"},
			expected: profile{
				path:                file2Path,
				hash:                hash2,
				shellVersion:        "invalidShell",
				isValidPath:         true,
				isValidHash:         true,
				isValidShellVersion: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LoadProfile(tt.line)
			if result.path != tt.expected.path {
				t.Errorf("LoadProfile().path = %v, expected %v", result.path, tt.expected.path)
			}
			if result.hash != tt.expected.hash {
				t.Errorf("LoadProfile().hash = %v, expected %v", result.hash, tt.expected.hash)
			}
			if result.shellVersion != tt.expected.shellVersion {
				t.Errorf("LoadProfile().shellVersion = %v, expected %v", result.shellVersion, tt.expected.shellVersion)
			}
			if result.isValidPath != tt.expected.isValidPath {
				t.Errorf("LoadProfile().isValidPath = %v, expected %v", result.isValidPath, tt.expected.isValidPath)
			}
			if result.isValidHash != tt.expected.isValidHash {
				t.Errorf("LoadProfile().isValidHash = %v, expected %v", result.isValidHash, tt.expected.isValidHash)
			}
			if result.isValidShellVersion != tt.expected.isValidShellVersion {
				t.Errorf("LoadProfile().isValidShellVersion = %v, expected %v", result.isValidShellVersion, tt.expected.isValidShellVersion)
			}
		})
	}
}
