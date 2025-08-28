package store

import (
	"gengaozo/app/models"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

var Paginations = make(map[string]*models.PaginationData)
var PaginationMutex sync.Mutex
var CleanupDelay = 30 * time.Second

func CleanupPagination(s *discordgo.Session, messageID, channelID string) {
	PaginationMutex.Lock()
	defer PaginationMutex.Unlock()

	pagination, exists := Paginations[messageID]
	if !exists {
		return
	}

	s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel:    channelID,
		ID:         messageID,
		Embed:      pagination.Pages[pagination.CurrentPage],
		Components: &[]discordgo.MessageComponent{},
	})

	delete(Paginations, messageID)
}
