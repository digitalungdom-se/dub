package pkg

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/digitalungdom-se/dub/internal"
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
		Database         *internal.Database
		Mailer           *internal.Mailer
		Bot              *discordgo.Member
		Ready            bool
	}

	channels struct {
		General *discordgo.Channel
		Bot     *discordgo.Channel
		Music   *discordgo.Channel
		Regler  *discordgo.Channel
	}

	roles struct {
		Admin *discordgo.Role
	}
)

func NewServer(config Config, discord *discordgo.Session, reactionListener *ReactionListener,
	database *internal.Database, mailer *internal.Mailer, commandHandler *CommandHandler) *Server {

	server := new(Server)

	server.Config = config
	server.CommandHandler = commandHandler
	server.Discord = discord
	server.Database = database
	server.Mailer = mailer
	server.ReactionListener = reactionListener
	server.Discord.State.MaxMessageCount = 5
	server.Ready = false

	return server
}

func (server *Server) Init() error {
	guild, err := server.Discord.Guild(server.Config.GuildID)
	if err != nil {
		return err
	}

	server.Discord.State.GuildAdd(guild)
	server.Guild, err = server.Discord.State.Guild(server.Config.GuildID)
	if err != nil {
		return err
	}

	discordStatus := new(discordgo.UpdateStatusData)
	discordStatus.Game = new(discordgo.Game)

	server.Discord.UpdateStatusComplex(*discordStatus)
	server.Status = discordStatus

	server.Channels = channels{}
	channels, err := server.Discord.GuildChannels(server.Config.GuildID)
	if err != nil {
		return err
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
		case "regler":
			server.Discord.State.ChannelAdd(channel)
			server.Channels.Regler, err = server.Discord.State.Channel(channel.ID)
		}

		if err != nil {
			return err
		}
	}

	for _, role := range server.Guild.Roles {
		if role.Name == "admin" {
			server.Discord.State.RoleAdd(server.Config.GuildID, role)
			server.Roles.Admin, err = server.Discord.State.Role(server.Config.GuildID, role.ID)

			if err != nil {
				return err
			}

			break
		}
	}

	if server.Channels.General == nil ||
		server.Channels.Music == nil ||
		server.Channels.Bot == nil ||
		server.Roles.Admin == nil {
		return errors.New("could not find channels or roles")
	}

	bot, err := server.Discord.User("@me")
	if err != nil {
		return err
	}
	for _, member := range server.Guild.Members {
		if member.User.ID == bot.ID {
			err = server.Discord.State.MemberAdd(member)
			if err != nil {
				return err
			}
		}
	}

	server.Bot, err = server.Discord.State.Member(server.Config.GuildID, bot.ID)
	if err != nil {
		return err
	}

	if server.Bot == nil {
		return errors.New("could not find bot")
	}

	server.Controller, err = NewController(server)

	server.Ready = true

	return err
}
