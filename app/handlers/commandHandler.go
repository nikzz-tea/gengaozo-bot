package handlers

import (
	"gengaozo/app/models"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const prefix = "!"

var commands = make(map[string]func(models.CommandProps))

func CommandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author == nil {
		return
	}
	if !strings.HasPrefix(m.Content, prefix) {
		return
	}
	if m.Author.ID == s.State.User.ID {
		return
	}

	args := strings.Split(m.Content[len(prefix):], " ")
	commandName := strings.ToLower(args[0])
	callback, exists := commands[commandName]
	if !exists {
		return
	}

	log.Printf("'%v' used '%v' command\n", m.Author.Username, commandName)

	callback(models.CommandProps{
		Args:    args[1:],
		Sess:    s,
		Message: m,
	})
}

func RegisterCommand(command models.CommandObject) {
	for _, alias := range command.Aliases {
		commands[alias] = command.Callback
	}
}
