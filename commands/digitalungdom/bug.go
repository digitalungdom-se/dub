package digitalungdom

import (
	"context"
	"strings"

	"github.com/digitalungdom-se/dub/pkg"
	"go.mongodb.org/mongo-driver/bson"
)

var Bug = pkg.Command{
	Name:        "bug",
	Description: "Anmäl ett bugg",
	Aliases:     []string{"bugg"},
	Group:       "digitalungdom",
	Usage:       "bug <message>",
	Example:     "bug boten spelar inte musik",
	ServerOnly:  false,
	AdminOnly:   false,

	Execute: func(ctx *pkg.Context) error {
		if len(ctx.Args) < 3 {
			ctx.Reply("Du måste ge en anledning")
			return nil
		}

		bug := bson.M{
			"type":    "bug",
			"where":   "discord",
			"message": strings.Join(ctx.Args, " "),
			"author":  ctx.Message.Author.ID}

		ctx.Server.Database.Collection("notifications").InsertOne(context.TODO(), bug)

		ctx.Reply("Du har nu anmält denna bugg till Digital Ungdom. Tack för din medverkan!")

		return nil
	},
}
