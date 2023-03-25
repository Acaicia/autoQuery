package credentials

import (
	"encoding/json"
	"os"
	// "path/filepath"
	// "autoQuery/common"
)

type Credentials struct {
	RedditApiKey     string `json:"reddit_api_key"`
	YouTubeAPIKey    string `json:"youtube_api_key"`
	TwitchAPIKey     string `json:"twitch_api_key"`
}

func LoadCredentials() (Credentials, error) {
	
	//file, err := os.Open(filepath.Join(common.GetProjectBaseFilepath(), "credentials.json"))
	file, err := os.Open(`/home/connor/go/src/autoQuery/credentials/credentials.json`)

	if err != nil {
		return Credentials{}, err
	}
	defer file.Close()

	var credentials Credentials
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&credentials)
	if err != nil {
		return Credentials{}, err
	}

	return credentials, nil
}
