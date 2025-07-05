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
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/language"
	m "golang.org/x/text/message"
)

var paginations = make(map[string]*models.PaginationData)
var paginationMutex sync.Mutex
var cleanupDelay = 30 * time.Second

func init() {
	handlers.RegisterCommand(models.CommandObject{
		Name:    "leaderboard",
		Aliases: []string{"lb"},
		Callback: func(props models.CommandProps) {
			sess, message, args := props.Sess, props.Message, props.Args

			sess.ChannelTyping(message.ChannelID)

			beatmapID := getBeatmapID(sess, message, args)

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
						sort.Slice(userScores.Scores, func(i, j int) bool {
							if userScores.Scores[i].PP == 0 {
								return userScores.Scores[i].Score > userScores.Scores[j].Score
							}
							return userScores.Scores[i].PP > userScores.Scores[j].PP
						})
						scores = append(scores, userScores.Scores[0])
					}
				case err := <-errors:
					log.Println(err)
				}
			}

			if len(scores) == 0 {
				sess.ChannelMessageSend(message.ChannelID, "🔴 **No scores were found**")
				return
			}

			sort.Slice(scores, func(i, j int) bool {
				if scores[i].PP == 0 {
					return scores[i].Score > scores[j].Score
				}
				return scores[i].PP > scores[j].PP
			})

			scoresPerPage := 5
			totalPages := (len(scores) + scoresPerPage - 1) / scoresPerPage
			var pages []*discordgo.MessageEmbed

			for i := 0; i < len(scores); i += scoresPerPage {
				pageScores := scores[i:min(i+scoresPerPage, len(scores))]

				title := fmt.Sprintf(
					"%v - %v [%v] %v*",
					beatmap.Beatmapset.Artist, beatmap.Beatmapset.Title, beatmap.Diffname, beatmap.StarRating,
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
				timer := time.AfterFunc(cleanupDelay, func() {
					cleanupPagination(sess, msg.ID, msg.ChannelID)
				})

				paginationMutex.Lock()
				paginations[msg.ID] = &models.PaginationData{
					Pages:       pages,
					CurrentPage: 0,
					LastUsed:    time.Now(),
					Timer:       timer,
				}
				paginationMutex.Unlock()
			}
		},
	})
	handlers.RegisterEvent(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionMessageComponent {
			return
		}

		messageID := i.Message.ID

		paginationMutex.Lock()
		pagination, exists := paginations[messageID]
		paginationMutex.Unlock()
		if !exists {
			return
		}

		pagination.Timer.Stop()
		pagination.Timer.Reset(cleanupDelay)
		pagination.LastUsed = time.Now()

		switch i.MessageComponentData().CustomID {
		case "page_first":
			pagination.CurrentPage = 0
		case "page_prev":
			pagination.CurrentPage--
		case "page_next":
			pagination.CurrentPage++
		case "page_last":
			pagination.CurrentPage = len(pagination.Pages) - 1
		default:
			return
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds:     []*discordgo.MessageEmbed{pagination.Pages[pagination.CurrentPage]},
				Components: utils.GetPaginationButtons(pagination.CurrentPage, len(pagination.Pages)),
			},
		})
	})
}

func getBeatmapID(s *discordgo.Session, m *discordgo.MessageCreate, args []string) string {
	var beatmapID string

	if len(args) > 0 {
		if parsed := utils.ParseBeatmapID(args[0]); parsed != "" {
			beatmapID = parsed
		} else if _, err := strconv.Atoi(args[0]); err == nil {
			beatmapID = args[0]
		}
	}
	if beatmapID != "" {
		return beatmapID
	}

	if m.MessageReference != nil {
		replied, err := s.ChannelMessage(m.MessageReference.ChannelID, m.MessageReference.MessageID)
		if err == nil && len(replied.Embeds) > 0 {
			beatmapID = utils.ParseEmbed(replied.Embeds[0])
		}
	}
	if beatmapID != "" {
		return beatmapID
	}

	database.DB.Model(&database.Map{}).
		Select("map_id").
		Where("channel_id = ?", m.ChannelID).
		Take(&beatmapID)
	if beatmapID != "" {
		return beatmapID
	}

	return beatmapID
}

func cleanupPagination(s *discordgo.Session, messageID, channelID string) {
	paginationMutex.Lock()
	defer paginationMutex.Unlock()

	pagination, exists := paginations[messageID]
	if !exists {
		return
	}

	s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel:    channelID,
		ID:         messageID,
		Embed:      pagination.Pages[pagination.CurrentPage],
		Components: &[]discordgo.MessageComponent{},
	})

	delete(paginations, messageID)
}
