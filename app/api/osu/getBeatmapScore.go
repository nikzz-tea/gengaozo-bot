package osu

import (
	"encoding/json"
	"gengaozo/app/models"
	"io"
	"net/http"
)

func GetBeatmapScore(userID string, beatmapID string) (models.BeatmapScore, error) {
	var scoreData models.BeatmapScore

	token, err := getToken()
	if err != nil {
		return scoreData, err
	}

	req, _ := http.NewRequest("GET", baseURL+"/beatmaps/"+beatmapID+"/scores/users/"+userID, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return scoreData, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &scoreData); err != nil {
		return scoreData, err
	}

	return scoreData, nil
}
