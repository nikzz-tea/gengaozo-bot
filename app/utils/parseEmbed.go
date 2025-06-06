package utils

import (
	"github.com/bwmarrin/discordgo"
)

func ParseEmbed(emb *discordgo.MessageEmbed) string {
	fields := []string{
		emb.Description, emb.URL,
	}

	if emb.Author != nil {
		fields = append(fields, emb.Author.URL)
	}

	for _, field := range fields {
		if id := ParseBeatmapID(field); id != "" {
			return id
		}
	}

	return ""
}
