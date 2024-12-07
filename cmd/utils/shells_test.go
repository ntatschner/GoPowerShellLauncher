package utils

import (
	"fmt"
	"os"
	"testing"

	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/types"
)

func mockValidatePath(path string) (bool, error) {
	if path == "C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe" {
		return true, nil
	}
	if path == "C:\\Program Files\\PowerShell\\7\\pwsh.exe" {
		return true, nil
	}
	return false, fmt.Errorf("invalid path")
}

func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestLoadShells(t *testing.T) {
	// Initialize the logger
	l.InitLogger("", "testslog.log", "debug")

	// Redirect logger output to a temporary file
	logFile, err := os.CreateTemp("", "log_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp log file: %v", err)
	}
	defer os.Remove(logFile.Name())
	defer logFile.Close()
	l.Logger.SetOutput(logFile)

	// Mock the validatePath function
	originalValidatePath := ValidatePath
	ValidatePath = mockValidatePath
	defer func() { ValidatePath = originalValidatePath }()

	tests := []struct {
		name        string
		expected    []types.ShellItem
		expectError bool
	}{
		{
			name: "Valid shells",
			expected: []types.ShellItem{
				{ItemTitle: "PowerShell", ItemDescription: "PowerShell", Name: "PowerShell", ShortName: []string{"powershell", "all"}, Path: "C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe"},
				{ItemTitle: "PowerShell Core", ItemDescription: "PowerShell Core", Name: "PowerShell Core", ShortName: []string{"pwsh", "all"}, Path: "C:\\Program Files\\PowerShell\\7\\pwsh.exe"},
			},
			expectError: false,
		},
		{
			name:        "No valid shells",
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shells, err := utils.LoadShells()
			if (err != nil) != tt.expectError {
				t.Errorf("LoadShells() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if len(shells) != len(tt.expected) {
				t.Errorf("LoadShells() = %v, expected %v", shells, tt.expected)
			}
			for i, shell := range shells {
				if shell.ItemTitle != tt.expected[i].ItemTitle ||
					shell.ItemDescription != tt.expected[i].ItemDescription ||
					shell.Name != tt.expected[i].Name ||
					!equalSlices(shell.ShortName, tt.expected[i].ShortName) ||
					shell.Path != tt.expected[i].Path {
					t.Errorf("LoadShells()[%d] = %v, expected %v", i, shell, tt.expected[i])
				}
			}
		})
	}
}
