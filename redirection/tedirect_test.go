package redirection

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestRedirect(t *testing.T) {
	t.Run("WrongMethodPOST", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(Redirect))
		resp, err := http.Post(server.URL, "", nil)

		if err != nil {
			t.Error(err)
		}
		
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})
}