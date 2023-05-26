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

func GetSearch(query string) (*SearchResults, error) {
	// Search for query
	duck_results, err := crawler.Search(duck_types.Search{Query: query})
	if err != nil {
		return nil, err
	}
	// Put the snippet of each result into a slice
	snippets := make([]string, len(duck_results))
	for i, result := range duck_results {
		snippets[i] = result.Snippet
	}
	// Semantic search for each snippet
	semantic_results, err := vectordb.SemanticSearch([]string{query}, snippets, 3, false)
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
	if len(sources) < 2 {
		return nil, fmt.Errorf("not enough results")
	}
	for i, result := range sources[:2] {
		similar_texts[i], err = crawler.Crawl(result.Link)
		if err != nil {
			return nil, err
		}
	}
	// Split texts by \n
	var similar_texts_split []string
	for _, text := range similar_texts {
		for _, line := range strings.Split(text, "\n") {
			// Get rid of short lines
			if len(line) > 70 {
				similar_texts_split = append(similar_texts_split, line)
			}
		}
	}
	// Semantic search for each line
	semantic_results, err = vectordb.SemanticSearch([]string{query}, similar_texts_split, 3, false)
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
