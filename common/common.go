package common

import (
	// "encoding/json"
	"path/filepath"
	"os"
	"time"
)

func GetProjectBaseFilepath() string {
	return filepath.Join(os.Getenv("GOPATH"), "src", "autoQuery")
}

type SearchResult struct {
	Platform     string    `json:"platform"`
	Title        string    `json:"title"`
	URL          string    `json:"url"`
	ThumbnailURL string    `json:"thumbnail_url"`
	Likes        int64     `json:"likes,omitempty"`
	Views        int64     `json:"views"`
	Comments     int64     `json:"comments"`
	UploadDate   time.Time `json:"upload_date"`
}