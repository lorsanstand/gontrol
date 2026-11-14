package executor

import (
	modelsAgent "github.com/lorsanstand/gontrol/internal/agent/models"
	"github.com/lorsanstand/gontrol/internal/agent/service/executor/os_utils"
	"github.com/lorsanstand/gontrol/internal/models"
)

var utils = os_utils.NewOSUtils()

var commands modelsAgent.Commands = map[string]modelsAgent.Command{
	"shutdown": {
		Exec: utils.Shutdown,
		Spec: models.CommandSpec{
			Command: "shutdown",
			Label:   "Выключение устройства",
			Args: []models.ArgSpec{
				{
					Name:     "delay",
					Label:    "Задержка выключения в минутах",
					Type:     models.ArgInt,
					Required: false,
					Min:      new(0),
					Max:      new(60),
				},
			},
		},
	},

	"reboot": {
		Exec: utils.Reboot,
		Spec: models.CommandSpec{
			Command: "reboot",
			Label:   "Перезагрузка устройства",
			Args: []models.ArgSpec{
				{
					Name:     "delay",
					Label:    "Задержка перезагрузки в минутах",
					Type:     models.ArgInt,
					Required: false,
					Min:      new(0),
					Max:      new(60),
				},
			},
		},
	},

	"cancel_shutdown": {
		Exec: utils.CancelShutdown,
		Spec: models.CommandSpec{
			Command: "cancel_shutdown",
			Label:   "Отмена перезагрузки или выключения устройства",
		},
	},

	"kill_process": {
		Exec: utils.KillProcess,
		Spec: models.CommandSpec{
			Command: "kill_process",
			Label:   "Завершение процесса",
			Args: []models.ArgSpec{
				{
					Name:     "pid",
					Label:    "PID процесса",
					Type:     models.ArgInt,
					Required: true,
					Min:      new(1),
				},
				{
					Name:     "signal",
					Label:    "Сигнал завершения (15 — мягко, 9 — принудительно)",
					Type:     models.ArgInt,
					Required: false,
					Min:      new(1),
					Max:      new(31),
				},
			},
		},
	},

	"disk_usage": {
		Exec: utils.DiskUsage,
		Spec: models.CommandSpec{
			Command: "disk_usage",
			Label:   "Информация о дисках",
		},
	},

	"memory_usage": {
		Exec: utils.MemoryUsage,
		Spec: models.CommandSpec{
			Command: "memory_usage",
			Label:   "Информация об оперативной памяти",
		},
	},

	"system_uptime": {
		Exec: utils.SystemUptime,
		Spec: models.CommandSpec{
			Command: "system_uptime",
			Label:   "Время работы и нагрузка системы",
		},
	},

	"network_info": {
		Exec: utils.NetworkInfo,
		Spec: models.CommandSpec{
			Command: "network_info",
			Label:   "Сетевые интерфейсы и IP-адреса",
		},
	},
}
