package model

// EnvVar represents a single environment variable entry.
type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ProcessInfo represents a running process entry.
type ProcessInfo struct {
	PID     int32  `json:"pid"`
	Name    string `json:"name"`
	ExePath string `json:"exePath"`
	Ports   string `json:"ports"`
}
