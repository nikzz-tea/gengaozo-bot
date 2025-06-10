package models

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type PaginationData struct {
	Pages       []*discordgo.MessageEmbed
	CurrentPage int
	LastUsed    time.Time
	Timer       *time.Timer
}
