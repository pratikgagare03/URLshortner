package redirection

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"urlshortner/short"

	"github.com/stretchr/testify/assert"
)

func TestRedirect(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(Redirect))
	t.Run("WrongMethodPOST", func(t *testing.T) {
		resp, err := http.Post(server.URL, "", nil)

		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("EmptyUrl", func(t *testing.T) {
		resp, err := http.Get(server.URL)

		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
	t.Run("WrongUrl", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/12vdsqA")

		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("validUrl", func(t *testing.T) {
		mux := http.NewServeMux()

		mux.HandleFunc("/makeshort", short.MakeShort)
		mux.HandleFunc("/{key}", Redirect)

		ts := httptest.NewServer(mux)
		defer ts.Close()

		inputURL := "https://www.youtube.com/watch?v=MDy7JQN5MN4"
		data := url.Values{}
		data.Set("url", inputURL)
		input := ts.URL + "/makeshort"
		resp1, _ := http.PostForm(input, data)

		databyte, _ := io.ReadAll(resp1.Body)
		shortURL := strings.TrimSuffix(strings.TrimPrefix(string(databyte),"Shortened URL: "),"\n")

		resp2, _ := http.Get(shortURL)
		
		assert.Equal(t,resp2.Request.URL.String(), inputURL)
	})
}



