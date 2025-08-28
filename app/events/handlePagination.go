package events

import (
	"gengaozo/app/handlers"
	"gengaozo/app/store"
	"gengaozo/app/utils"
	"time"

	"github.com/bwmarrin/discordgo"
)

func init() {
	handlers.RegisterEvent(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionMessageComponent {
			return
		}

		messageID := i.Message.ID

		store.PaginationMutex.Lock()
		pagination, exists := store.Paginations[messageID]
		store.PaginationMutex.Unlock()
		if !exists {
			return
		}

		pagination.Timer.Stop()
		pagination.Timer.Reset(store.CleanupDelay)
		pagination.LastUsed = time.Now()

		switch i.MessageComponentData().CustomID {
		case "page_first":
			pagination.CurrentPage = 0
		case "page_prev":
			pagination.CurrentPage--
		case "page_next":
			pagination.CurrentPage++
		case "page_last":
			pagination.CurrentPage = len(pagination.Pages) - 1
		default:
			return
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds:     []*discordgo.MessageEmbed{pagination.Pages[pagination.CurrentPage]},
				Components: utils.GetPaginationButtons(pagination.CurrentPage, len(pagination.Pages)),
			},
		})
	})
}
