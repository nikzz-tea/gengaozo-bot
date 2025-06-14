package commands

import (
	"fmt"
	"gengaozo/app/handlers"
	"gengaozo/app/models"
)

func init() {
	handlers.RegisterCommand(models.CommandObject{
		Name:    "ping",
		Aliases: []string{"p"},
		Callback: func(props models.CommandProps) {
			sess, message := props.Sess, props.Message

			latency := sess.HeartbeatLatency().Milliseconds()

			sess.ChannelMessageSend(message.ChannelID, fmt.Sprintf("pong (%vms)", latency))
		},
	})
}
