package events

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/digitalungdom-se/dub/internal"
)

func AddReactionHandler(server *internal.Server) func(*discordgo.Session, *discordgo.MessageReactionAdd) {
	return func(discord *discordgo.Session, message *discordgo.MessageReactionAdd) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("%v | RECOVERED ANNA FROM %v\n", time.Now().Format("2006-01-02 15:04:05"), r)
			}
		}()

		user, err := discord.User(message.UserID)
		if err != nil {
			log.Println("Error getting user,", err)
			return
		}

		if user.Bot {
			return
		}

		if message.ChannelID == server.Channels.Regler.ID {
			if server.ReactionListener.Messages[message.MessageID] == nil {
				discord.MessageReactionsRemoveAll(message.ChannelID, message.MessageID)
			}
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
}

func RemoveReactionHandler(server *internal.Server) func(discord *discordgo.Session, message *discordgo.MessageReactionRemove) {
	return func(discord *discordgo.Session, message *discordgo.MessageReactionRemove) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("%v | RECOVERED ANNA FROM %v\n", time.Now().Format("2006-01-02 15:04:05"), r)
			}
		}()

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
}
