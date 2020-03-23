package digitalungdom

import (
	"strings"

	"github.com/digitalungdom-se/dub/pkg"
	"go.mongodb.org/mongo-driver/bson"
)

var Idea = pkg.Command{
	Name:        "idea",
	Description: "Föreslå något till Digital Ungdom",
	Aliases:     []string{"förslag"},
	Group:       "digitalungdom",
	Usage:       "idea <idea>",
	Example:     "idea skaffa programerings tutorials",
	ServerOnly:  false,
	AdminOnly:   false,

	Execute: func(ctx *pkg.Context) error {
		if len(ctx.Args) < 3 {
			_, err := ctx.Reply("Du måste ge ett förslag på minst tre ord.")
			return err
		}

		idea := bson.M{
			"type":    "idea",
			"where":   "discord",
			"message": strings.Join(ctx.Args, " "),
			"author":  ctx.Message.Author.ID}

		err := ctx.Server.Database.InsertNotification(idea)
		if err != nil {
			return err
		}

		_, err = ctx.Reply("Du har nu skickat in ditt förslag. Tack för din medverkan!")

		return err
	},
}
