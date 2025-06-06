package osu

import (
	"encoding/json"
	"gengaozo/app/models"
	"io"
	"log"
	"net/http"
)

func GetUser(id string) (models.User, error) {
	token, err := getToken()
	if err != nil {
		log.Fatal(err)
	}

	req, _ := http.NewRequest("GET", baseURL+"/users/"+id, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var userData models.User

	if err := json.Unmarshal(body, &userData); err != nil {
		return userData, err
	}

	return userData, nil
}
