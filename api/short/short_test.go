package short

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"
	"urlshortner/database"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateURLString(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}
	os.Setenv("DBNO", "0")

	dbno, _ := strconv.Atoi(os.Getenv("DBNO"))
	rdb := database.CreateClient(dbno)
	defer rdb.FlushAll()

	t.Run("TestEmptyInput", func(t *testing.T) {
		input := ""
		_, err := GenerateURLString(input)

		require.NotNil(t, err, "Expected error for empty input")
	})

	t.Run("validInput", func(t *testing.T) {
		input := "https://www.youtube.com/"
		output, err := GenerateURLString(input)

		assert.Equal(t, 7, len(output), "Output should match expected")
		require.Nil(t, err, "Unexpected error for valid input")
	})

	t.Run("UniqueKeyTest", func(t *testing.T) {
		input := "https://www.reddit.com/"
		output, _ := GenerateURLString(input)
		_, err := rdb.HGet("OrignalToKey", output).Result()

		if err != redis.Nil && err != nil {
			fmt.Println("Database connection failed", err)
			return
		}

		assert.Equal(t, redis.Nil, err, "Output should match expected")
	})

}

func TestMakeShort(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(MakeShort))

	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}
	os.Setenv("DBNO", "0")

	dbno, _ := strconv.Atoi(os.Getenv("DBNO"))
	rdb := database.CreateClient(dbno)
	defer rdb.FlushDB()

	t.Run("WrongMethodGET", func(t *testing.T) {
		resp, err := http.Get(server.URL)

		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, resp.StatusCode, http.StatusMethodNotAllowed)
	})

	t.Run("InvalidUrl", func(t *testing.T) {
		data := url.Values{}
		data.Set("url", "ht/w.youtube.com/watch?v=MDy7JQN5MN")

		resp, err := http.PostForm(server.URL, data)

		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Same Domain", func(t *testing.T) {
		data := url.Values{}
		data.Set("url", server.URL)
		fmt.Println(server.URL)
		resp, err := http.PostForm(server.URL, data)

		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("RepeatRequestSameURL", func(t *testing.T) {
		data := url.Values{}
		data.Set("url", "https://www.youtube.com/watch?v=MDy7JQN5MN")
		resp1, err := http.PostForm(server.URL, data)
		if err != nil {
			t.Error(err)
		}

		resp2, err := http.PostForm(server.URL, data)
		if err != nil {
			t.Error(err)
		}
		resp1body, _ := io.ReadAll(resp1.Body)
		resp2body, _ := io.ReadAll(resp2.Body)
		assert.Equal(t, resp1body, resp2body)
	})
	
	t.Run("Database Connection Failed", func(t *testing.T) {
		os.Setenv("DB_ADDR", "3000")
		data := url.Values{}
		data.Set("url", "https://www.youtube.com/watch?v=MDy7JQN5MN")
		resp1, err := http.PostForm(server.URL, data)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, http.StatusBadGateway, resp1.StatusCode)
	})

}
