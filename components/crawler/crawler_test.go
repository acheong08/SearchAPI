package crawler_test

import (
	"SearchAPI/components/crawler"
	"fmt"
	"testing"

	duck_types "github.com/acheong08/DuckDuckGo-API/typings"
)

func TestSearch(t *testing.T) {
	results, err := crawler.Search(duck_types.Search{Query: "test"})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(results)
	if len(results) == 0 {
		t.Error("No results")
	}
}

func TestCrawl(t *testing.T) {
	results, err := crawler.Crawl("https://github.com/ggerganov/llama.cpp")
	if err != nil {
		t.Error(err)
	}
	if results != "" {
		fmt.Println(results)
	} else {
		t.Error("No results")
	}
}
