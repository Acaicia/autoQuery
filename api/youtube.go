package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
	c "autoQuery/common"
	// cred "autoQuery/credentials"
)

type YouTubeResponse struct {
	Items []struct {
		ID struct {
			VideoID string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			Title       string `json:"title"`
			PublishedAt string `json:"publishedAt"`
			Thumbnails  struct {
				Default struct {
					URL string `json:"url"`
				} `json:"default"`
			} `json:"thumbnails"`
			ChannelTitle string `json:"channelTitle"`
		} `json:"snippet"`
	} `json:"items"`
}

type YouTubeVideoResponse struct {
	Items []struct {
		Statistics struct {
			ViewCount    string `json:"viewCount"`
			LikeCount    string `json:"likeCount"`
			CommentCount string `json:"commentCount"`
		} `json:"statistics"`
	} `json:"items"`
}

func SearchYouTube(query string, apiKey string) []c.SearchResult {
	searchURL := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?part=snippet&q=%s&maxResults=1000&key=%s&type=video", url.QueryEscape(query), apiKey)

	resp, err := http.Get(searchURL)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var searchResponse YouTubeResponse
	err = json.Unmarshal(body, &searchResponse)
	if err != nil {
		return nil
	}

	results := make([]c.SearchResult, 0, len(searchResponse.Items))
	for _, item := range searchResponse.Items {
		videoResponse, err := getVideoDetails(item.ID.VideoID, apiKey)
		if err != nil {
			continue
		}

		likes := int64(0)
		views := int64(0)
		comments := int64(0)
		if len(videoResponse.Items) > 0 {
			likes = parseInt64(videoResponse.Items[0].Statistics.LikeCount)
			views = parseInt64(videoResponse.Items[0].Statistics.ViewCount)
			comments = parseInt64(videoResponse.Items[0].Statistics.CommentCount)
		}

		uploadDate, _ := time.Parse(time.RFC3339, item.Snippet.PublishedAt)

		results = append(results, c.SearchResult{
			Platform:     "YouTube",
			Title:        item.Snippet.Title,
			URL:          "https://www.youtube.com/watch?v=" + item.ID.VideoID,
			ThumbnailURL: item.Snippet.Thumbnails.Default.URL,
			Likes:        likes,
			Views:        views,
			Comments:     comments,
			UploadDate:   uploadDate,
			Uploader:     item.Snippet.ChannelTitle,
		})
	}

	return results
}

func getVideoDetails(videoID, apiKey string) (*YouTubeVideoResponse, error) {
	videoURL := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?id=%s&part=statistics&key=%s", videoID, apiKey)

	resp, err := http.Get(videoURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var videoResponse YouTubeVideoResponse
	err = json.Unmarshal(body, &videoResponse)
	if err != nil {
		return nil, err
	}

	return &videoResponse, nil
}