package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger
var File *os.File

func init() {
	File, err := os.OpenFile(
		"myapp.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	File.WriteString("\n")
	if err != nil {
		panic(err)
	}

	Logger = zerolog.New(File).With().Timestamp().Logger()
	// Logger.Info().Msg("Logger Initialized")

}
