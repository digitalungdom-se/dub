package music

import (
	"github.com/digitalungdom-se/dub/internal"
	"github.com/digitalungdom-se/dub/pkg"
)

var PauseResume = internal.Command{
	Name:        "pr",
	Description: "Pausar eller återupptar musiken",
	Aliases:     []string{"pause", "resume", "pausa", "fortsätt"},
	Group:       "music",
	Usage:       "pr",
	Example:     "pr",
	ServerOnly:  true,
	AdminOnly:   false,

	Execute: func(ctx *pkg.Context, server *internal.Server) error {
		ctx.Delete()
		server.Controller.PauseResume()

		return nil
	},
}
