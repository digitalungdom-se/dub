package music

import (
	"github.com/digitalungdom-se/dub/internal"
	"github.com/digitalungdom-se/dub/pkg"
)

var Stop = internal.Command{
	Name:        "stop",
	Description: "Stannar botens nuvarande musik",
	Aliases:     []string{"stanna", "st"},
	Group:       "music",
	Usage:       "stop",
	Example:     "stop",
	ServerOnly:  true,
	AdminOnly:   false,

	Execute: func(ctx *pkg.Context, server *internal.Server) error {
		ctx.Delete()
		server.Controller.Stop()

		return nil
	},
}
