package osu

import (
	"encoding/json"
	"gengaozo/app/models"
	"io"
	"net/http"
)

func GetBeatmap(id string) (models.Beatmap, error) {
	var beatmapData models.Beatmap

	token, err := getToken()
	if err != nil {
		return beatmapData, err
	}

	req, _ := http.NewRequest("GET", baseURL+"/beatmaps/"+id, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return beatmapData, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &beatmapData); err != nil {
		return beatmapData, err
	}

	return beatmapData, nil
}
