package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func Init(verbose bool) {
	level := zerolog.InfoLevel
	if verbose {
		level = zerolog.DebugLevel
	}
	log = zerolog.New(os.Stderr).
		With().
		Timestamp().
		Logger().
		Level(level)
}

func Info(msg string)  { log.Info().Msg(msg) }
func Debug(msg string) { log.Debug().Msg(msg) }
func Error(msg string) { log.Error().Msg(msg) }

func Infof(msg string, fields map[string]any) {
	e := log.Info()
	for k, v := range fields {
		e = e.Interface(k, v)
	}
	e.Msg(msg)
}
