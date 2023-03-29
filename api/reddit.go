package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
    c "autoQuery/common"
    cred "autoQuery/credentials"
)

type RedditAccessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

type RedditData struct {
	Data struct {
		Children []struct {
			Data struct {
				Title        string  `json:"title"`
				URL          string  `json:"url"`
				ThumbnailURL string  `json:"thumbnail"`
				Score        int     `json:"score"`
				NumComments  int     `json:"num_comments"`
				Created      float64 `json:"created_utc"`
				Author       string  `json:"author"`
				IsVideo      bool    `json:"is_video"`
			} `json:"data"`
		} `json:"children"`
		After string `json:"after"`
	} `json:"data"`
}

func fetchRedditData(apiURL string, headers map[string]string) ([]byte, error) {
	resp, err := fetchData(apiURL, headers)
	if err != nil {
		return nil, err
	}

	data, err := readResponseBody(resp)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func getRedditAccessToken(creds cred.Credentials) (string, error) {
	apiURL := "https://www.reddit.com/api/v1/access_token"
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(creds.RedditApiKey, "")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "VideoSearch/0.1 by YourUsername")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var accessToken RedditAccessToken
	err = json.NewDecoder(resp.Body).Decode(&accessToken)
	if err != nil {
		return "", err
	}

	return accessToken.AccessToken, nil
}

func SearchReddit(query string, page int, credentials cred.Credentials) ([]c.SearchResult, error) {
	accessToken, err := getRedditAccessToken(credentials)
	if err != nil {
		return nil, fmt.Errorf("failed to get Reddit access token: %v", err)
	}

	var results []c.SearchResult
	after := ""

	for {
		apiURL := fmt.Sprintf("https://oauth.reddit.com/search.json?q=%s&type=link&limit=100&after=%s", url.QueryEscape(query), after)
		headers := map[string]string{
			"Authorization": "Bearer " + accessToken,
			"User-Agent":    "VideoSearch/0.1 by YourUsername",
		}
		data, err := fetchRedditData(apiURL, headers)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch Reddit data: %v", err)
		}

		var redditData RedditData
		err = parseJSONResponse(data, &redditData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Reddit data: %v", err)
		}

		for _, post := range redditData.Data.Children {
			result := c.SearchResult{
				Platform:     "Reddit",
				Title:        post.Data.Title,
				URL:          post.Data.URL,
				ThumbnailURL: post.Data.ThumbnailURL,
				Views:        int64(post.Data.Score),
				Likes:        int64(post.Data.Score),
				Comments:     int64(post.Data.NumComments),
				UploadDate:   time.Unix(int64(post.Data.Created), 0),
				Uploader:     post.Data.Author,
			}
			results = append(results, result)
		}

		// Check if there are more results to fetch
		if redditData.Data.After != "" {
			after = redditData.Data.After
		} else {
			break
		}

		// Optional: add a delay to avoid hitting the rate limits
		time.Sleep(2 * time.Second)
	}

	return results, nil
}