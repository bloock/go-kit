package log

import (
	"github.com/rs/zerolog/pkgerrors"
	"os"

	"github.com/rs/zerolog"
)

func InitLogger(app string, debug bool) zerolog.Logger {
	var l zerolog.Logger
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	if debug {
		l = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Str("service", app).Logger()
	} else {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		l = zerolog.New(os.Stderr).With().Timestamp().Str("service", app).Logger()
	}

	return l
}
