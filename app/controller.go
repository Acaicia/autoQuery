package app

import (
	"net/http"
	"path/filepath"
	"html/template"
	"github.com/gin-gonic/gin"
	"autoQuery/common"
	cred "autoQuery/credentials"
	api "autoQuery/api"
)

type searchParams struct {
	Query string `form:"query" json:"query" binding:"required"`
	Page  int    `form:"page" json:"page"`
}

func ServeWeb(credentials cred.Credentials) {
    router := gin.Default()

    html := template.Must(template.ParseGlob(filepath.Join("app", "index.html")))
    router.SetHTMLTemplate(html)

    // Add a route to serve static files
    router.StaticFile("/main.js", filepath.Join("app", "main.js"))

    router.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", nil)
    })


    router.POST("/search", func(c *gin.Context) {
        var searchParams searchParams
        err := c.BindJSON(&searchParams)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
    
        query := searchParams.Query
        page := searchParams.Page
        if page < 1 {
            page = 1
        }
        results := search(query, page, credentials)
        c.JSON(http.StatusOK, results)
    })

    router.Run(":8080")
}

func search(query string, page int, credentials cred.Credentials) (results []common.SearchResult) { 
    redditResults, err := api.SearchReddit(query, page, credentials)
    if err == nil {
        results = append(results, redditResults...)
    }
	results = append(results, api.SearchYouTube(query, credentials.YouTubeAPIKey)...)
	results = append(results, api.SearchTwitch(query, credentials)...)
	return
}
