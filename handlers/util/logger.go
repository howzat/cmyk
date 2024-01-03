package util

import (
	"github.com/rs/zerolog"
	"os"
	"runtime/debug"
	"time"
)

func NewDevLogger(level zerolog.Level) zerolog.Logger {
	buildInfo, _ := debug.ReadBuildInfo()

	// colourise for consoles
	return zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(level).
		With().
		Timestamp().
		Caller().
		Int("pid", os.Getpid()).
		Str("go_version", buildInfo.GoVersion).
		Logger()
}
