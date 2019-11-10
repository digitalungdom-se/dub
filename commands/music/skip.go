package music

import (
	"github.com/digitalungdom-se/dub/pkg"
)

var Skip = pkg.Command{
	Name:        "skip",
	Description: "Skippar den nuvarande låten",
	Aliases:     []string{"skippa", "byt", "sk"},
	Group:       "music",
	Usage:       "skip",
	Example:     "skip",
	ServerOnly:  true,
	AdminOnly:   false,

	Execute: func(context *pkg.Context) error {
		context.Delete()
		context.Server.Controller.Skip()

		return nil
	},
}
