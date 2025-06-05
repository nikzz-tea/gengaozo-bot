package osu

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

const tokenURL = "https://osu.ppy.sh/oauth/token"

func requestToken(clientId string, clientSecret string, tokenData *TokenData) error {
	form := url.Values{}
	form.Set("client_id", clientId)
	form.Set("client_secret", clientSecret)
	form.Set("grant_type", "client_credentials")
	form.Set("scope", "public")

	resp, err := http.PostForm(tokenURL, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var data struct {
		Token   string `json:"access_token"`
		Expires int    `json:"expires_in"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	*tokenData = TokenData{
		accessToken:  data.Token,
		expires:      time.Now().Add(time.Duration(data.Expires-60) * time.Second),
		clientId:     clientId,
		clientSecret: clientSecret,
	}

	return nil
}
