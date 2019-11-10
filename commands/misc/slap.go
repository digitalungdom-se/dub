package misc

import (
	"fmt"

	"github.com/digitalungdom-se/dub/pkg"
)

var Slap = pkg.Command{
	Name:        "slap",
	Description: "Smiska någon som har varit stygg",
	Aliases:     []string{"smisk"},
	Group:       "misc",
	Usage:       "smisk @<user>",
	Example:     "smisk @Ippyson#6200 ",
	ServerOnly:  true,
	AdminOnly:   false,

	Execute: func(context *pkg.Context) error {
		if len(context.Message.Mentions) == 0 {
			context.Reply("Vem har varit stygg?")
			return nil
		}

		slapped, err := context.GetMentions()

		if err != nil {
			return err
		}

		context.Reply(fmt.Sprintf("Du har varit riktigt stygg <@%v>", slapped[0]))

		return nil
	},
}
