package commands

import (
	"gengaozo/app/handlers"
	"gengaozo/app/models"
	"strings"
)

func init() {
	handlers.RegisterCommand(models.CommandObject{
		Aliases: []string{"ping", "p"},
		Callback: func(props models.CommandProps) {
			sess, message, args := props.Sess, props.Message, props.Args

			sess.ChannelMessageSend(message.ChannelID, "pong "+strings.Join(args, " "))
		},
	})
}
