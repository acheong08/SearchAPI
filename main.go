package main

import (
	"SearchAPI/components/search"
	"fmt"
	"strconv"

	gin "github.com/gin-gonic/gin"
)

func parse_get(c *gin.Context) (string, string, int, int, error) {
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
		return "", "", 0, 0, fmt.Errorf("slimit must be an integer")
	}
	rlimit_int, err := strconv.Atoi(results_limit)
	if err != nil {
		return "", "", 0, 0, fmt.Errorf("rlimit must be an integer")
	}
	web_query := c.Query("wq")
	semantic_query := c.Query("sq")
	if web_query == "" {
		return "", "", 0, 0, fmt.Errorf("web query must be provided")
	}
	if semantic_query == "" {
		semantic_query = web_query
	}
	return web_query, semantic_query, slimit_int, rlimit_int, nil
}
func parse_post(c *gin.Context) (string, string, int, int, error) {
	type PostBody struct {
		WebQuery      string `json:"wq"`
		SemanticQuery string `json:"sq"`
		SourcesLimit  int    `json:"slimit"`
		ResultsLimit  int    `json:"rlimit"`
	}
	var post_body PostBody
	err := c.BindJSON(&post_body)
	if err != nil {
		return "", "", 0, 0, fmt.Errorf("invalid json")
	}
	if post_body.WebQuery == "" {
		return "", "", 0, 0, fmt.Errorf("web query must be provided")
	}
	if post_body.SemanticQuery == "" {
		post_body.SemanticQuery = post_body.WebQuery
	}
	if post_body.SourcesLimit == 0 {
		post_body.SourcesLimit = 3
	}
	if post_body.ResultsLimit == 0 {
		post_body.ResultsLimit = 2
	}
	return post_body.WebQuery, post_body.SemanticQuery, post_body.SourcesLimit, post_body.ResultsLimit, nil
}

func main() {
	handler := gin.Default()
	handler.GET("/search", search_handler)
	handler.POST("/search", search_handler)

	handler.Run()

}

func search_handler(c *gin.Context) {
	var web_query, semantic_query string
	var slimit_int, rlimit_int int
	var err error
	if c.Request.Method == "GET" {
		web_query, semantic_query, slimit_int, rlimit_int, err = parse_get(c)
	} else {
		web_query, semantic_query, slimit_int, rlimit_int, err = parse_post(c)
	}
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	results, err := search.GetSearch(web_query, semantic_query, slimit_int, rlimit_int)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, results)
}
