package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
)

func shutdown(ctx context.Context, args map[string]json.RawMessage) ([]byte, error) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// /s — выключение, /t 0 — немедленно (0 секунд задержки)
		cmd = exec.Command("shutdown", "/s", "/t", "0")
	case "linux":
		// now — немедленно (требуются права root/sudo)
		cmd = exec.Command("shutdown", "-h", "now")
	case "darwin": // macOS
		cmd = exec.Command("shutdown", "-h", "now")
	default:
		return []byte{}, fmt.Errorf("unsupported: %s", runtime.GOOS)
	}

	return []byte{}, cmd.Run()
}
