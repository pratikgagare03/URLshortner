package redirection

import (
	"net/http"
	"urlshortner/database"
	"urlshortner/logger"

	"github.com/go-redis/redis"
)

func Redirect(w http.ResponseWriter, r *http.Request) {
	logger.Logs.Info().Msg("Entered in Redirect Fuction ")

	rdb := database.CreateClient(0)
	defer rdb.Close()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		logger.Logs.Error().Msgf("Got wrong method for Redirect request %s", r.Method)
		return
	}

	key := r.PathValue("key")
	originalURL, err := rdb.HGet("KeyToOrignal", key).Result()
	if err == redis.Nil {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	} else if err != nil {
		logger.Logs.Error().Msgf("Database connection failed %s", err)
		http.Error(w, "Database Connection Failed", http.StatusBadGateway)
		return
	}
	http.Redirect(w, r, originalURL, http.StatusFound)
	logger.Logs.Info().Msg("Exited from redirect Fuction")
}
