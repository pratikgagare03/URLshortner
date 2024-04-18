package short

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)



func TestGenerateURLString(t *testing.T) {
	t.Run("TestEmptyInput", func (t* testing.T)  {
		input := ""
		_, err := GenerateURLString(input)

		require.NotNil(t, err, "Expected error for empty input")
	})

	t.Run("validInput", func (t* testing.T)  {
		input := "https://www.youtube.com/"
		output, err := GenerateURLString(input)

		assert.Equal(t,7,len(output),"Output should match expected")
		require.Nil(t, err, "Unexpected error for valid input")
	})

	t.Run("UniqueKeyTest", func (t* testing.T)  {
		input := "https://www.reddit.com/"
		output, err := GenerateURLString(input)
		_, exists := KeyToOrignal[output]

		assert.Equal(t, false,exists,"Output should match expected")
		require.Nil(t, err, "Unexpected error for valid input")
	})

	
}

func TestMakeShort(t *testing.T) {
	t.Run("WrongMethodGET", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(MakeShort))
		resp, err := http.Get(server.URL)

		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, resp.StatusCode, http.StatusMethodNotAllowed)
	})

	t.Run("InvalidUrl", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(MakeShort))
		data := url.Values{}
		data.Set("url", "ht/w.youtube.com/watch?v=MDy7JQN5MN4")
		
		resp, err := http.PostForm(server.URL,data)

		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("RepeatRequestSameURL", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(MakeShort))
		data := url.Values{}
		data.Set("url", "https://www.youtube.com/watch?v=MDy7JQN5MN4")
		
		resp1, err := http.PostForm(server.URL,data)
		if err != nil {
			t.Error(err)
		}

		resp2, err := http.PostForm(server.URL,data)
		if err != nil {
			t.Error(err)
		}
		resp1body,_ := io.ReadAll(resp1.Body)
		resp2body,_ := io.ReadAll(resp2.Body)
		assert.Equal(t,resp1body , resp2body)
	})

}