package digitalungdom

import (
	"strings"

	"github.com/digitalungdom-se/dub/internal"
	"github.com/digitalungdom-se/dub/pkg"
	"go.mongodb.org/mongo-driver/bson"
)

var Bug = internal.Command{
	Name:        "bug",
	Description: "Anmäl ett bugg",
	Aliases:     []string{"bugg"},
	Group:       "digitalungdom",
	Usage:       "bug <message>",
	Example:     "bug boten spelar inte musik",
	ServerOnly:  false,
	AdminOnly:   false,

	Execute: func(ctx *pkg.Context, server *internal.Server) error {
		if len(ctx.Args) < 3 {
			_, err := ctx.Reply("Du måste ge en anledning på minst tre ord.")
			return err
		}

		bug := bson.M{
			"type":    "bug",
			"where":   "discord",
			"message": strings.Join(ctx.Args, " "),
			"author":  ctx.Message.Author.ID}

		err := server.Database.InsertNotification(bug)
		if err != nil {
			return err
		}

		_, err = ctx.Reply("Du har nu anmält denna bugg till Digital Ungdom. Tack för din medverkan!")
		if err != nil {
			return err
		}

		return nil
	},
}
