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

	// В Windows параметр /t принимает задержку в секундах
	seconds := 0
	if delay > 0 {
		seconds = delay * 60
	}

	return u.runCmd(ctx, "shutdown", "/s", "/t", strconv.Itoa(seconds))
}

// 2. Перезагрузка системы
func (u *utils) Reboot(ctx context.Context, args map[string]json.RawMessage) ([]byte, error) {
	var delay int
	if rawDelay, ok := args["delay"]; ok && len(rawDelay) > 0 {
		if err := json.Unmarshal(rawDelay, &delay); err != nil {
			return nil, fmt.Errorf("invalid delay: %w", err)
		}
	}

	seconds := 0
	if delay > 0 {
		seconds = delay * 60
	}

	return u.runCmd(ctx, "shutdown", "/r", "/t", strconv.Itoa(seconds))
}

// 3. Отмена запланированного выключения или перезагрузки
func (u *utils) CancelShutdown(ctx context.Context, args map[string]json.RawMessage) ([]byte, error) {
	return u.runCmd(ctx, "shutdown", "/a")
}

// 5. Завершение процесса по PID (с опциональным сигналом: "9" принудительно, остальные — мягко)
func (u *utils) KillProcess(ctx context.Context, args map[string]json.RawMessage) ([]byte, error) {
	var pid int
	if err := json.Unmarshal(args["pid"], &pid); err != nil || pid <= 0 {
		return nil, fmt.Errorf("valid pid is required: %w", err)
	}

	signal := 15 // По умолчанию мягкое завершение
	if rawSig, ok := args["signal"]; ok && len(rawSig) > 0 {
		if err := json.Unmarshal(rawSig, &sig); err != nil {
			return nil, fmt.Errorf("invalid signal: %w", err)
		}
		signal = sig
	}

	cmdArgs := []string{"/PID", strconv.Itoa(pid)}
	// В Windows аналог SIGKILL (9) — флаг /F (принудительное завершение)
	if signal == 9 {
		cmdArgs = append(cmdArgs, "/F")
	}

	return u.runCmd(ctx, "taskkill", cmdArgs...)
}

// 6. Информация о дисках (аналог df -h через PowerShell)
func (u *utils) DiskUsage(ctx context.Context, args map[string]json.RawMessage) ([]byte, error) {
	psScript := `Get-PSDrive -PSProvider FileSystem | Select-Object Name, @{Name="Used(GB)";Expression={[math]::Round($_.Used/1GB,2)}}, @{Name="Free(GB)";Expression={[math]::Round($_.Free/1GB,2)}} | Format-Table -AutoSize`
	return u.runCmd(ctx, "powershell", "-NoProfile", "-Command", psScript)
}

// 7. Информация об оперативной памяти в МБ (аналог free -m через PowerShell)
func (u *utils) MemoryUsage(ctx context.Context, args map[string]json.RawMessage) ([]byte, error) {
	psScript := `$os = Get-CimInstance Win32_OperatingSystem; [PSCustomObject]@{ TotalMB = [math]::Round($os.TotalVisibleMemorySize / 1024); FreeMB = [math]::Round($os.FreePhysicalMemory / 1024); UsedMB = [math]::Round(($os.TotalVisibleMemorySize - $os.FreePhysicalMemory) / 1024) } | Format-Table -AutoSize`
	return u.runCmd(ctx, "powershell", "-NoProfile", "-Command", psScript)
}

// 8. Аптайм системы (аналог uptime через PowerShell)
func (u *utils) SystemUptime(ctx context.Context, args map[string]json.RawMessage) ([]byte, error) {
	psScript := `(Get-Date) - (Get-CimInstance Win32_OperatingSystem).LastBootUpTime | Select-Object Days, Hours, Minutes, Seconds | Format-Table -AutoSize`
	return u.runCmd(ctx, "powershell", "-NoProfile", "-Command", psScript)
}

// 9. Сетевые интерфейсы и IP-адреса (аналог ip -brief addr через PowerShell)
func (u *utils) NetworkInfo(ctx context.Context, args map[string]json.RawMessage) ([]byte, error) {
	psScript := `Get-NetIPAddress -AddressFamily IPv4 | Select-Object InterfaceAlias, IPAddress, PrefixLength | Format-Table -AutoSize`
	return u.runCmd(ctx, "powershell", "-NoProfile", "-Command", psScript)
}
