package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var Logs zerolog.Logger
var File *os.File

func init() {
	File, err := os.OpenFile(
		"myapp.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		panic(err)
	}
	
	File.WriteString("\n")
	Logs = zerolog.New(File).With().Timestamp().Logger()
	Logs.Info().Msg("Logger Initialized")

}
