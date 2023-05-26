package main

import (
	"SearchAPI/components/search"
	"strconv"

	gin "github.com/gin-gonic/gin"
)

func main() {
	handler := gin.Default()
	handler.GET("/search", func(c *gin.Context) {
		sources_limit := c.Query("slimit")
		results_limit := c.Query("rlimit")
		// Check if blank
		if sources_limit == "" {
			sources_limit = "3"
		}
		if results_limit == "" {
			results_limit = "2"
		}
		// Convert to int
		slimit_int, err := strconv.Atoi(sources_limit)
		if err != nil {
			c.JSON(400, gin.H{"error": "slimit must be an integer"})
			return
		}
		rlimit_int, err := strconv.Atoi(results_limit)
		if err != nil {
			c.JSON(400, gin.H{"error": "rlimit must be an integer"})
			return
		}
		results, err := search.GetSearch(c.Query("query"), slimit_int, rlimit_int)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, results)
	})

	handler.Run()

}
