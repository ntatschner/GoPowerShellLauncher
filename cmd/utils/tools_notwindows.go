//go:build !windows
// +build !windows

package utils

func executeCommandWithPowershell(encodedCmd string) error {
	return nil
}
