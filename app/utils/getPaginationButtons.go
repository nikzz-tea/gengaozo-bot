package utils

import "github.com/bwmarrin/discordgo"

func GetPaginationButtons(i, pages int) []discordgo.MessageComponent {
	if pages == 1 {
		return []discordgo.MessageComponent{}
	}
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Emoji:    &discordgo.ComponentEmoji{Name: "⏪"},
					Style:    discordgo.SecondaryButton,
					CustomID: "page_first",
					Disabled: i == 0,
				},
				discordgo.Button{
					Emoji:    &discordgo.ComponentEmoji{Name: "◀️"},
					Style:    discordgo.SecondaryButton,
					CustomID: "page_prev",
					Disabled: i == 0,
				},
				discordgo.Button{
					Emoji:    &discordgo.ComponentEmoji{Name: "▶️"},
					Style:    discordgo.SecondaryButton,
					CustomID: "page_next",
					Disabled: i >= pages-1,
				},
				discordgo.Button{
					Emoji:    &discordgo.ComponentEmoji{Name: "⏩"},
					Style:    discordgo.SecondaryButton,
					CustomID: "page_last",
					Disabled: i >= pages-1,
				},
			},
		},
	}
}
