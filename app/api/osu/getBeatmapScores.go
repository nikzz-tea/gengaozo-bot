package osu

import (
	"encoding/json"
	"gengaozo/app/models"
	"io"
	"net/http"
)

func GetBeatmapScores(userID string, beatmapID string) (models.BeatmapScores, error) {
	var scoresData models.BeatmapScores

	token, err := getToken()
	if err != nil {
		return scoresData, err
	}

	req, _ := http.NewRequest(
		"GET", baseURL+"/beatmaps/"+beatmapID+"/scores/users/"+userID+"/all", nil,
	)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("x-api-version", "20220704")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return scoresData, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &scoresData); err != nil {
		return scoresData, err
	}

	return scoresData, nil
}
