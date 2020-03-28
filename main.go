package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/digitalungdom-se/dub/commands/admin"
	"github.com/digitalungdom-se/dub/commands/digitalungdom"
	"github.com/digitalungdom-se/dub/commands/misc"
	"github.com/digitalungdom-se/dub/commands/music"
	"github.com/digitalungdom-se/dub/events"
	"github.com/digitalungdom-se/dub/internal"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	botToken := os.Getenv("BOT_TOKEN")

	discord, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatal("Error creating Discord session,", err)
	}

	config := internal.LoadConfig()

	server := internal.NewServer(config, discord)
	if err != nil {
		log.Fatal("Error creating server,", err)
	}

	registerCommands(&server.CommandHandler)

	discord.AddHandler(events.GuildMemberAddHandler(server))
	discord.AddHandler(events.MessageHandler(server))
	discord.AddHandler(events.AddReactionHandler(server))
	discord.AddHandler(events.RemoveReactionHandler(server))
	discord.AddHandlerOnce(events.StartHandler(server))

	err = discord.Open()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v | STARTED ANNA IN DISCORD\n", time.Now().Format("2006-01-02 15:04:05"))
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	err = discord.Close()
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}

func registerCommands(commandHandler *internal.CommandHandler) {
	commandHandler.Register(admin.DubStatus)
	commandHandler.Register(admin.Join)
	commandHandler.Register(admin.SendGuild)
	commandHandler.Register(admin.SendDM)

	commandHandler.Register(digitalungdom.Bug)
	commandHandler.Register(digitalungdom.Idea)
	commandHandler.Register(digitalungdom.Report)
	commandHandler.Register(digitalungdom.Status)
	commandHandler.Register(digitalungdom.Verify)
	commandHandler.Register(digitalungdom.Whois)

	commandHandler.Register(misc.Ping)
	commandHandler.Register(misc.Slap)
	commandHandler.Register(misc.Help)
	commandHandler.Register(misc.Info)

	commandHandler.Register(music.Play)
	commandHandler.Register(music.Skip)
	commandHandler.Register(music.Stop)
	commandHandler.Register(music.PauseResume)
	commandHandler.Register(music.Controller)
}
