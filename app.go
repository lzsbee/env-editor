package main

import (
	"context"

	"env-editor/internal/env"
	"env-editor/internal/model"
	"env-editor/internal/process"
)

// Keep type aliases in main so Wails bindings stay under main.*.
type EnvVar = model.EnvVar
type ProcessInfo = model.ProcessInfo

// App exposes environment variable operations to the frontend.
type App struct{}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	_ = ctx
}

func (a *App) ListUserVars() ([]EnvVar, error) {
	return env.List(env.ScopeUser)
}

func (a *App) ListSystemVars() ([]EnvVar, error) {
	return env.List(env.ScopeSystem)
}

func (a *App) SetUserVar(name, value string) error {
	return env.Set(env.ScopeUser, name, value)
}

func (a *App) SetSystemVar(name, value string) error {
	return env.Set(env.ScopeSystem, name, value)
}

func (a *App) DeleteUserVar(name string) error {
	return env.Delete(env.ScopeUser, name)
}

func (a *App) DeleteSystemVar(name string) error {
	return env.Delete(env.ScopeSystem, name)
}

func (a *App) ListProcesses() ([]ProcessInfo, error) {
	return process.List()
}

func (a *App) GetProcessEnvVars(pid int32) ([]EnvVar, error) {
	return process.EnvVars(pid)
}

func (a *App) KillProcess(pid int32) error {
	return process.Kill(pid)
}
