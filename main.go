package main

import (
	"SearchAPI/components/crawler"
	"strings"

	duck_types "github.com/acheong08/DuckDuckGo-API/typings"
	"github.com/acheong08/vectordb"

	gin "github.com/gin-gonic/gin"
)

func main() {
	handler := gin.Default()
	handler.GET("/search", func(c *gin.Context) {
		// Get query from request
		query := c.Query("query")
		// Search for query
		duck_results, err := crawler.Search(duck_types.Search{Query: query})
		if err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}
		// Put the snippet of each result into a slice
		snippets := make([]string, len(duck_results))
		for i, result := range duck_results {
			snippets[i] = result.Snippet
		}
		// Semantic search for each snippet
		semantic_results, err := vectordb.SemanticSearch([]string{query}, snippets, 3, false)
		if err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}
		// return the 3 most similar snippets
		similar_results := make([]duck_types.Result, len(semantic_results[0]))
		for i, result := range semantic_results[0] {
			similar_results[i] = duck_results[result.CorpusID]
		}
		c.JSON(200, gin.H{"sources": similar_results})
		// Crawl top 2 results
		similar_texts := make([]string, len(similar_results))
		if len(similar_results) < 2 {
			c.JSON(500, gin.H{"error": "Not enough results", "results": similar_results})
			return
		}
		for i, result := range similar_results[:2] {
			similar_texts[i], err = crawler.Crawl(result.Link)
			if err != nil {
				c.JSON(500, gin.H{"error": err})
				return
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
			c.JSON(500, gin.H{"error": err})
			return
		}
		// return the 3 most similar lines
		test_results := make([]string, len(semantic_results[0]))
		for i, result := range semantic_results[0] {
			test_results[i] = similar_texts_split[result.CorpusID]
		}
		c.JSON(200, gin.H{"results": test_results})
	})

	handler.Run()

}
