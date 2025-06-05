package osu

import (
	"time"
)

const baseURL = "https://osu.ppy.sh/api/v2"

type TokenData struct {
	accessToken  string
	expires      time.Time
	clientId     string
	clientSecret string
}

var tokenData TokenData

func Auth(clientId string, clientSecret string) error {
	err := requestToken(clientId, clientSecret, &tokenData)
	if err != nil {
		return err
	}

	return nil
}

func getToken() (string, error) {
	if time.Now().Before(tokenData.expires) {
		return tokenData.accessToken, nil
	}

	err := requestToken(tokenData.clientId, tokenData.clientSecret, &tokenData)
	if err != nil {
		return "", err
	}

	return tokenData.accessToken, nil
}
