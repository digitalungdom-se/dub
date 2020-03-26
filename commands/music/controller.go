package music

import (
	"github.com/digitalungdom-se/dub/pkg"
)

var Controller = pkg.Command{
	Name:        "controller",
	Description: "Skickar en ny controller i musik kanalen",
	Aliases:     []string{"kontroll"},
	Group:       "music",
	Usage:       "controller",
	Example:     "controller",
	ServerOnly:  true,
	AdminOnly:   false,

	Execute: func(ctx *pkg.Context) error {
		ctx.Delete()
		err := ctx.Server.Controller.NewControllerMessage()

		return err
	},
}
