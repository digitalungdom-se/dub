package digitalungdom

import (
	"fmt"

	"github.com/digitalungdom-se/dub/pkg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Status = pkg.Command{
	Name:        "status",
	Description: "Kolla till statusen för Digital Ungdom",
	Aliases:     []string{},
	Group:       "digitalungdom",
	Usage:       "status",
	Example:     "status",
	ServerOnly:  false,
	AdminOnly:   false,

	Execute: func(ctx *pkg.Context) error {
		kelvin, _ := primitive.ObjectIDFromHex("5c62e491a754cd4c7c9c7fa7")
		douglas, _ := primitive.ObjectIDFromHex("5c6b29066157e7938c63aec9")
		simon, _ := primitive.ObjectIDFromHex("5c6bbae9f302e343d917a510")
		charles, _ := primitive.ObjectIDFromHex("5d1dea408fab54e111828f3e")

		filter := bson.M{"_id": bson.M{"$in": bson.A{
			kelvin,
			douglas,
			simon,
			charles,
		}}}

		board, err := ctx.Server.Database.GetUsersByID(filter)
		if err != nil {
			return err
		}

		var boardStatus string
		for _, member := range board {
			name := member.Details.Name
			status := "Inget 🤷"

			if member.Profile.Status != "" {
				status = member.Profile.Status
			}

			boardStatus += fmt.Sprintf("**%v**: %v\n", name, status)

		}

		serverStatus := "**digitalungdom.se**: online\n"
		serverStatus += "**dub**: online"

		var memberStatus string
		digitalungdomCount, err := ctx.Server.Database.GetMemberCount()
		if err != nil {
			return err
		}

		discordCount := ctx.Server.Guild.MemberCount

		memberStatus = fmt.Sprintf("**digitalungdom.se**: %v st\n", digitalungdomCount)
		memberStatus += fmt.Sprintf("**discord**: %v st", discordCount)

		embed := pkg.NewEmbed().
			SetTitle("__**STATUS**__").
			SetColor(4086462).
			AddField("__**STYRELSE**__", boardStatus).
			AddField("__**SERVRAR**__", serverStatus).
			AddField("__**MEDLEMMAR**__", memberStatus).
			MessageEmbed

		_, err = ctx.ReplyEmbed(embed)

		return err
	},
}
