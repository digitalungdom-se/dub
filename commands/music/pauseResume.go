package music

import (
	"github.com/digitalungdom-se/dub/pkg"
)

var PauseResume = pkg.Command{
	Name:        "pr",
	Description: "Pausar eller återupptar musiken",
	Aliases:     []string{"pause", "resume", "pausa", "fortsätt"},
	Group:       "music",
	Usage:       "pr",
	Example:     "pr",
	ServerOnly:  true,
	AdminOnly:   false,

	Execute: func(context *pkg.Context) error {
		context.Delete()
		context.Server.Controller.PauseResume()

		return nil
	},
}
