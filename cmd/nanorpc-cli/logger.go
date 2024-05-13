package main

import (
	"os"

	"darvaza.org/sidecar/pkg/logger/zerolog"
	"darvaza.org/slog"
)

func newLogger(level slog.LogLevel) slog.Logger {
	return zerolog.New(os.Stderr, level)
}
