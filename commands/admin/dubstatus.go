package admin

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/digitalungdom-se/dub/internal"
	"github.com/digitalungdom-se/dub/pkg"
)

var DubStatus = internal.Command{
	Name:        "dubstatus",
	Description: "Ändrar statusen av boten",
	Aliases:     []string{},
	Group:       "admin",
	Usage:       "dubstatus <PLAYING|STREAMING|LISTENING|WATCHING> <status>",
	Example:     "dubstatus WATCHING kelvin's cat",
	ServerOnly:  true,
	AdminOnly:   true,

	Execute: func(ctx *pkg.Context, server *internal.Server) error {
		if len(ctx.Args) == 0 {
			ctx.Reply("Du måste ange en status.")
			return nil
		}

		status := strings.ToUpper(ctx.Args[0])
		statusText := strings.Join(ctx.Args[1:], " ")
		discordStatus := discordgo.UpdateStatusData{}
		discordStatus.Game = new(discordgo.Game)

		switch status {
		case "PLAYING":
			discordStatus.Game.Type = discordgo.GameTypeGame
		case "STREAMING":
			discordStatus.Game.Type = discordgo.GameTypeStreaming
		case "LISTENING":
			discordStatus.Game.Type = discordgo.GameTypeListening
		case "WATCHING":
			discordStatus.Game.Type = discordgo.GameTypeWatching
		default:
			_, err := ctx.Reply("Felaktig status")
			if err != nil {
				return err
			}
			return nil
		}

		discordStatus.Status = statusText
		discordStatus.Game.Name = statusText

		err := ctx.Discord.UpdateStatusComplex(discordStatus)
		if err != nil {
			return err
		}

		server.Status = discordStatus

		ctx.Reply(fmt.Sprintf("Sätter statusen till `%v %v`", status, statusText))

		return nil
	},
}
