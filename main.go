package main

import (
	// "encoding/json"
	"log"
	// "AutoQuery/common"
	// api "AutoQuery/api"
	cred "autoQuery/credentials"
	"autoQuery/app"
)

func main() {
	credentials, err := cred.LoadCredentials()
	if err != nil {
		log.Fatalf("Error loading credentials: %v", err)
	}

	app.ServeWeb(credentials)
}