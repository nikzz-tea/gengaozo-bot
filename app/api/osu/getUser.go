package osu

import (
	"encoding/json"
	"gengaozo/app/models"
	"io"
	"net/http"
)

func GetUser(id string) (models.User, error) {
	var userData models.User

	token, err := getToken()
	if err != nil {
		return userData, err
	}

	req, _ := http.NewRequest("GET", baseURL+"/users/"+id, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return userData, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &userData); err != nil {
		return userData, err
	}

	return userData, nil
}
