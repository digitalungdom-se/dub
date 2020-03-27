package music

import (
	"github.com/digitalungdom-se/dub/internal"
	"github.com/digitalungdom-se/dub/pkg"
)

var Play = internal.Command{
	Name:        "play",
	Description: "Spelar en låt",
	Aliases:     []string{"spela", "pl"},
	Group:       "music",
	Usage:       "play <youtube link>",
	Example:     "play https://www.youtube.com/watch?v=dQw4w9WgXcQ",
	ServerOnly:  true,
	AdminOnly:   false,

	Execute: func(ctx *pkg.Context, server *internal.Server) error {
		ctx.Delete()

		if len(ctx.Args) == 0 {
			ctx.Reply("Du måste skicka en youtube länk.")
			return nil
		}

		err := server.Controller.AddToQueue(ctx)

		return err
	},
}
