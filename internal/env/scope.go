package env

// Scope identifies user-level or system-level environment variables.
type Scope int

const (
	ScopeUser Scope = iota
	ScopeSystem
)
