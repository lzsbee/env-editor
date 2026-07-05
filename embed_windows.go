//go:build windows

package main

import (
	_ "embed"

	"env-editor/internal/elevated"
)

//go:embed elevhelper.exe
var elevHelperExe []byte

func init() {
	elevated.SetHelperBytes(elevHelperExe)
}
