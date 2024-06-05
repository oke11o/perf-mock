package logger

import (
	"log/slog"

	"github.com/oke11o/wslog"
)

func New() *slog.Logger {
	return wslog.New(true, slog.LevelDebug)
}
