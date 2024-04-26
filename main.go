package main

import (
	"fmt"
	http "net/http"
	"urlshortner/logger"
	metrics "urlshortner/metrics"
	redirection "urlshortner/redirection"
	short "urlshortner/short"
)

func main() {
	defer logger.File.Close()
	logger.Logger.Info().Msg("Started Main")
	http.HandleFunc("/makeshort", short.MakeShort)
	http.HandleFunc("/", redirection.Redirect)
	http.HandleFunc("/metrics", metrics.GetMetrics)

	fmt.Printf("localhost started at port:8080")
	err := http.ListenAndServe(":8080", nil)
	if(err!=nil){
		logger.Logger.Error().Err(err)
		return
	}
	logger.Logger.Info().Msg("Main Function over")
}
