package music

import (
	"github.com/digitalungdom-se/dub/pkg"
)

var Play = pkg.Command{
	Name:        "play",
	Description: "Spelar en låt",
	Aliases:     []string{"spela", "pl"},
	Group:       "music",
	Usage:       "play <youtube link>",
	Example:     "play https://www.youtube.com/watch?v=dQw4w9WgXcQ",
	ServerOnly:  true,
	AdminOnly:   false,

	Execute: func(ctx *pkg.Context) error {
		ctx.Delete()

		if len(ctx.Args) == 0 {
			ctx.Reply("Du måste skicka en youtube länk.")
			return nil
		}

		err := ctx.Server.Controller.AddToQueue(ctx)

		return err
	},
}
