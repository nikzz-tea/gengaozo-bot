package commands

import (
	"fmt"
	"gengaozo/app/api/osu"
	"gengaozo/app/database"
	"gengaozo/app/handlers"
	"gengaozo/app/models"
	"log"
	"strconv"
	"strings"
)

func init() {
	handlers.RegisterCommand(models.CommandObject{
		Aliases: []string{"set"},
		Callback: func(props models.CommandProps) {
			sess, message, args := props.Sess, props.Message, props.Args
			if len(args) == 0 {
				sess.ChannelMessageSend(message.ChannelID, "🔴 **Provide a user id or name**")
				return
			}

			user, err := osu.GetUser(strings.Join(args, " "))
			if err != nil {
				log.Fatal(err)
				return
			}

			if user.Username == "" {
				sess.ChannelMessageSend(message.ChannelID, "🔴 **No user was found**")
				return
			}

			database.DB.Save(database.User{
				Discord_id: message.Author.ID,
				Osu_id:     strconv.Itoa(user.Id),
			})

			sess.ChannelMessageSend(message.ChannelID, fmt.Sprintf("`%v`'s osu! username is set to `%v`", message.Message.Author.Username, user.Username))
		},
	})
}
