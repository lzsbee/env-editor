//go:build windows

package elevated

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"env-editor/internal/winutil"
)

var helperBytes []byte

func SetHelperBytes(data []byte) {
	helperBytes = data
}

func ensureElevHelper() (string, error) {
	dir, err := workDir()
	if err != nil {
		return "", err
	}

	helperPath := filepath.Join(dir, "elevhelper.exe")
	sum := sha256.Sum256(helperBytes)

	if data, err := os.ReadFile(helperPath); err == nil {
		if sha256.Sum256(data) == sum {
			return helperPath, nil
		}
	}

	if err := os.WriteFile(helperPath, helperBytes, 0755); err != nil {
		return "", fmt.Errorf("write elevhelper: %w", err)
	}
	return helperPath, nil
}

func launchHelper(args string) error {
	helperPath, err := ensureElevHelper()
	if err != nil {
		return err
	}
	return winutil.ShellExecuteRunAs(helperPath, args)
}

func workDir() (string, error) {
	dir := filepath.Join(os.Getenv("ProgramData"), "EnvEditor")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

func newResultFile() (string, error) {
	dir, err := workDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, fmt.Sprintf("result-%d.txt", time.Now().UnixNano())), nil
}

func RunSystemSet(name, value string) error {
	resultFile, err := newResultFile()
	if err != nil {
		return err
	}
	defer os.Remove(resultFile)

	nameB64 := base64.StdEncoding.EncodeToString([]byte(name))
	valueB64 := base64.StdEncoding.EncodeToString([]byte(value))
	args := fmt.Sprintf(`--elevated-set-system %s %s "%s"`, nameB64, valueB64, resultFile)

	if err := launchHelper(args); err != nil {
		return err
	}
	return winutil.WaitForResultFile(resultFile, 2*time.Minute)
}

func RunSystemDelete(name string) error {
	resultFile, err := newResultFile()
	if err != nil {
		return err
	}
	defer os.Remove(resultFile)

	nameB64 := base64.StdEncoding.EncodeToString([]byte(name))
	args := fmt.Sprintf(`--elevated-delete-system %s "%s"`, nameB64, resultFile)

	if err := launchHelper(args); err != nil {
		return err
	}
	return winutil.WaitForResultFile(resultFile, 2*time.Minute)
}
