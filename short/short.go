package short

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

var KeyToOrignal = make(map[string]string)
var OrignalToKey = make(map[string]string)
var DomainCounter = make(map[string]int)

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

func MakeShort(w http.ResponseWriter, r *http.Request){

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	inputURL := r.FormValue("url")

	parsedURL, err := url.ParseRequestURI(inputURL)
	
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	key, ok := OrignalToKey[inputURL]
	for !ok {
		key, _ = GenerateURLString(inputURL)
		_, ok = KeyToOrignal[key]
		ok = !ok
	}

	shortURL := "http://" + r.Host+ "/" +key

	DomainCounter[parsedURL.Hostname()]++
	OrignalToKey[inputURL] = key
	KeyToOrignal[key] = inputURL

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Shortened URL: %s\n", shortURL)
}
