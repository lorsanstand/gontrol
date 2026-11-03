package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

func Init(level slog.Level, format string) *slog.Logger {
	switch format {
	case "json":
		return jsonLogger(level)
	case "console":
		return consoleLogger(level)
	default:
		if isatty.IsTerminal(os.Stdout.Fd()) {
			return consoleLogger(level)
		}
		return jsonLogger(level)
	}
}

func jsonLogger(l slog.Level) *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: l}))
}

func consoleLogger(l slog.Level) *slog.Logger {
	return slog.New(tint.NewTextHandler(os.Stdout, &tint.Options{
		AddSource:   true,
		Level:       l,
		ReplaceAttr: nil,
		TimeFormat:  "15:04:05",
		NoColor:     false,
	}))
}
