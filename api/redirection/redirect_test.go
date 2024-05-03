package redirection

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"urlshortner/database"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestRedirect(t *testing.T) {
	//test db
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}
	os.Setenv("DBNO", "0")
	rdb := database.CreateClient(0)
	defer rdb.FlushDB()

	//test go
	mux := http.NewServeMux()

	mux.HandleFunc("/{key}", Redirect)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	t.Run("WrongMethodPOST", func(t *testing.T) {
		data := url.Values{}
		resp, err := http.PostForm(ts.URL+"/abc", data)

		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("EmptyUrl", func(t *testing.T) {
		resp, err := http.Get(ts.URL)

		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
	t.Run("WrongUrl", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/12vdsqA")

		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("validUrl", func(t *testing.T) {
		inputURL := "https://www.youtube.com/watch?v=MDy7JQN5MN4"
		rdb.HSet("KeyToOrignal", "qb5tyqA", inputURL)
		shortURL := ts.URL + "/" + "qb5tyqA"
		resp2, _ := http.Get(shortURL)
		assert.Equal(t, http.StatusFound, resp2.Request.Response.StatusCode)
	})

	t.Run("Database Connection Failed", func(t *testing.T) {
		os.Setenv("DB_ADDR", "3000")

		inputURL := "https://www.youtube.com/watch?v=MDy7JQN5MN4"
		rdb.HSet("KeyToOrignal", "qb5tyqA", inputURL)
		shortURL := ts.URL + "/" + "qb5tyqA"
		resp1, _ := http.Get(shortURL)
		assert.Equal(t, http.StatusBadGateway, resp1.StatusCode)
	})

}
