package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/digitalungdom-se/dub/commands/admin"
	"github.com/digitalungdom-se/dub/commands/digitalungdom"
	"github.com/digitalungdom-se/dub/commands/misc"
	"github.com/digitalungdom-se/dub/commands/music"
	"github.com/digitalungdom-se/dub/events"
	"github.com/digitalungdom-se/dub/internal"
	"github.com/digitalungdom-se/dub/pkg"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	config := pkg.LoadConfig()

	commandHandler := pkg.NewCommandHandler()
	registerCommands(commandHandler)

	var mongoClient *mongo.Client
	mongoClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("DB_URI")))
	if err != nil {
		log.Fatal(err)
	}

	err = mongoClient.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	mongoDatabase := mongoClient.Database(os.Getenv("DB_DATABASE"))
	database := internal.NewDatabase(mongoDatabase)

	reactionListener := pkg.NewReactionListener(discord)
	mailer := internal.NewMailer(mongoDatabase.Collection("emails"))

	server := pkg.NewServer(config, discord, reactionListener, database, mailer, commandHandler)
	if err != nil {
		log.Fatal("Error creating server,", err)
	}

	discord.AddHandler(events.GuildMemberAddHandler(server))
	discord.AddHandler(events.MessageHandler(server))
	discord.AddHandler(events.AddReactionHandler(server))
	discord.AddHandler(events.RemoveReactionHandler(server))
	discord.AddHandler(events.StartHandler(server))

	err = discord.Open()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	discord.Close()
}

func registerCommands(commandHandler *pkg.CommandHandler) {
	commandHandler.Register(admin.DubStatus)
	commandHandler.Register(admin.Join)

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
}
