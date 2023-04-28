package crawler

import (
	"io"
	"net/http"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/acheong08/DuckDuckGo-API/duckduckgo"
	"github.com/acheong08/DuckDuckGo-API/typings"
)

var converter = md.NewConverter("", true, nil)

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
	if err != nil {
		return "", err
	}
	markdown, err := converter.ConvertString(string(body))
	return markdown, err
}
