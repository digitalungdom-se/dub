package digitalungdom

import (
	"context"
	"fmt"

	"github.com/digitalungdom-se/dub/pkg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

		findOptions := options.Find()
		findOptions.Projection = bson.M{
			"details.name":   true,
			"profile.status": true,
		}

		cur, err := ctx.Server.Database.Collection("users").Find(context.TODO(), filter, findOptions)
		if err != nil {
			return err
		}

		var boardStatus string
		for cur.Next(context.TODO()) {
			var user bson.M

			err := cur.Decode(&user)
			if err != nil {
				return err
			}

			name := user["details"].(primitive.M)["name"].(string)
			status := "Inget 🤷"

			if user["profile"].(primitive.M)["status"] != nil {
				status = user["profile"].(primitive.M)["status"].(string)
			}

			boardStatus += fmt.Sprintf("**%v**: %v\n", name, status)
		}
		cur.Close(context.TODO())

		serverStatus := "**digitalungdom.se**: online\n"
		serverStatus += "**dub**: online"

		var memberStatus string
		digitalungdomCount, err := ctx.Server.Database.Collection("users").CountDocuments(context.TODO(), bson.M{})
		if err := cur.Err(); err != nil {
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

		ctx.ReplyEmbed(embed)

		return nil
	},
}
