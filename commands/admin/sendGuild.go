package admin

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/cbroglie/mustache"
	"github.com/digitalungdom-se/dub/internal"
	"github.com/digitalungdom-se/dub/pkg"
)

var SendGuild = internal.Command{
	Name:        "sendguild",
	Description: "Skickar meddelande till medlemmarna via Discord kanalen",
	Aliases:     []string{},
	Group:       "admin",
	Usage:       "sendguild <channelID> <groupID|userID|\"EVERYONE\"> <message>",
	Example:     "sendguild 573654197703278593 228889878861971456 **Hej** {{{mention}}}, hur mår du?",
	ServerOnly:  true,
	AdminOnly:   true,

	Execute: func(ctx *pkg.Context, server *internal.Server) error {
		var err error

		if len(ctx.Args) < 3 {
			ctx.Reply("Du måste ange channelID, groupID|userID|\"EVERYONE\" och ett meddelande.")
			return nil
		}

		if !pkg.IsSnowflake(ctx.Args[0]) {
			ctx.Reply("channelID måste vara en snowflake.")
			return nil
		}

		var channel *discordgo.Channel
		channel, err = ctx.Discord.Channel(ctx.Args[0])
		if err != nil || channel.GuildID != server.Guild.ID {
			ctx.Reply("ChannelID måste vara ett giltigt ID.")
			return nil
		}

		if !pkg.IsSnowflake(ctx.Args[1]) && strings.ToUpper(ctx.Args[1]) != "EVERYONE" {
			ctx.Reply("Du måste ange en snowflake eller EVERYONE.")
			return nil
		}

		messageTemplate := strings.Join(ctx.Args[2:], " ")
		var message string
		var mention string

		if strings.ToUpper(ctx.Args[1]) == "EVERYONE" {
			mention = "@everyone"
		}

		if pkg.IsSnowflake(ctx.Args[1]) {
			for _, role := range server.Guild.Roles {
				if role.ID == ctx.Args[1] {
					mention = role.Mention()
				}
			}

			if mention == "" {
				for _, member := range server.Guild.Members {
					if member.User.ID == ctx.Args[1] {
						mention = member.Mention()
					}
				}
			}

			if mention == "" {
				ctx.Reply("Du måste ange ett giltigt channelID|userID")
				return nil
			}
		}

		message, err = mustache.Render(messageTemplate, map[string]string{"mention": mention})
		if err != nil {
			ctx.Reply("Error med att kompilera meddelandet.")
			return err
		}

		ctx.Discord.ChannelMessageSend(ctx.Args[0], message)

		return nil
	},
}
