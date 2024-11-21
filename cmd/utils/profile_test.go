package utils

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
			line: []string{file1Path, hash1, "pwsh", "A valid description"},
			expected: profile{
				path:                file1Path,
				hash:                hash1,
				shellVersion:        "pwsh",
				description:         "A valid description",
				isValidPath:         true,
				isValidHash:         true,
				isValidShellVersion: true,
				isValidDescription:  true,
			},
		},
		{
			name: "Invalid path",
			line: []string{"", hash1, "pwsh", "A valid description"},
			expected: profile{
				path:                "",
				hash:                hash1,
				shellVersion:        "pwsh",
				description:         "A valid description",
				isValidPath:         false,
				isValidHash:         false,
				isValidShellVersion: true,
				isValidDescription:  true,
			},
		},
		{
			name: "Invalid hash",
			line: []string{file1Path, "", "pwsh", "A valid description"},
			expected: profile{
				path:                file1Path,
				hash:                "",
				shellVersion:        "pwsh",
				description:         "A valid description",
				isValidPath:         true,
				isValidHash:         false,
				isValidShellVersion: true,
				isValidDescription:  true,
			},
		},
		{
			name: "Invalid shell version",
			line: []string{file2Path, hash2, "invalidShell", "A valid description"},
			expected: profile{
				path:                file2Path,
				hash:                hash2,
				shellVersion:        "invalidShell",
				description:         "A valid description",
				isValidPath:         true,
				isValidHash:         true,
				isValidShellVersion: false,
				isValidDescription:  true,
			},
		},
		{
			name: "Invalid description",
			line: []string{file2Path, hash2, "pwsh", "A very long description that exceeds the maximum allowed length of 100 characters. This description should be considered invalid."},
			expected: profile{
				path:                file2Path,
				hash:                hash2,
				shellVersion:        "pwsh",
				description:         "A very long description that exceeds the maximum allowed length of 100 characters. This description should be considered invalid.",
				isValidPath:         true,
				isValidHash:         true,
				isValidShellVersion: true,
				isValidDescription:  false,
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
			if result.description != tt.expected.description {
				t.Errorf("LoadProfile().description = %v, expected %v", result.description, tt.expected.description)
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
			if result.isValidDescription != tt.expected.isValidDescription {
				t.Errorf("LoadProfile().isValidDescription = %v, expected %v", result.isValidDescription, tt.expected.isValidDescription)
			}
		})
	}
}
