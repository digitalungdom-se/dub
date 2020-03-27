package misc

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/digitalungdom-se/dub/internal"
	"github.com/digitalungdom-se/dub/pkg"
)

var Ping = internal.Command{
	Name:        "ping",
	Description: "ping, pong!",
	Aliases:     []string{},
	Group:       "misc",
	Usage:       "ping",
	Example:     "ping",
	ServerOnly:  false,
	AdminOnly:   false,

	Execute: func(ctx *pkg.Context, server *internal.Server) error {
		messageTime, err := discordgo.SnowflakeTimestamp(ctx.Message.ID)
		if err != nil {
			return err
		}

		timeNow := time.Now().UnixNano() / 1000000

		ping := timeNow - (messageTime.UnixNano() / 1000000)

		embed := pkg.NewEmbed().
			SetTitle(":ping_pong:").
			SetDescription(fmt.Sprintf("%vms", ping)).
			SetColor(4086462).MessageEmbed

		_, err = ctx.ReplyEmbed(embed)

		return err
	},
}
