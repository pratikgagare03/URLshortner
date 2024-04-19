package redirection

import (
	"net/http"
	shorten "urlshortner/short"
)

func Redirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	key := r.URL.Path[1:]
	originalURL, exists := shorten.KeyToOrignal[key]
	if !exists{
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, originalURL, http.StatusFound)
}
