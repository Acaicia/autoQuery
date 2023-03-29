package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
	c "autoQuery/common"
	credentials "autoQuery/credentials"
)

type twitchResponse struct {
	Data []struct {
		DisplayName      string `json:"display_name"`
		URL              string `json:"url"`
		BroadcasterLogin string `json:"broadcaster_login"`
	} `json:"data"`
}


func SearchTwitch(query string, credentials credentials.Credentials) []c.SearchResult {
	twitchAPIKey := credentials.TwitchAPIKey
	twitchURL := "https://api.twitch.tv/helix/search/channels"

	req, err := http.NewRequest("GET", twitchURL+"?query="+query+"&first=10", nil)
	if err != nil {
		return nil
	}

	req.Header.Set("Client-ID", twitchAPIKey)
	req.Header.Set("Authorization", "Bearer "+twitchAPIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var response twitchResponse
	json.Unmarshal(body, &response)

	var results []c.SearchResult
	for _, item := range response.Data {
		results = append(results, c.SearchResult{
			Platform:     "Twitch",
			Title:        item.DisplayName,
			URL:          item.URL,
			ThumbnailURL: "", // Add an appropriate thumbnail URL
			Likes:        0,   // Add an appropriate likes count
			Views:        0,   // Add an appropriate views count
			Comments:     0,   // Add an appropriate comments count
			UploadDate:   time.Time{}, // Add an appropriate upload date
			Uploader:     item.BroadcasterLogin,
		})
	}

	return results
}