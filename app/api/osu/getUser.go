package osu

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type user struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

func GetUser(id string) (user, error) {
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

	var userData user

	if err := json.Unmarshal(body, &userData); err != nil {
		return userData, err
	}

	return userData, nil
}
