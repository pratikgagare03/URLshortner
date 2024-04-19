package metrics

import (
	"fmt"
	"net/http"
	"sort"
	short "urlshortner/short"
)

func GetMetrics(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

}
