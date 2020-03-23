package music

import (
	"github.com/digitalungdom-se/dub/pkg"
)

var Stop = pkg.Command{
	Name:        "stop",
	Description: "Stannar botens nuvarande musik",
	Aliases:     []string{"stanna", "st"},
	Group:       "music",
	Usage:       "stop",
	Example:     "stop",
	ServerOnly:  true,
	AdminOnly:   false,

	Execute: func(ctx *pkg.Context) error {
		ctx.Delete()
		ctx.Server.Controller.Stop()

		return nil
	},
}
