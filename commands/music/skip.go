package music

import (
	"github.com/digitalungdom-se/dub/internal"
	"github.com/digitalungdom-se/dub/pkg"
)

var Skip = internal.Command{
	Name:        "skip",
	Description: "Skippar den nuvarande låten",
	Aliases:     []string{"skippa", "byt", "sk"},
	Group:       "music",
	Usage:       "skip",
	Example:     "skip",
	ServerOnly:  true,
	AdminOnly:   false,

	Execute: func(ctx *pkg.Context, server *internal.Server) error {
		ctx.Delete()
		server.Controller.Skip()

		return nil
	},
}
