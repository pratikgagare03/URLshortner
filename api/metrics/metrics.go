package metrics

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"urlshortner/database"
	"urlshortner/logger"
)

func GetMetrics(w http.ResponseWriter, r *http.Request) {
	logger.Logs.Info().Msg("Entered in GetMetrics Fuction")
	rdb := database.CreateClient(0)
	defer rdb.Close()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		logger.Logs.Error().Msgf("Got wrong method for GetMetrics request %s", r.Method)
		return
	}

	mapLen := min(3, rdb.HLen("DomainCounter").Val())
	if mapLen == 0 {
		fmt.Fprint(w, "No domain shortened")
		return
	} else {
		fmt.Fprintf(w, "Top %d most shortened domains are :\n", mapLen)
	}

	DomainCounter := rdb.HGetAll("DomainCounter").Val()
	
	type kv struct {
		Key   string
		Value int64
	}

	var ss []kv
	for k, v := range DomainCounter {
		vint, _ := strconv.ParseInt(v, 10, 0)
		ss = append(ss, kv{k, vint})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	for _, kv := range ss[:mapLen] {
		fmt.Fprintf(w, "%s: %d\n", kv.Key, kv.Value)
	}

	logger.Logs.Info().Msg("Exited from GetMetrics Fuction")
}
