//go:build !windows
// +build !windows

package launcher

var TempFiles []string

func CreateTempFileWithContent(content string, pattern string) (string, error) {
	return "", nil
}

func ExecutePowerShell(encodedCmd string) error {
	return nil
}

func ExecutePowerShellProcess(finalProfile string, shellPath string) error {
	return nil
}
