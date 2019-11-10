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

	Execute: func(context *pkg.Context) error {
		context.Delete()
		context.Server.Controller.Stop()

		return nil
	},
}
