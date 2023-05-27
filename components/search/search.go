package search

import (
	"SearchAPI/components/crawler"
	"fmt"
	"strings"

	duck_types "github.com/acheong08/DuckDuckGo-API/typings"
	"github.com/acheong08/vectordb"
)

type SearchResults struct {
	Sources []duck_types.Result `json:"sources"`
	Results []string            `json:"results"`
}

func GetSearch(web_query, semantic_query string, source_limit, result_limit int) (*SearchResults, error) {
	// Search for query
	duck_results, err := crawler.Search(duck_types.Search{Query: web_query})
	if err != nil {
		return nil, err
	}
	// Put the snippet of each result into a slice
	snippets := make([]string, len(duck_results))
	for i, result := range duck_results {
		snippets[i] = result.Snippet
	}
	// Check if there are enough results
	if len(snippets) < source_limit {
		return nil, fmt.Errorf("not enough results")
	}
	// Semantic search for each snippet
	semantic_results, err := vectordb.SemanticSearch([]string{semantic_query}, snippets, source_limit, false)
	if err != nil {
		return nil, err
	}
	// return the 3 most similar snippets
	sources := make([]duck_types.Result, len(semantic_results[0]))
	for i, result := range semantic_results[0] {
		sources[i] = duck_results[result.CorpusID]
	}
	// c.JSON(200, gin.H{"sources": similar_results})
	// Crawl top 2 results
	similar_texts := make([]string, len(sources))
	if len(sources) < result_limit {
		return nil, fmt.Errorf("not enough results")
	}
	for i, result := range sources[:result_limit] {
		similar_texts[i], err = crawler.Crawl(result.Link)
		if err != nil {
			return nil, err
		}
	}
	// Split texts by \n
	similar_texts_split := split_text(strings.Join(similar_texts, "\n"))
	// Semantic search for each line
	semantic_results, err = vectordb.SemanticSearch([]string{semantic_query}, similar_texts_split, result_limit, false)
	if err != nil {
		return nil, err
	}
	// return the 3 most similar lines
	text_results := make([]string, len(semantic_results[0]))
	for i, result := range semantic_results[0] {
		text_results[i] = similar_texts_split[result.CorpusID]
	}
	return &SearchResults{Sources: sources, Results: text_results}, nil
}
func split_text(text string) []string {
	var chunks []string
	lines := strings.Split(text, "\n")

	var chunk string
	for _, line := range lines {
		if len(line) < 50 {
			continue
		}
		if len(chunk)+len(line) > 1000 {
			chunks = append(chunks, chunk)
			chunk = ""
		}
		if len(chunk)+len(line) > 1200 {
			chunks = append(chunks, chunk)
			chunk = ""
		}
		chunk += line + "\n"
	}
	return chunks
}
