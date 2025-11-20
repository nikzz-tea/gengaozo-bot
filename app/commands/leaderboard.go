package commands

import (
	"fmt"
	"gengaozo/app/api/osu"
	"gengaozo/app/database"
	"gengaozo/app/handlers"
	"gengaozo/app/models"
	"gengaozo/app/utils"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/language"
	m "golang.org/x/text/message"
)

func init() {
	handlers.RegisterCommand(models.CommandObject{
		Name:    "leaderboard",
		Aliases: []string{"lb"},
		Callback: func(props models.CommandProps) {
			sess, message, args := props.Sess, props.Message, props.Args

			sess.ChannelTyping(message.ChannelID)

			beatmapID, mods := handleArguments(sess, message, args)

			if beatmapID == "" {
				sess.ChannelMessageSend(message.ChannelID, "🔴 **Provide a beatmap id**")
				return
			}

			beatmap, err := osu.GetBeatmap(beatmapID)
			if err != nil || beatmap.Beatmapset.Artist == "" {
				sess.ChannelMessageSend(message.ChannelID, "🔴 **No beatmap was found**")
				return
			}

			members, err := sess.GuildMembers(message.GuildID, "", 1000)
			if err != nil {
				log.Println(err)
				return
			}
			memberIDs := make([]string, 0, len(members))
			for _, member := range members {
				memberIDs = append(memberIDs, member.User.ID)
			}
			var osuUsers []database.User
			database.DB.Model(&database.User{}).
				Select("osu_id", "osu_username").
				Where("discord_id IN ?", memberIDs).
				Find(&osuUsers)

			results := make(chan models.BeatmapScores, len(osuUsers))
			errors := make(chan error, len(osuUsers))
			for _, user := range osuUsers {
				go func(user database.User) {
					scores, err := osu.GetBeatmapScores(user.Osu_id, beatmapID)
					if err != nil {
						errors <- err
						return
					}
					for i := range scores.Scores {
						scores.Scores[i].User = user
					}
					results <- scores
				}(user)
			}

			var scores []models.Score
			for range osuUsers {
				select {
				case userScores := <-results:
					if len(userScores.Scores) > 0 {
						filteredScores := filterByMods(userScores.Scores, mods)
						if len(filteredScores) > 0 {
							sortScores(filteredScores)
							scores = append(scores, filteredScores[0])
						}
					}
				case err := <-errors:
					log.Println(err)
				}
			}

			if len(scores) == 0 {
				sess.ChannelMessageSend(message.ChannelID, "🔴 **No scores were found**")
				return
			}

			sortScores(scores)

			scoresPerPage := 5
			totalPages := (len(scores) + scoresPerPage - 1) / scoresPerPage
			var pages []*discordgo.MessageEmbed

			for i := 0; i < len(scores); i += scoresPerPage {
				pageScores := scores[i:min(i+scoresPerPage, len(scores))]
				starRating := fmt.Sprintf("%.2f", beatmap.StarRating)

				title := fmt.Sprintf(
					"%v - %v [%v] %v\\*",
					beatmap.Beatmapset.Artist, beatmap.Beatmapset.Title, beatmap.Diffname, starRating,
				)
				footer := fmt.Sprintf("Page %v/%v", len(pages)+1, totalPages)
				embed := &discordgo.MessageEmbed{
					Title:     title,
					URL:       "https://osu.ppy.sh/b/" + beatmapID,
					Thumbnail: &discordgo.MessageEmbedThumbnail{URL: beatmap.Beatmapset.Covers.List},
					Color:     0x637191,
					Footer:    &discordgo.MessageEmbedFooter{Text: footer},
					Fields:    []*discordgo.MessageEmbedField{},
				}

				for j, score := range pageScores {
					mods := "NM"
					if len(score.Mods) > 0 {
						mods = strings.Join(score.Mods, "")
					}
					pp := math.Round(score.PP)
					accuracy := fmt.Sprintf("%.2f", score.Accuracy*100)
					date, _ := time.Parse(time.RFC3339, score.Date)
					timestamp := fmt.Sprintf("<t:%v:R>", date.Unix())
					locale := m.NewPrinter(language.English)
					formattedScore := locale.Sprintf("%d", score.Score)

					name := fmt.Sprintf(
						"**#%v** %v `%vpp` +%v",
						scoresPerPage*len(pages)+j+1, score.User.Osu_username, pp, mods,
					)
					value := fmt.Sprintf(
						"%v (%v%%) • %v • **x%v/%v** • [%v/%v/%v/%v]\nScore set %v",
						utils.GetRankEmote(score.Rank), accuracy, formattedScore, score.MaxCombo,
						beatmap.MaxCombo, score.Hits.Count300, score.Hits.Count100,
						score.Hits.Count50, score.Hits.CountMiss, timestamp,
					)

					embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{Name: name, Value: value})
				}

				pages = append(pages, embed)
			}

			msg, _ := sess.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
				Embeds:     []*discordgo.MessageEmbed{pages[0]},
				Components: utils.GetPaginationButtons(0, totalPages),
			})

			if totalPages > 1 {
				utils.CreatePagination(sess, pages, msg.ID, msg.ChannelID)
			}
		},
	})
}

func handleArguments(s *discordgo.Session, m *discordgo.MessageCreate, args []string) (string, []string) {
	var beatmapID string
	var mods []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "+") {
			raw := strings.TrimPrefix(arg, "+")
			if len(raw)%2 != 0 {
				continue
			}
			for i := 0; i < len(raw); i += 2 {
				mods = append(mods, strings.ToUpper(raw[i:i+2]))
			}
		} else if parsed := utils.ParseBeatmapID(arg); parsed != "" {
			beatmapID = parsed
		} else if _, err := strconv.Atoi(arg); err == nil {
			beatmapID = arg
		}
	}

	if beatmapID == "" && m.MessageReference != nil {
		replied, err := s.ChannelMessage(m.MessageReference.ChannelID, m.MessageReference.MessageID)
		if err == nil && len(replied.Embeds) > 0 {
			beatmapID = utils.ParseEmbed(replied.Embeds[0])
		}
	}

	if beatmapID == "" {
		database.DB.Model(&database.Map{}).
			Select("map_id").
			Where("channel_id = ?", m.ChannelID).
			Take(&beatmapID)
	}

	return beatmapID, mods
}

func sortScores(scores []models.Score) {
	sort.Slice(scores, func(i, j int) bool {
		if scores[i].PP == 0 {
			return scores[i].Score > scores[j].Score
		}
		return scores[i].PP > scores[j].PP
	})
}

func filterByMods(scores []models.Score, mods []string) []models.Score {
	if len(mods) == 0 {
		return scores
	}

	var filteredScores []models.Score

	for _, score := range scores {
		if len(score.Mods) == 0 {
			score.Mods = append(score.Mods, "NM")
		}
		match := true
		modsMap := make(map[string]bool)
		for _, mod := range score.Mods {
			mod := strings.Replace(mod, "NC", "DT", 1)
			modsMap[mod] = true
		}
		for _, mod := range mods {
			mod := strings.Replace(mod, "NC", "DT", 1)
			if !modsMap[mod] {
				match = false
				break
			}
		}
		if match {
			filteredScores = append(filteredScores, score)
		}
	}

	return filteredScores
}
