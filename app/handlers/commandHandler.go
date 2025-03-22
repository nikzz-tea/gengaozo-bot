package handlers

import (
	"gengaozo/app/models"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const prefix = "!"

var commands = make(map[string]func(models.CommandProps))

func CommandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if !strings.HasPrefix(m.Content, prefix) {
		return
	}

	args := strings.Split(m.Content[len(prefix):], " ")
	commandName := strings.ToLower(args[0])
	callback, exists := commands[commandName]
	if !exists {
		return
	}

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
