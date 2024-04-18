package main

import (
	"fmt"
	http "net/http"
	metrics "urlshortner/metrics"
	redirection "urlshortner/redirection"
	short "urlshortner/short"
)

func main() {

	http.HandleFunc("/makeshort", short.MakeShort)
	http.HandleFunc("/", redirection.Redirect)
	http.HandleFunc("/metrics", metrics.GetMetrics)

	fmt.Printf("localhost started at port:8080")
	http.ListenAndServe(":8080", nil)
}
