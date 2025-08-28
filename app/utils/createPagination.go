package utils

import (
	"gengaozo/app/models"
	"gengaozo/app/store"
	"time"

	"github.com/bwmarrin/discordgo"
)

func CreatePagination(sess *discordgo.Session, pages []*discordgo.MessageEmbed, messageID, channelID string) {
	timer := time.AfterFunc(store.CleanupDelay, func() {
		store.CleanupPagination(sess, messageID, channelID)
	})

	store.PaginationMutex.Lock()
	store.Paginations[messageID] = &models.PaginationData{
		Pages:       pages,
		CurrentPage: 0,
		LastUsed:    time.Now(),
		Timer:       timer,
	}
	store.PaginationMutex.Unlock()
}
