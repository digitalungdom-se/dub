package misc

import (
	"fmt"
	"time"

	"github.com/digitalungdom-se/dub/pkg"
	"github.com/hassieswift621/discord-goflake"
)

var Ping = pkg.Command{
	Name:        "ping",
	Description: "ping, pong!",
	Aliases:     []string{},
	Group:       "misc",
	Usage:       "ping",
	Example:     "ping",
	ServerOnly:  false,
	AdminOnly:   false,

	Execute: func(context *pkg.Context) error {
		snowflake, err := dgoflake.ParseString(context.Message.ID)

		if err != nil {
			return err
		}

		messageTime := snowflake.Timestamp().UnixNano() / 1000000
		timeNow := time.Now().UnixNano() / 1000000

		ping := timeNow - messageTime

		embed := pkg.NewEmbed().
			SetTitle(":ping_pong:").
			SetDescription(fmt.Sprintf("%vms", ping)).
			SetColor(4086462).MessageEmbed

		context.ReplyEmbed(embed)

		return nil
	},
}
