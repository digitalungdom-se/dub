package misc

import (
	"github.com/digitalungdom-se/dub/internal"
	"github.com/digitalungdom-se/dub/pkg"
)

var Slap = internal.Command{
	Name:        "slap",
	Description: "Smiska någon som har varit stygg",
	Aliases:     []string{"smisk"},
	Group:       "misc",
	Usage:       "smisk @<user>",
	Example:     "smisk @Ippyson#6200 ",
	ServerOnly:  true,
	AdminOnly:   false,

	Execute: func(ctx *pkg.Context, server *internal.Server) error {
		mentions := ctx.Message.Mentions
		if len(mentions) == 0 {
			return nil
		}

		ctx.ReplyNoMention("Du har varit riktigt stygg " + ctx.Message.Mentions[0].Mention())

		return nil
	},
}
