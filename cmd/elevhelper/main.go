//go:build windows

package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"golang.org/x/sys/windows/registry"
)

const regPath = `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`

func main() {
	if len(os.Args) < 2 {
		os.Exit(2)
	}

	var resultFile string
	var err error

	switch os.Args[1] {
	case "--elevated-set-system":
		if len(os.Args) != 5 {
			os.Exit(2)
		}
		resultFile = os.Args[4]
		name, decErr := base64.StdEncoding.DecodeString(os.Args[2])
		if decErr != nil {
			err = fmt.Errorf("decode name: %w", decErr)
			break
		}
		value, decErr := base64.StdEncoding.DecodeString(os.Args[3])
		if decErr != nil {
			err = fmt.Errorf("decode value: %w", decErr)
			break
		}
		err = setSystemEnvVar(string(name), string(value))
	case "--elevated-delete-system":
		if len(os.Args) != 4 {
			os.Exit(2)
		}
		resultFile = os.Args[3]
		name, decErr := base64.StdEncoding.DecodeString(os.Args[2])
		if decErr != nil {
			err = fmt.Errorf("decode name: %w", decErr)
			break
		}
		err = deleteSystemEnvVar(string(name))
	default:
		os.Exit(2)
	}

	writeResult(resultFile, err)
	if err != nil {
		os.Exit(1)
	}
}

func setSystemEnvVar(name, value string) error {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, regPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("open registry: %w", err)
	}
	defer key.Close()

	if err := key.SetStringValue(name, value); err != nil {
		return fmt.Errorf("set value: %w", err)
	}
	return nil
}

func deleteSystemEnvVar(name string) error {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, regPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("open registry: %w", err)
	}
	defer key.Close()

	if err := key.DeleteValue(name); err != nil {
		return fmt.Errorf("delete value: %w", err)
	}
	return nil
}

func writeResult(resultFile string, err error) {
	content := "OK"
	if err != nil {
		content = err.Error()
	}
	_ = os.WriteFile(resultFile, []byte(content), 0644)
}
