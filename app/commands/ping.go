package commands

import (
	"gengaozo/app/handlers"

	"github.com/bwmarrin/discordgo"
)

func init() {
	handlers.RegisterCommand("ping", func(s *discordgo.Session, m *discordgo.MessageCreate) {
		s.ChannelMessageSend(m.ChannelID, "pong "+m.Author.Username)
	})
}
