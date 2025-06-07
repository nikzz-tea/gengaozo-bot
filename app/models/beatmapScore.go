package models

type BeatmapScore struct {
	Score *Score `json:"score"`
}

type Score struct {
	Accuracy float32  `json:"accuracy"`
	Score    int      `json:"score"`
	PP       float64  `json:"pp"`
	MaxCombo int      `json:"max_combo"`
	Rank     string   `json:"rank"`
	Date     string   `json:"created_at"`
	Mods     []string `json:"mods"`
	Hits     Hits     `json:"statistics"`
	User     User     `json:"user"`
}

type Hits struct {
	Count300  int `json:"count_300"`
	Count100  int `json:"count_100"`
	Count50   int `json:"count_50"`
	CountMiss int `json:"count_miss"`
}
