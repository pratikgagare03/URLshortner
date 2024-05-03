package main

import (
	"fmt"
	http "net/http"
	"os"
	"urlshortner/logger"
	metrics "urlshortner/metrics"
	redirection "urlshortner/redirection"
	short "urlshortner/short"

	"github.com/joho/godotenv"
)

func setupRoutes() {
	http.HandleFunc("/makeshort", short.MakeShort)
	http.HandleFunc("/{key}", redirection.Redirect)
	http.HandleFunc("/metrics", metrics.GetMetrics)
}

func startServer() {
	fmt.Printf("localhost started at port%v", os.Getenv("APP_PORT"))
	err := http.ListenAndServe(os.Getenv("APP_PORT"), nil)
	if err != nil {
		logger.Logs.Error().Err(err)
		return
	}
}

func main() {
	logger.Logs.Info().Msg("Started Main")
	err := godotenv.Load(".env")
	if err != nil {
		logger.Logs.Error().Err(err)
	}
	defer logger.File.Close()


	setupRoutes()
	startServer()

	logger.Logs.Info().Msg("Main Function over")
}
