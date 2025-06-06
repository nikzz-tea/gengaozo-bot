package osu

import (
	"encoding/json"
	"gengaozo/app/models"
	"io"
	"log"
	"net/http"
)

func GetBeatmap(id string) (models.Beatmap, error) {
	token, err := getToken()
	if err != nil {
		log.Fatal(err)
	}

	req, _ := http.NewRequest("GET", baseURL+"/beatmaps/"+id, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var beatmapData models.Beatmap

	if err := json.Unmarshal(body, &beatmapData); err != nil {
		return beatmapData, err
	}

	return beatmapData, nil
}
