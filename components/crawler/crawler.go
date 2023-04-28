package crawler

import (
	"io"
	"net/http"

	"github.com/acheong08/DuckDuckGo-API/duckduckgo"
	"github.com/acheong08/DuckDuckGo-API/typings"
)

func Search(search typings.Search) ([]typings.Result, error) {
	return duckduckgo.Get_results(search)
}

func Crawl(url string) (string, error) {
	// Construct get request
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return string(body), err
}
