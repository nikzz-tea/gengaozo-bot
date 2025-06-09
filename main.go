package main

import (
	"gengaozo/app/api/osu"
	_ "gengaozo/app/commands"
	"gengaozo/app/database"
	_ "gengaozo/app/events"
	"gengaozo/app/handlers"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	token := os.Getenv("TOKEN")

	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	sess.AddHandler(handlers.CommandHandler)
	handlers.EventHandler(sess)

	osuClientId := os.Getenv("OSU_APP_ID")
	osuClientSecret := os.Getenv("OSU_SECRET")
	err = osu.Auth(osuClientId, osuClientSecret)
	if err != nil {
		log.Fatal(err)
	}

	database.Init()

	log.Println("Logged as " + sess.State.User.Username + "#" + sess.State.User.Discriminator)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
