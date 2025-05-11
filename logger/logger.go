package logger

import (
	"log/slog"
	"os"
)

func LogInit() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}
