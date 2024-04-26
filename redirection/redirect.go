package redirection

import (
	"net/http"
	"urlshortner/logger"
	shorten "urlshortner/short"
)

func Redirect(w http.ResponseWriter, r *http.Request) {
	logger.Logs.Info().Msg("Entered in Redirect Fuction ")
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		logger.Logs.Error().Msgf("Got wrong method for Redirect request %s", r.Method)
		return
	}

	key := r.URL.Path[1:]
	originalURL, exists := shorten.KeyToOrignal[key]
	if !exists {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, originalURL, http.StatusFound)
	logger.Logs.Info().Msg("Exited from redirect Fuction")
}
