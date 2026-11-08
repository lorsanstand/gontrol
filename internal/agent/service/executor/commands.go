package executor

import (
	modelsAgent "github.com/lorsanstand/gontrol/internal/agent/models"
	"github.com/lorsanstand/gontrol/internal/models"
)

var commands modelsAgent.Commands = map[string]modelsAgent.Command{
	"shutdown": {
		Exec: shutdown,
		Spec: models.CommandSpec{
			Command: "shutdown",
			Label:   "Выключение устройства",
			Args: []models.ArgSpec{
				{
					Name:      "delay",
					Label:     "Задержка выключения в секундах",
					Type:      models.ArgInt,
					Required:  false,
					Min:       new(0),
					MaxLength: new(120),
				},
			},
		},
	},
}
