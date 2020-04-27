package misc

import (
	"fmt"

	"github.com/digitalungdom-se/dub/assets"
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
		var maintainers string
		var contributors string

		for _, contributor := range assets.Contributors {
			if contributor.Role == assets.Maintainer {
				if maintainers == "" {
					maintainers += fmt.Sprintf("<@%v>", contributor.DiscordID)
				} else {
					maintainers += fmt.Sprintf("<@%v>\n", contributor.DiscordID)
				}
			}

			if contributor.Role == assets.Contributor {
				if contributors == "" {
					contributors += fmt.Sprintf("<@%v>", contributor.DiscordID)
				} else {
					contributors += fmt.Sprintf(", <@%v>", contributor.DiscordID)
				}
			}
		}

		embed := pkg.NewEmbed().
			SetTitle("__**INFORMATION OM BOTEN**__").
			AddField("VERSION", "1.1").
			AddField("KÄLLKOD", "[github](https://github.com/digitalungdom-se/dub)").
			AddField("MAINTAINERS", maintainers).
			AddField("MEDARBETARE", contributors).
			InlineAllFields().
			SetColor(4086462).MessageEmbed

		_, err := ctx.ReplyEmbed(embed)

		return err
	},
}
