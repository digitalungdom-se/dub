package misc

import (
	"github.com/digitalungdom-se/dub/internal"
	"github.com/digitalungdom-se/dub/pkg"
)

var Info = internal.Command{
	Name:        "info",
	Description: "Få information om boten",
	Aliases:     []string{},
	Group:       "misc",
	Usage:       "info",
	Example:     "info",
	ServerOnly:  false,
	AdminOnly:   false,

	Execute: func(ctx *pkg.Context, server *internal.Server) error {
		embed := pkg.NewEmbed().
			SetTitle("__**INFORMATION OM BOTEN**__").
			AddField("VERSION", "1.0").
			AddField("KÄLLKOD", "[github](https://github.com/digitalungdom-se/dub)").
			AddField("MEDARBETARE", "<@217632464531619852>").
			InlineAllFields().
			SetColor(4086462).MessageEmbed

		_, err := ctx.ReplyEmbed(embed)

		return err
	},
}
