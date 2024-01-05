package util

import (
	"github.com/rs/zerolog"
	"io"
	"os"
	"runtime/debug"
	"time"
)

func NewDevLogger(level zerolog.Level) zerolog.Logger {
	buildInfo, _ := debug.ReadBuildInfo()

	// colourise for consoles
	var output io.Writer = zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	if len(os.Getenv("STAGE")) != 0 {
		output = os.Stdout
	}

	return zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Caller().
		Int("pid", os.Getpid()).
		Str("go_version", buildInfo.GoVersion).
		Logger()
}
