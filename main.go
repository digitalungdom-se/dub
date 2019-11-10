package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/digitalungdom-se/dub/commands/admin"
	"github.com/digitalungdom-se/dub/commands/digitalungdom"
	"github.com/digitalungdom-se/dub/commands/misc"
	"github.com/digitalungdom-se/dub/commands/music"
	"github.com/digitalungdom-se/dub/pkg"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/gomail.v2"
)

var (
	server *pkg.Server
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	botToken := os.Getenv("BOT_TOKEN")

	discord, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
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

	database := mongoClient.Database("digitalungdom")
	reactionListener := pkg.NewReactionListener(discord)
	dialer := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("NOREPLY_EMAIL"), os.Getenv("NOREPLY_PASSWORD"))

	err = discord.Open()
	if err != nil {
		log.Fatal(err)
	}

	server, err = pkg.NewServer(config, discord, reactionListener, database, dialer, commandHandler)
	if err != nil {
		log.Fatal("Error creating server,", err)
	}

	discord.AddHandler(messageHandler)
	discord.AddHandler(addReactionHandler)
	discord.AddHandler(removeReactionHandler)

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}

func messageHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	channel, err := discord.State.Channel(message.ChannelID)
	if err != nil {
		log.Println("Error getting channel,", err)
		return
	}

	if channel.Name == "music" {
		if message.ID != server.Controller.Message.ID {
			err := discord.ChannelMessageDelete(message.ChannelID, message.ID)
			if err != nil {
				log.Print("Error deleting message", err)
				return
			}
		}
	}

	if message.Author.Bot {
		return
	}

	if len(message.Mentions) > 0 {
		member, err := discord.GuildMember(server.Config.GuildID, message.Mentions[0].ID)
		if err != nil {
			log.Print("Error adding member")
			return
		}

		memberAdd := &discordgo.GuildMemberAdd{Member: member}

		guildMemberAddHandler(discord, memberAdd)
		return
	}

	if !pkg.StringInSlice(string([]rune(message.Content)[0]), server.Config.Prefix) {
		return
	}

	args := strings.Fields(message.Content)
	commandName := strings.ToLower(args[0][1:])
	args = args[1:]

	command, found := server.CommandHandler.GetCommand(commandName)

	if !found {
		discord.ChannelMessageSend(message.ChannelID,
			"Kunde inte hitta kommandot. Testa `$help` för att se alla kommandon.")
		return
	}

	if command.ServerOnly {
		if channel.Type != discordgo.ChannelTypeGuildText {
			discord.ChannelMessageSend(message.ChannelID,
				"Kommandot är bara tillgängligt i servern, var snäll och använd den där.")
			return
		}
	}

	if command.AdminOnly {
		for _, member := range server.Guild.Members {
			if member.User.ID == message.Author.ID {
				if !pkg.StringInSlice(server.Roles.Admin.ID, member.Roles) {
					discord.ChannelMessageSend(message.ChannelID,
						"Detta kommando är bara tillgänglig för admins.")
					return
				}
			}
		}
	}

	if command.Group == "music" {
		var botVS *discordgo.VoiceState
		var userVS *discordgo.VoiceState

		for _, vs := range server.Guild.VoiceStates {
			switch vs.UserID {
			case message.Author.ID:
				userVS = vs
			case server.Bot.User.ID:
				botVS = vs
			}

		}

		if userVS == nil {
			discord.ChannelMessageSend(message.ChannelID,
				"Du måste vara i en ljudkanal för att styra boten.")

			return
		}

		if botVS != nil && botVS.ChannelID != userVS.ChannelID {
			discord.ChannelMessageSend(message.ChannelID,
				"Du måste vara i samma ljudkanal som boten för att styra den.")

			return

		}
	}

	ctx := pkg.NewContext(discord, server, channel, message, args)

	err = command.Execute(ctx)
	if err != nil {
		log.Print(fmt.Sprintf("Error executing: %v ", message.Content), err)
	}
}

func guildMemberAddHandler(discord *discordgo.Session, member *discordgo.GuildMemberAdd) {
	gifbuff, err := pkg.NameToGif(member.User.Username, member.User.AvatarURL("128"))
	if err != nil {
		log.Print("Error creating user gif:", err)
		return
	}
	reader := bytes.NewReader((*gifbuff).Bytes())

	_, err = discord.ChannelFileSendWithMessage(server.Channels.General.ID,
		fmt.Sprintf("Välkommen till Digital Ungdoms Discord server <@%v>!", member.User.ID),
		"welcome.gif",
		reader)
	if err != nil {
		log.Print("Error sending welcome message:", err)
		return
	}

	var privateDM *discordgo.Channel

	privateDM, err = discord.UserChannelCreate(member.User.ID)
	if err != nil {
		log.Print("Error creating private channel:", err)
		return
	}

	content := fmt.Sprintf("Hej **%v** och välkommen till *Digital Ungdoms* Discord server."+
		" Jag är boten som hjälper till med kanalen. För att se alla mina funktioner skriv `$help` till mig.\n\n"+
		"Om du inte redan har ett *Digital Ungdom* konto så rekommenderar jag starkt att du skaffar ett."+
		" Som medlem kan du bland annat skriva på vårt forum (https://digitalungdom.se/agora)."+
		" Du kan enkelt bli medlem genom följande länk: https://digitalungdom.se/bli-medlem\n\n"+
		"Om du redan är medlem så kan du koppla ditt *Digital Ungdom* konto till ditt Discord konto"+
		" genom att skriva `$verify` och sedan ditt användarnamn eller epost. Till exempel `$verify username`\n\n"+
		"Vi synns där inne!",
		member.User.Username)

	_, err = discord.ChannelMessageSend(privateDM.ID, content)
	if err != nil {
		log.Print("Error sending private welcome message:", err)
		return
	}
}

func addReactionHandler(discord *discordgo.Session, message *discordgo.MessageReactionAdd) {
	user, err := discord.User(message.UserID)
	if err != nil {
		log.Println("Error getting user,", err)
		return
	}

	isBot := user.Bot
	if isBot {
		return
	}

	if server.ReactionListener.Messages[message.MessageID] == nil {
		return
	}

	err = server.ReactionListener.React(message.MessageReaction)

	if err != nil {
		log.Print("Error reacting to message:", err)
		return
	}
}

func removeReactionHandler(discord *discordgo.Session, message *discordgo.MessageReactionRemove) {
	user, err := discord.User(message.UserID)
	if err != nil {
		log.Println("Error getting user,", err)
		return
	}

	isBot := user.Bot
	if isBot {
		return
	}

	if server.ReactionListener.Messages[message.MessageID] == nil ||
		!server.ReactionListener.Messages[message.MessageID].ListenToRemove {
		return
	}

	err = server.ReactionListener.React(message.MessageReaction)
	if err != nil {
		log.Print("Error reacting to message:", err)
		return
	}
}

func registerCommands(commandHandler *pkg.CommandHandler) {
	commandHandler.Register(&admin.DubStatus)

	commandHandler.Register(&digitalungdom.Bug)
	commandHandler.Register(&digitalungdom.Idea)
	commandHandler.Register(&digitalungdom.Report)
	commandHandler.Register(&digitalungdom.Status)
	commandHandler.Register(&digitalungdom.Verify)

	commandHandler.Register(&misc.Ping)
	commandHandler.Register(&misc.Slap)
	commandHandler.Register(&misc.Help)
	commandHandler.Register(&misc.Info)

	commandHandler.Register(&music.Play)
	commandHandler.Register(&music.Skip)
	commandHandler.Register(&music.Stop)
	commandHandler.Register(&music.PauseResume)
}
