package admin

import (
	"github.com/bwmarrin/discordgo"
	"github.com/digitalungdom-se/dub/events"
	"github.com/digitalungdom-se/dub/internal"
	"github.com/digitalungdom-se/dub/pkg"
)

var Join = internal.Command{
	Name:        "join",
	Description: "Simulerar en användare joinar",
	Aliases:     []string{},
	Group:       "admin",
	Usage:       "join @<user>",
	Example:     "join @kelszo#6200",
	ServerOnly:  true,
	AdminOnly:   true,

	Execute: func(ctx *pkg.Context, server *internal.Server) error {
		var member *discordgo.Member

		if len(ctx.Message.Mentions) == 0 {
			ctx.Reply("Du måste ange en användare som ska joina.")
			return nil
		}

		for _, guildMember := range server.Guild.Members {
			if guildMember.User.ID == ctx.Message.Mentions[0].ID {
				member = guildMember
			}
		}

		guildMemberAdd := &discordgo.GuildMemberAdd{Member: member}

		events.GuildMemberAddHandler(server)(ctx.Discord, guildMemberAdd)

		return nil
	},
}
