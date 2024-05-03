package metrics

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"urlshortner/database"

	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestGetMetrics(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}
	os.Setenv("DBNO", "0")

	dbno, _ := strconv.Atoi(os.Getenv("DBNO"))
	rdb := database.CreateClient(dbno)
	defer rdb.FlushDB()

	t.Run("EmptyMap", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(GetMetrics))
		rdb.FlushAll()
		resp, err := http.Get(server.URL)

		if err != nil {
			t.Error(err)
		}

		respBody, err := io.ReadAll(resp.Body)

		if err != nil {
			t.Error(err)
		}
		expected := "No domain shortened"

		assert.Equal(t, expected, string(respBody))

	})

	t.Run("WrongMethodPOST", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(GetMetrics))
		resp, err := http.Post(server.URL, "", nil)

		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("CheckSort", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(GetMetrics))
		rdb.HIncrBy("DomainCounter", "youtube.com", 5)
		rdb.HIncrBy("DomainCounter", "gmail.com", 2)
		rdb.HIncrBy("DomainCounter", "whatsapp.com", 8)
		rdb.HIncrBy("DomainCounter", "reddit.com", 3)

		resp, err := http.Get(server.URL)
		if err != nil {
			t.Error(err)
		}

		respBody, _ := io.ReadAll(resp.Body)
		expected := "Top 3 most shortened domains are :\nwhatsapp.com: 8\nyoutube.com: 5\nreddit.com: 3\n"
		assert.Equal(t, expected, string(respBody))
	})

}
