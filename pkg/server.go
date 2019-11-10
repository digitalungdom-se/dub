package pkg

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
	gomail "gopkg.in/gomail.v2"
)

type (
	Server struct {
		Config           Config
		Controller       *Controller
		CommandHandler   *CommandHandler
		Discord          *discordgo.Session
		Guild            *discordgo.Guild
		Status           *discordgo.UpdateStatusData
		Channels         channels
		Roles            roles
		ReactionListener *ReactionListener
		Database         *mongo.Database
		Dialer           *gomail.Dialer
		Bot              *discordgo.Member
	}

	channels struct {
		General *discordgo.Channel
		Bot     *discordgo.Channel
		Music   *discordgo.Channel
	}

	roles struct {
		Admin *discordgo.Role
	}
)

func NewServer(config Config, discord *discordgo.Session, reactionListener *ReactionListener,
	database *mongo.Database, dialer *gomail.Dialer, commandHandler *CommandHandler) (*Server, error) {

	server := new(Server)
	var err error

	server.Config = config
	server.CommandHandler = commandHandler
	server.Discord = discord
	server.Database = database
	server.Dialer = dialer
	server.ReactionListener = reactionListener
	server.Discord.State.MaxMessageCount = 5

	var guild *discordgo.Guild
	guild, err = discord.Guild(config.GuildID)
	if err != nil {
		return nil, err
	}

	server.Discord.State.GuildAdd(guild)
	server.Guild, err = server.Discord.State.Guild(config.GuildID)
	if err != nil {
		return nil, err
	}

	discordStatus := new(discordgo.UpdateStatusData)
	discordStatus.Game = new(discordgo.Game)

	discord.UpdateStatusComplex(*discordStatus)
	server.Status = discordStatus

	server.Channels = channels{}
	channels, err := discord.GuildChannels(config.GuildID)
	if err != nil {
		return nil, err
	}

	for _, channel := range channels {
		switch channel.Name {
		case "music":
			server.Discord.State.ChannelAdd(channel)
			server.Channels.Music, err = server.Discord.State.Channel(channel.ID)
		case "bot":
			server.Discord.State.ChannelAdd(channel)
			server.Channels.Bot, err = server.Discord.State.Channel(channel.ID)
		case "general":
			server.Discord.State.ChannelAdd(channel)
			server.Channels.General, err = server.Discord.State.Channel(channel.ID)
		}

		if err != nil {
			return nil, err
		}
	}

	for _, role := range server.Guild.Roles {
		if role.Name == "admin" {
			server.Discord.State.RoleAdd(config.GuildID, role)
			server.Roles.Admin, err = server.Discord.State.Role(config.GuildID, role.ID)

			if err != nil {
				return nil, err
			}

			break
		}
	}

	if server.Channels.General == nil ||
		server.Channels.Music == nil ||
		server.Channels.Bot == nil ||
		server.Roles.Admin == nil {
		return nil, errors.New("Could not find channels or roles.")
	}

	bot, err := discord.User("@me")
	if err != nil {
		return nil, err
	}
	for _, member := range server.Guild.Members {
		if member.User.ID == bot.ID {
			err = discord.State.MemberAdd(member)
			if err != nil {
				return nil, err
			}
		}
	}

	server.Bot, err = discord.State.Member(config.GuildID, bot.ID)
	if err != nil {
		return nil, err
	}

	if server.Bot == nil {
		return nil, errors.New("Could not find bot")
	}

	server.Controller, err = NewController(server)

	return server, err
}
