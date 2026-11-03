package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)


func Init(level slog.Level) *slog.Logger {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		logger := slog.New(tint.NewTextHandler(os.Stdout, &tint.Options{
			AddSource:   true,
			Level:       level,
			ReplaceAttr: nil,
			TimeFormat:  "15:04:05",
			NoColor:     false,
		}))
		return logger
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	return logger
}
