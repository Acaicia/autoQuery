package app

import (
	"net/http"
	"path/filepath"
	"html/template"
	"fmt"
	"github.com/gin-gonic/gin"
	"autoQuery/common"
	cred "autoQuery/credentials"
	api "autoQuery/api"
)

type searchParams struct {
	Query string `form:"query" json:"query" binding:"required"`
}

func ServeWeb(credentials cred.Credentials) {
    router := gin.Default()

    tmpl := template.Must(template.ParseGlob(filepath.Join("app", "index.tmpl")))
    router.SetHTMLTemplate(tmpl)

    router.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.tmpl", nil)
    })

    router.POST("/search", func(c *gin.Context) {
        var searchParams searchParams
        err := c.BindJSON(&searchParams)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
    
        query := searchParams.Query
        fmt.Println("Received search query:", query)
    
        results := search(query, credentials)
    
        fmt.Println("Search results count:", len(results))
        for _, result := range results {
            fmt.Printf("Result: %+v\n", result)
        }
    
        c.JSON(http.StatusOK, results)
    })

    router.Run(":8080")
}

func search(query string, credentials cred.Credentials) (results []common.SearchResult) { 
	redditResults, err := api.SearchReddit(query, credentials)
	if err == nil {
		results = append(results, redditResults...)
	}
	results = append(results, api.SearchYouTube(query, credentials.YouTubeAPIKey)...)
	results = append(results, api.SearchTwitch(query, credentials)...)
	return
}
