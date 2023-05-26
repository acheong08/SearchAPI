package main

import (
	"SearchAPI/components/search"

	gin "github.com/gin-gonic/gin"
)

func main() {
	handler := gin.Default()
	handler.GET("/search", func(c *gin.Context) {
		results, err := search.GetSearch(c.Query("query"))
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, results)
	})

	handler.Run()

}
