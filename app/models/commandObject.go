package models

import "github.com/bwmarrin/discordgo"

type CommandObject struct {
	Aliases  []string
	Callback func(CommandProps)
}

type CommandProps struct {
	Sess    *discordgo.Session
	Message *discordgo.MessageCreate
	Args    []string
}
