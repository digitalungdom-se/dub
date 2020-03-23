package misc

import (
	"github.com/digitalungdom-se/dub/pkg"
)

var Slap = pkg.Command{
	Name:        "slap",
	Description: "Smiska någon som har varit stygg",
	Aliases:     []string{"smisk"},
	Group:       "misc",
	Usage:       "smisk @<user>",
	Example:     "smisk @Ippyson#6200 ",
	ServerOnly:  true,
	AdminOnly:   false,

	Execute: func(ctx *pkg.Context) error {
		mentions := ctx.Message.Mentions
		if len(mentions) == 0 {
			return nil
		}

		ctx.ReplyNoMention("Du har varit riktigt stygg " + ctx.Message.Mentions[0].Mention())

		return nil
	},
}
