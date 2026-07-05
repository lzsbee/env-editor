//go:build !windows

package process

import (
	"fmt"

	"env-editor/internal/model"
)

func List() ([]model.ProcessInfo, error) {
	return nil, fmt.Errorf("process list is only supported on Windows")
}

func EnvVars(pid int32) ([]model.EnvVar, error) {
	return nil, fmt.Errorf("process environment is only supported on Windows")
}

func Kill(pid int32) error {
	return fmt.Errorf("kill process is only supported on Windows")
}
