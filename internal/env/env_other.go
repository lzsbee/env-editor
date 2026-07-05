//go:build !windows

package env

import (
	"fmt"

	"env-editor/internal/model"
)

func List(scope Scope) ([]model.EnvVar, error) {
	return nil, fmt.Errorf("persistent environment editing is only supported on Windows")
}

func Set(scope Scope, name, value string) error {
	return fmt.Errorf("persistent environment editing is only supported on Windows")
}

func Delete(scope Scope, name string) error {
	return fmt.Errorf("persistent environment editing is only supported on Windows")
}
