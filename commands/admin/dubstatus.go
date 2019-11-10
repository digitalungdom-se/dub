package admin

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/digitalungdom-se/dub/pkg"
)

var DubStatus = pkg.Command{
	Name:        "dubstatus",
	Description: "Ändrar statusen av boten",
	Aliases:     []string{},
	Group:       "admin",
	Usage:       "dubstatus <PLAYING|STREAMING|LISTENING|WATCHING> <status>",
	Example:     "dubstatus WATCHING kelvin's cat",
	ServerOnly:  true,
	AdminOnly:   true,

	Execute: func(context *pkg.Context) error {
		status := strings.ToUpper(context.Args[0])
		statusText := strings.Join(context.Args[1:], " ")
		discordStatus := new(discordgo.UpdateStatusData)
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
			context.Reply("Felaktig status")
			return nil
		}

		discordStatus.Status = statusText
		discordStatus.Game.Name = statusText

		context.Discord.UpdateStatusComplex(*discordStatus)
		context.Server.Status = discordStatus

		return nil
	},
}
