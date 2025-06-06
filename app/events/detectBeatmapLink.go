package events

import (
	"gengaozo/app/database"
	"gengaozo/app/handlers"
	"gengaozo/app/utils"

	"github.com/bwmarrin/discordgo"
)

func init() {
	handlers.RegisterEvent(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if len(m.Embeds) == 0 {
			return
		}

		if id := utils.ParseEmbed(m.Embeds[0]); id != "" {
			database.DB.Save(&database.Map{
				Channel_id: m.ChannelID,
				Map_id:     id,
			})
		}
	})
}
