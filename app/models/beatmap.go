package models

type Beatmap struct {
	StarRating float32    `json:"difficulty_rating"`
	Diffname   string     `json:"version"`
	MaxCombo   int        `json:"max_combo"`
	Beatmapset Beatmapset `json:"beatmapset"`
}

type Beatmapset struct {
	Artist string `json:"artist"`
	Title  string `json:"title"`
	Covers Covers `json:"covers"`
}

type Covers struct {
	List string `json:"list"`
}
