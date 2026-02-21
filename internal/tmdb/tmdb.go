package tmdb

import (
	"io"
	"net/http"
	"time"
)

type cacheItem struct {
	Data      []byte
	ExpiresAt time.Time
}

var cache = make(map[string]cacheItem)

func Fetch(url string) ([]byte, error) {

	if item, found := cache[url]; found && time.Now().Before(item.ExpiresAt) {
		return item.Data, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	cache[url] = cacheItem{
		Data:      data,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	return data, nil
}