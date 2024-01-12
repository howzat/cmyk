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
	var output io.Writer = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	return NewZeroLog(level, output).
		With().
		Int("pid", os.Getpid()).
		Str("go_version", buildInfo.GoVersion).
		Logger()
}

func NewProdLogger(level zerolog.Level) zerolog.Logger {
	buildInfo, _ := debug.ReadBuildInfo()
	return NewZeroLog(level, os.Stdout).
		With().
		Int("pid", os.Getpid()).
		Str("go_version", buildInfo.GoVersion).
		Logger()
}

func NewZeroLog(level zerolog.Level, out io.Writer) zerolog.Logger {
	return zerolog.New(out).
		Level(level).
		With().
		Timestamp().
		Caller().
		Logger()
}
