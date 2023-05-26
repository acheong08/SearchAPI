package main

import (
	"SearchAPI/components/crawler"

	duck_types "github.com/acheong08/DuckDuckGo-API/typings"
	"github.com/acheong08/vectordb"

	gin "github.com/gin-gonic/gin"
)

func main() {
	handler := gin.Default()
	handler.GET("/search", func(c *gin.Context) {
		// Get query from request
		query := c.Query("query")
		println(query)
		// Search for query
		duck_results, err := crawler.Search(duck_types.Search{Query: query})
		if err != nil {
			c.JSON(500, gin.H{"error": err})
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
		}
		// return the 3 most similar snippets
		similar_snippets := make([]string, len(semantic_results[0]))
		for i, result := range semantic_results[0] {
			similar_snippets[i] = snippets[result.CorpusID]
		}
		c.JSON(200, gin.H{"results": similar_snippets})
	})

	handler.Run()

}
