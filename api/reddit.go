package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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
			Data c.SearchResult `json:"data"`
		} `json:"children"`
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

func SearchReddit(query string, creds cred.Credentials) ([]c.SearchResult, error) {
	accessToken, err := getRedditAccessToken(creds)
	if err != nil {
		return nil, fmt.Errorf("failed to get Reddit access token: %v", err)
	}

	apiURL := fmt.Sprintf("https://oauth.reddit.com/search.json?q=%s&type=video&limit=10", url.QueryEscape(query))
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

	results := make([]c.SearchResult, 0, len(redditData.Data.Children))
	for _, child := range redditData.Data.Children {
		results = append(results, child.Data)
	}

	return results, nil
}


