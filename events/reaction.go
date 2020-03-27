package events

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/digitalungdom-se/dub/internal"
)

func AddReactionHandler(server *internal.Server) func(*discordgo.Session, *discordgo.MessageReactionAdd) {
	return func(discord *discordgo.Session, message *discordgo.MessageReactionAdd) {
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
