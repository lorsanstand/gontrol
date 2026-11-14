package os_utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type utils struct{}

// Вспомогательная функция для запуска команд и возврата их вывода
func (u *utils) runCmd(ctx context.Context, name string, arg ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, arg...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("command %s failed: %w (stderr: %s)", name, err, strings.TrimSpace(stderr.String()))
	}
	return stdout.Bytes(), nil
}

// 1. Выключение системы (с поддержкой задержки в минутах)
func (u *utils) Shutdown(ctx context.Context, args map[string]json.RawMessage) ([]byte, error) {
	var delay int
	if rawDelay, ok := args["delay"]; ok && len(rawDelay) > 0 {
		if err := json.Unmarshal(rawDelay, &delay); err != nil {
			return nil, fmt.Errorf("invalid delay: %w", err)
		}
	}

	timeArg := "now"
	if delay > 0 {
		timeArg = fmt.Sprintf("+%d", delay)
	}

	return u.runCmd(ctx, "shutdown", "-h", timeArg)
}

// 2. Перезагрузка системы
func (u *utils) Reboot(ctx context.Context, args map[string]json.RawMessage) ([]byte, error) {
	var delay int
	if rawDelay, ok := args["delay"]; ok && len(rawDelay) > 0 {
		if err := json.Unmarshal(rawDelay, &delay); err != nil {
			return nil, fmt.Errorf("invalid delay: %w", err)
		}
	}

	timeArg := "now"
	if delay > 0 {
		timeArg = fmt.Sprintf("+%d", delay)
	}

	return u.runCmd(ctx, "shutdown", "-r", timeArg)
}

// 3. Отмена запланированного выключения или перезагрузки
func (u *utils) CancelShutdown(ctx context.Context, args map[string]json.RawMessage) ([]byte, error) {
	return u.runCmd(ctx, "shutdown", "-c")
}

// 5. Завершение процесса по PID (с опциональным сигналом, например "9" или "15")
func (u *utils) KillProcess(ctx context.Context, args map[string]json.RawMessage) ([]byte, error) {
	var pid int
	if err := json.Unmarshal(args["pid"], &pid); err != nil || pid <= 0 {
		return nil, fmt.Errorf("valid pid is required: %w", err)
	}

	signal := "15" // По умолчанию SIGTERM (мягкое завершение)
	if rawSig, ok := args["signal"]; ok && len(rawSig) > 0 {
		var sig int
		if err := json.Unmarshal(rawSig, &sig); err != nil {
			return nil, fmt.Errorf("invalid signal: %w", err)
		}
		signal = strconv.Itoa(sig)
	}

	return u.runCmd(ctx, "kill", "-"+signal, strconv.Itoa(pid))
}

// 6. Информация о дисках (df -h)
func (u *utils) DiskUsage(ctx context.Context, args map[string]json.RawMessage) ([]byte, error) {
	return u.runCmd(ctx, "df", "-h")
}

// 7. Информация об оперативной памяти (free -m)
func (u *utils) MemoryUsage(ctx context.Context, args map[string]json.RawMessage) ([]byte, error) {
	return u.runCmd(ctx, "free", "-m")
}

// 8. Аптайм и средняя нагрузка на систему (uptime)
func (u *utils) SystemUptime(ctx context.Context, args map[string]json.RawMessage) ([]byte, error) {
	return u.runCmd(ctx, "uptime")
}

// 9. Сетевые интерфейсы и IP-адреса (ip -brief addr)
func (u *utils) NetworkInfo(ctx context.Context, args map[string]json.RawMessage) ([]byte, error) {
	return u.runCmd(ctx, "ip", "-brief", "addr")
}
