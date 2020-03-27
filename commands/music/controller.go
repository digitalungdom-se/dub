package music

import (
	"github.com/digitalungdom-se/dub/internal"
	"github.com/digitalungdom-se/dub/pkg"
)

var Controller = internal.Command{
	Name:        "controller",
	Description: "Skickar en ny controller i musik kanalen",
	Aliases:     []string{"kontroll"},
	Group:       "music",
	Usage:       "controller",
	Example:     "controller",
	ServerOnly:  true,
	AdminOnly:   false,

	Execute: func(ctx *pkg.Context, server *internal.Server) error {
		ctx.Delete()
		err := server.Controller.NewControllerMessage(&server.ReactionListener)

		return err
	},
}
