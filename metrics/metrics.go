package metrics

import (
	"fmt"
	"net/http"
	"sort"
	"urlshortner/logger"
	short "urlshortner/short"
)

func GetMetrics(w http.ResponseWriter, r *http.Request) {
	logger.Logger.Info().Msg("Entered in GetMetrics Fuction")
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		logger.Logger.Error().Msgf("Got wrong method for GetMetrics request %s", r.Method)
		return
	}

	mapLen := min(3, len(short.DomainCounter))
	if mapLen == 0 {
		fmt.Fprint(w, "No domain shortened")
		return
	} else {
		fmt.Fprintf(w, "Top %d most shortened domains are :\n", mapLen)
	}

	type kv struct {
		Key   string
		Value int
	}

	var ss []kv
	for k, v := range short.DomainCounter {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	for _, kv := range ss[:mapLen] {
		fmt.Fprintf(w, "%s: %d\n", kv.Key, kv.Value)
	}

	logger.Logger.Info().Msg("Exited from GetMetrics Fuction")
}
