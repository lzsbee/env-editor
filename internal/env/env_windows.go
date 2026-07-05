//go:build windows

package env

import (
	"fmt"
	"sort"
	"strings"

	"env-editor/internal/elevated"
	"env-editor/internal/model"
	"env-editor/internal/winutil"

	"golang.org/x/sys/windows/registry"
)

func registryPath(scope Scope) (registry.Key, string, uint32) {
	switch scope {
	case ScopeUser:
		return registry.CURRENT_USER, `Environment`, registry.READ | registry.WRITE
	case ScopeSystem:
		return registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.READ | registry.WRITE
	default:
		return 0, "", 0
	}
}

func openKey(scope Scope, access uint32) (registry.Key, error) {
	root, path, defaultAccess := registryPath(scope)
	if access == 0 {
		access = defaultAccess
	}
	return registry.OpenKey(root, path, access)
}

func List(scope Scope) ([]model.EnvVar, error) {
	key, err := openKey(scope, registry.READ)
	if err != nil {
		return nil, fmt.Errorf("open registry: %w", err)
	}
	defer key.Close()

	names, err := key.ReadValueNames(-1)
	if err != nil {
		return nil, fmt.Errorf("read value names: %w", err)
	}

	vars := make([]model.EnvVar, 0, len(names))
	for _, name := range names {
		value, _, err := key.GetStringValue(name)
		if err != nil {
			continue
		}
		vars = append(vars, model.EnvVar{Name: name, Value: value})
	}

	sort.Slice(vars, func(i, j int) bool {
		return strings.ToLower(vars[i].Name) < strings.ToLower(vars[j].Name)
	})

	return vars, nil
}

func validateName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("变量名不能为空")
	}
	if strings.ContainsAny(name, "=\x00") {
		return fmt.Errorf("变量名不能包含 = 或空字符")
	}
	return nil
}

func Set(scope Scope, name, value string) error {
	if err := validateName(name); err != nil {
		return err
	}

	if scope == ScopeSystem && !winutil.IsProcessElevated() {
		if err := elevated.RunSystemSet(name, value); err != nil {
			return err
		}
		return winutil.BroadcastEnvChange()
	}

	key, err := openKey(scope, registry.WRITE)
	if err != nil {
		if scope == ScopeSystem {
			return fmt.Errorf("无法修改系统变量: %w", err)
		}
		return fmt.Errorf("open registry for write: %w", err)
	}
	defer key.Close()

	if err := key.SetStringValue(name, value); err != nil {
		return fmt.Errorf("set value: %w", err)
	}

	return winutil.BroadcastEnvChange()
}

func Delete(scope Scope, name string) error {
	if err := validateName(name); err != nil {
		return err
	}

	if scope == ScopeSystem && !winutil.IsProcessElevated() {
		if err := elevated.RunSystemDelete(name); err != nil {
			return err
		}
		return winutil.BroadcastEnvChange()
	}

	key, err := openKey(scope, registry.WRITE)
	if err != nil {
		if scope == ScopeSystem {
			return fmt.Errorf("无法修改系统变量: %w", err)
		}
		return fmt.Errorf("open registry for write: %w", err)
	}
	defer key.Close()

	if err := key.DeleteValue(name); err != nil {
		return fmt.Errorf("delete value: %w", err)
	}

	return winutil.BroadcastEnvChange()
}
