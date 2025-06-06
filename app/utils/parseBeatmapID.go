package utils

import "regexp"

func ParseBeatmapID(str string) string {
	regex := regexp.MustCompile(`https://osu\.ppy\.sh/(?:b/(\d+)|beatmapsets/\d+/?#?osu/(\d+)|beatmaps/(\d+))`)

	matches := regex.FindStringSubmatch(str)
	if matches == nil {
		return ""
	}

	for i := 1; i <= 3; i++ {
		if matches[i] != "" {
			return matches[i]
		}
	}

	return ""
}
