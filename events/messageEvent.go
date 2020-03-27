package events

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/digitalungdom-se/dub/internal"
	"github.com/digitalungdom-se/dub/pkg"
)

func MessageHandler(server *internal.Server) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(discord *discordgo.Session, message *discordgo.MessageCreate) {
		if !server.Ready {
			return
		}

		channel, err := discord.State.Channel(message.ChannelID)
		if err != nil {
			log.Println("Error getting channel,", err)

			return
		}

		if channel.Name == "music" {
			time.AfterFunc(5*time.Second, func() {
				if message.ID != server.Controller.Message.ID {
					err := discord.ChannelMessageDelete(message.ChannelID, message.ID)
					if err != nil {
						log.Print("Error deleting message", err)

						return
					}
				}
			})

		}

		if message.Author.Bot {
			return
		}

		if len([]rune(message.Content)) == 0 || !pkg.StringInSlice(string([]rune(message.Content)[0]), server.Config.Prefix) {
			return
		}

		args := strings.Fields(message.Content)
		commandName := strings.ToLower(args[0][1:])
		args = args[1:]

		command, found := server.CommandHandler.GetCommand(commandName)

		if !found {
			replyMessage := "Kunde inte hitta kommandot. Testa `$help` för att se alla kommandon."

			if channel.Type == discordgo.ChannelTypeGuildText {
				replyMessage = message.Author.Mention() + ", " + strings.ToLower(string(replyMessage[0])) + replyMessage[1:]
			}

			discord.ChannelMessageSend(message.ChannelID, replyMessage)

			return
		}

		if command.ServerOnly {
			if channel.Type != discordgo.ChannelTypeGuildText {
				replyMessage := "Kommandot är bara tillgängligt i servern, var snäll och använd den där."

				if channel.Type == discordgo.ChannelTypeGuildText {
					replyMessage = message.Author.Mention() + ", " + strings.ToLower(string(replyMessage[0])) + replyMessage[1:]
				}

				discord.ChannelMessageSend(message.ChannelID, replyMessage)

				return
			}
		}

		if command.AdminOnly {
			for _, member := range server.Guild.Members {
				if member.User.ID == message.Author.ID {
					if !pkg.StringInSlice(server.Roles.Admin.ID, member.Roles) {
						replyMessage := "Detta kommando är bara tillgänglig för admins."

						if channel.Type == discordgo.ChannelTypeGuildText {
							replyMessage = message.Author.Mention() + ", " + strings.ToLower(string(replyMessage[0])) + replyMessage[1:]
						}

						discord.ChannelMessageSend(message.ChannelID, replyMessage)

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
				replyMessage := "Du måste vara i en ljudkanal för att styra boten."

				if channel.Type == discordgo.ChannelTypeGuildText {
					replyMessage = message.Author.Mention() + ", " + strings.ToLower(string(replyMessage[0])) + replyMessage[1:]
				}

				discord.ChannelMessageSend(message.ChannelID, replyMessage)

				return
			}

			if botVS != nil && botVS.ChannelID != userVS.ChannelID {
				replyMessage := "Du måste vara i samma ljudkanal som boten för att styra den."

				if channel.Type == discordgo.ChannelTypeGuildText {
					replyMessage = message.Author.Mention() + ", " + strings.ToLower(string(replyMessage[0])) + replyMessage[1:]
				}

				discord.ChannelMessageSend(message.ChannelID, replyMessage)

				return
			}
		}

		ctx := pkg.NewContext(discord, channel, message, args)

		err = command.Execute(ctx, server)
		if err != nil {
			log.Print(fmt.Sprintf("Error executing: %v ", message.Content), err)
		}
	}
}
