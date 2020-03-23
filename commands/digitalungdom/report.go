package digitalungdom

import (
	"strings"

	"github.com/digitalungdom-se/dub/pkg"
	"go.mongodb.org/mongo-driver/bson"
)

var Report = pkg.Command{
	Name:        "report",
	Description: "Anmäl en användare",
	Aliases:     []string{"anmäl"},
	Group:       "digitalungdom",
	Usage:       "report @<user> <reason>",
	Example:     "report @Ippyson#6200 han är taskig",
	ServerOnly:  false,
	AdminOnly:   false,

	Execute: func(ctx *pkg.Context) error {
		if !ctx.IsDM() {
			err := ctx.Delete()
			if err != nil {
				return err
			}
		}

		if len(ctx.Message.Mentions) == 0 {
			ctx.Reply("Du måste anmäla någon")
			return nil
		}

		if len(ctx.Args[1:]) < 3 {
			ctx.Reply("Du måste ge en anledning")
			return nil
		}

		report := bson.M{
			"type":     "report",
			"where":    "discord",
			"message":  strings.Join(ctx.Args[1:], " "),
			"author":   ctx.Message.Author.ID,
			"reported": ctx.Message.Mentions[0].ID}

		err := ctx.Server.Database.InsertNotification(report)
		if err != nil {
			return err
		}

		err = ctx.DM("Du har nu anmält denna person till Digital Ungdom. Tack för din medverkan!")

		return err
	},
}
