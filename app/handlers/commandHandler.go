package handlers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

const prefix = "!"

var commands = make(map[string]func(s *discordgo.Session, m *discordgo.MessageCreate))

func CommandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if !strings.HasPrefix(m.Content, prefix) {
		return
	}

	commandName := strings.Fields(m.Content[len(prefix):])[0]

	if callback, exists := commands[commandName]; exists {
		callback(s, m)
	}
}

func RegisterCommand(name string, command func(s *discordgo.Session, m *discordgo.MessageCreate)) {
	commands[name] = command
}
