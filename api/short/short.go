package short

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"urlshortner/database"
	"urlshortner/logger"

	"github.com/go-redis/redis"
)


func GenerateURLString(inputURL string) (string, error) {

	if len(inputURL) == 0 {
		return "", errors.New("recieved empty Url")
	}

	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz"
	charsetLen := len(charset)

	for _, char := range "/.:" {
		inputURL = strings.ReplaceAll(inputURL, string(char), string(charset[rand.Intn(charsetLen)]))
	}
	inputURL += strings.ToUpper(inputURL)

	if len(inputURL) >= 20 {
		charset = inputURL
	}

	const keyLength = 7
	charsetLen = len(charset)
	key := make([]byte, keyLength)
	for i := range key {
		key[i] = charset[rand.Intn(charsetLen)]
	}

	return string(key), nil
}

func MakeShort(w http.ResponseWriter, r *http.Request) {
	logger.Logs.Info().Msg("Entered in Short Fuction")

	rdb := database.CreateClient(0)
	defer rdb.Close()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		logger.Logs.Error().Msgf("Got wrong method for MakeShort request%s", r.Method)
		return
	}

	inputURL := r.FormValue("url")

	parsedURL, err := url.ParseRequestURI(inputURL)

	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		logger.Logs.Error().Msgf("Got invalid URL %s", err)
		return
	}

	key, err := rdb.HGet("OrignalToKey", inputURL).Result()
	if err == redis.Nil {
		key, _ = GenerateURLString(inputURL)
		// logger.Logs.Debug().Msgf("Required %v iterations for generating unique key", cnt)
		// _, err := rdb.HGet("OrignalToKey", inputURL).Result()
		// ok = !ok
	} else if err != nil {
		logger.Logs.Error().Msgf("Database connection failed %s", err)
		http.Error(w, "Database Connection Failed", http.StatusBadGateway)
		return
	}

	logger.Logs.Debug().Msgf("Recieved a key of type %T from GenerateURLString", key)

	shortURL := "http://" + r.Host + "/" + key

	rdb.HIncrBy("DomainCounter", parsedURL.Hostname(), 1)
	rdb.HSet("OrignalToKey", inputURL, key)
	rdb.HSet("KeyToOrignal", key, inputURL)

	// DomainCounter[parsedURL.Hostname()]++
	// OrignalToKey[inputURL] = key
	// KeyToOrignal[key] = inputURL

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Shortened URL: %s\n", shortURL)
	logger.Logs.Info().Msg("Exited from Short Fuction")
}
