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
		results, err := crawler.Search(duck_types.Search{Query: query})
		if err != nil {
			c.JSON(500, gin.H{"error": err})
		}
		// Put the snippet of each result into a slice
		snippets := make([]string, len(results))
		for i, result := range results {
			snippets[i] = result.Snippet
		}
		// Semantic search for each snippet
		semantic_results, err := vectordb.SemanticSearch([]string{query}, snippets, 3, false)
		if err != nil {
			c.JSON(500, gin.H{"error": err})
		}
		crawled_results := make([]string, len(semantic_results))
		// Crawl the URL of each result
		for i, result := range results {
			crawled_results[i], err = crawler.Crawl(result.Link)
			if err != nil {
				c.JSON(500, gin.H{"error": err})
			}
		}
		// Split the crawled results into sentences
		sentences := make([][]string, len(crawled_results))
		for i, result := range crawled_results {
			sentences[i] = splitSentences(result)
		}
		// Semantic search the sentences for top 10 results
		semantic_sentences := make([]string, 0)
		for _, sentence := range sentences {
			semantic_sentences = append(semantic_sentences, sentence...)
		}
		semantic_sentences_results, err := vectordb.SemanticSearch([]string{query}, semantic_sentences, 10, false)
		if err != nil {
			c.JSON(500, gin.H{"error": err})
		}
		// Return the top 10 results
		c.JSON(200, gin.H{"results": semantic_sentences_results})
	})

}

func splitSentences(text string) []string {
	sentences := make([]string, 0)
	punctuationMarks := map[rune]bool{
		'.': true,
		'!': true,
		'?': true,
	}

	// Replace newlines and multiple spaces with a single space
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "  ", " ")

	// Trim leading and trailing spaces
	text = strings.TrimSpace(text)

	// Split text into sentences based on punctuation marks
	sentenceStart := 0
	for i, char := range text {
		if punctuationMarks[char] {
			sentence := text[sentenceStart : i+1]
			sentences = append(sentences, strings.TrimSpace(sentence))
			sentenceStart = i + 1
		}
	}

	// Add the last sentence if there's any text remaining
	lastSentence := text[sentenceStart:]
	if len(lastSentence) > 0 {
		sentences = append(sentences, strings.TrimSpace(lastSentence))
	}

	// Group sentences into chunks of three joined by space
	chunks := make([]string, 0)
	chunkSize := 3
	numSentences := len(sentences)

	for i := 0; i < numSentences; i += chunkSize {
		end := i + chunkSize
		if end > numSentences {
			end = numSentences
		}

		chunk := strings.Join(sentences[i:end], " ")
		chunks = append(chunks, chunk)
	}

	return chunks
}
