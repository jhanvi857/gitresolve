package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func Init(verbose bool) {
	level := zerolog.WarnLevel
	if verbose {
		level = zerolog.InfoLevel
	}
	InitWithLevel(level)
}

func InitWithLevel(level zerolog.Level) {
	InitWithLevelAndOutput(level, os.Stderr)
}

func InitWithLevelAndOutput(level zerolog.Level, out io.Writer) {
	log = zerolog.New(out).
		With().
		Timestamp().
		Logger().
		Level(level)
}

func Info() *zerolog.Event  { return log.Info() }
func Debug() *zerolog.Event { return log.Debug() }
func Error() *zerolog.Event { return log.Error() }
func Warn() *zerolog.Event  { return log.Warn() }
func Trace() *zerolog.Event { return log.Trace() }

// Legacy helpers for existing string-only calls
func InfoMsg(msg string)  { log.Info().Msg(msg) }
func DebugMsg(msg string) { log.Debug().Msg(msg) }
func ErrorMsg(msg string) { log.Error().Msg(msg) }
