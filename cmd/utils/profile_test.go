package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"testing"

	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/types"
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

type MockHashValidator struct{}

func (m MockHashValidator) ValidateHash(expectedHash, filePath string) (bool, error) {
	// Mock validation logic
	if expectedHash == "" {
		return false, fmt.Errorf("invalid hash")
	}
	return true, nil
}

func TestLoadProfile(t *testing.T) {
	// Create temporary files with known content
	_ = l.InitLogger("", "testslog.log", "debug")
	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		t.Fatalf("Failed to open os.DevNull: %v", err)
	}
	defer devNull.Close()
	l.Logger.SetOutput(devNull)
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
		expected types.ProfileItem
	}{
		{
			name: "Valid profile",
			line: []string{file1Path, hash1, "powershell", "A valid description"},
			expected: types.ProfileItem{
				Path:                file1Path,
				Hash:                hash1,
				Shell:               "powershell",
				ItemDescription:     "A valid description",
				IsValidPath:         true,
				IsValidShellVersion: true,
				IsValidDescription:  true,
			},
		},
		{
			name: "Invalid path",
			line: []string{"", hash1, "powershell", "A valid description"},
			expected: types.ProfileItem{
				Path:                "",
				Hash:                hash1,
				Shell:               "powershell",
				ItemDescription:     "A valid description",
				IsValidPath:         false,
				IsValidShellVersion: true,
				IsValidDescription:  true,
			},
		},
		{
			name: "Invalid shell version",
			line: []string{file2Path, hash2, "invalidShell", "A valid description"},
			expected: types.ProfileItem{
				Path:                file2Path,
				Hash:                hash2,
				Shell:               "invalidShell",
				ItemDescription:     "A valid description",
				IsValidPath:         true,
				IsValidShellVersion: false,
				IsValidDescription:  true,
			},
		},
		{
			name: "Invalid description",
			line: []string{file2Path, hash2, "powershell", "A very long description that exceeds the maximum allowed length of 100 characters. This description should be considered invalid."},
			expected: types.ProfileItem{
				Path:                file2Path,
				Hash:                hash2,
				Shell:               "powershell",
				ItemDescription:     "A very long description that exceeds the maximum allowed length of 100 characters. This description should be considered invalid.",
				IsValidPath:         true,
				IsValidShellVersion: true,
				IsValidDescription:  false,
			},
		},
		{
			name: "Empty profile",
			line: []string{"", "", "", ""},
			expected: types.ProfileItem{
				Path:                "",
				Hash:                "",
				Shell:               "",
				ItemDescription:     "",
				IsValidPath:         false,
				IsValidShellVersion: false,
				IsValidDescription:  true,
			},
		},
		{
			name: "Valid profile with different shell",
			line: []string{file1Path, hash1, "bash", "Another valid description"},
			expected: types.ProfileItem{
				Path:                file1Path,
				Hash:                hash1,
				Shell:               "bash",
				ItemDescription:     "Another valid description",
				IsValidPath:         true,
				IsValidShellVersion: false,
				IsValidDescription:  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LoadProfile(tt.line)
			if result.Path != tt.expected.Path {
				t.Errorf("LoadProfile().Path = %v, expected %v", result.Path, tt.expected.Path)
			}
			if result.Hash != tt.expected.Hash {
				t.Errorf("LoadProfile().Hash = %v, expected %v", result.Hash, tt.expected.Hash)
			}
			if result.Shell != tt.expected.Shell {
				t.Errorf("LoadProfile().Shell = %v, expected %v", result.Shell, tt.expected.Shell)
			}
			if result.ItemDescription != tt.expected.ItemDescription {
				t.Errorf("LoadProfile().ItemDescription = %v, expected %v", result.ItemDescription, tt.expected.ItemDescription)
			}
			if result.IsValidPath != tt.expected.IsValidPath {
				t.Errorf("LoadProfile().IsValidPath = %v, expected %v", result.IsValidPath, tt.expected.IsValidPath)
			}
			if result.IsValidShellVersion != tt.expected.IsValidShellVersion {
				t.Errorf("LoadProfile().IsValidShellVersion = %v, expected %v", result.IsValidShellVersion, tt.expected.IsValidShellVersion)
			}
			if result.IsValidDescription != tt.expected.IsValidDescription {
				t.Errorf("LoadProfile().IsValidDescription = %v, expected %v", result.IsValidDescription, tt.expected.IsValidDescription)
			}
		})
	}
}
