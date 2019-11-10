package digitalungdom

import (
	"context"
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/cbroglie/mustache"
	"github.com/dchest/uniuri"
	"github.com/digitalungdom-se/dub/pkg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	gomail "gopkg.in/gomail.v2"
)

var Verify = pkg.Command{
	Name:        "verify",
	Description: "Koppla ditt discord konto till ditt Digital Ungdom konto",
	Aliases:     []string{"verifiera"},
	Group:       "digitalungdom",
	Usage:       "verify <username|email|token>",
	Example:     "verify kelvin",
	ServerOnly:  false,
	AdminOnly:   false,

	Execute: func(ctx *pkg.Context) error {
		ctx.Delete()

		for _, role := range ctx.Server.Guild.Roles {
			if role.Name == "verified" {
				verifiedID := role.ID
				for _, member := range ctx.Server.Guild.Members {
					if member.User.ID == ctx.Message.Author.ID {
						if pkg.StringInSlice(verifiedID, member.Roles) {
							ctx.DM("Du är redan verifierad")
							return nil
						}
					}
				}
			}
		}

		if len(ctx.Args) == 0 {
			ctx.DM("Du måste ange ditt användarnamn eller email.")
			return nil
		}

		var userType string
		var filter bson.M
		queryOptions := new(options.FindOneAndUpdateOptions)

		var returnDocument options.ReturnDocument = options.After
		queryOptions.ReturnDocument = &returnDocument

		queryOptions.Projection = bson.M{
			"details.email": true,
		}

		if govalidator.IsEmail(ctx.Args[0]) {
			userType = "emailet"
			email, _ := govalidator.NormalizeEmail(ctx.Args[0])
			filter = bson.M{"details.email": email}

			var check bson.M

			findOneOptions := new(options.FindOneOptions)
			findOneOptions.Projection = bson.M{
				"connectedApps.discord": true,
			}

			err := ctx.Server.Database.Collection("users").FindOne(
				context.Background(),
				filter,
				findOneOptions).Decode(&check)

			if err != nil {
				ctx.DM("Inget konto kunde hittas")
				return nil
			}

			if check["connectedApps"].(primitive.M)["discord"] != nil {
				ctx.DM(fmt.Sprintf("Digital Ungdom kontot är redan kopplat till ett discord konto %v", userType))
				return nil
			}

		} else if len(ctx.Args[0]) == 16 {
			filter = bson.M{"discordToken": ctx.Args[0]}
			update := bson.M{
				"$unset": bson.M{
					"discordToken": true,
				}, "$set": bson.M{
					"connectedApps.discord": ctx.Message.Author.ID,
				},
			}
			var user bson.M

			err := ctx.Server.Database.Collection("users").FindOneAndUpdate(
				context.Background(),
				filter,
				update,
				queryOptions).Decode(&user)
			if err != nil {
				ctx.DM("Inget konto kunde hittas")
				return nil
			}

			email := user["details"].(primitive.M)["email"].(string)

			filter = bson.M{"type": "discordVerificationConfirmation"}
			var emailTemplate bson.M
			err = ctx.Server.Database.Collection("emails").FindOne(
				context.Background(),
				filter).Decode(&emailTemplate)
			if err != nil {
				ctx.DM("Kunde inte hitta ett konto kopplat till tokenet")
				return nil
			}

			emailTemplateSTR := emailTemplate["email"].(string)

			data, err := mustache.Render(emailTemplateSTR, map[string]string{"name": ctx.Message.Author.Username})
			if err != nil {
				return err
			}

			m := gomail.NewMessage()
			m.SetHeader("From", m.FormatAddress("noreply@digitalungdom.se", "Digital Ungdom"))
			m.SetHeader("To", email)
			m.SetHeader("Subject", "Koppla ditt Discord konto till Digital Ungdom")
			m.SetBody("text/html", data)

			err = ctx.Server.Dialer.DialAndSend(m)
			if err != nil {
				return err
			}

			var verifiedID string

			for _, role := range ctx.Server.Guild.Roles {
				if role.Name == "verified" {
					verifiedID = role.ID
				}
			}

			ctx.Discord.GuildMemberRoleAdd(ctx.Server.Guild.ID, ctx.Message.Author.ID, verifiedID)

			ctx.DM("Grattis ditt Discord konto är nu kopplat till Digital Ungdom")

			return nil
		} else {
			userType = "användarnamnet"
			filter = bson.M{"details.username": ctx.Args[0]}

			collation := new(options.Collation)
			collation.Locale = "en"
			collation.Strength = 2
			queryOptions.Collation = collation

			findOneOptions := new(options.FindOneOptions)
			findOneOptions.Projection = bson.M{
				"connectedApps.discord": true,
			}
			var check bson.M

			err := ctx.Server.Database.Collection("users").FindOne(
				context.Background(),
				filter,
				findOneOptions).Decode(&check)

			if err != nil {
				err = ctx.DM("Inget konto kunde hittas")
				return err
			}

			if check["connectedApps"].(primitive.M)["discord"] != nil {
				ctx.DM(fmt.Sprintf("Digital Ungdom kontot är redan kopplat till ett discord konto %v", userType))
				return nil
			}
		}

		var user bson.M

		token := uniuri.NewLen(16)

		update := bson.M{
			"$set": bson.M{
				"discordToken": token,
			},
		}

		err := ctx.Server.Database.Collection("users").FindOneAndUpdate(
			context.Background(),
			filter,
			update,
			queryOptions).Decode(&user)
		if err != nil {
			ctx.DM("Inget konto kunde hittas")
			return nil
		}

		if user["details"].(primitive.M)["email"].(string) == "" {
			ctx.DM(fmt.Sprintf("Ingen användare kunde hittas med det %v", userType))
			return nil
		}

		email := user["details"].(primitive.M)["email"].(string)

		filter = bson.M{"type": "discordVerification"}
		var emailTemplate bson.M
		err = ctx.Server.Database.Collection("emails").FindOne(
			context.Background(),
			filter).Decode(&emailTemplate)
		if err != nil {
			return err
		}

		emailTemplateSTR := emailTemplate["email"].(string)

		data, err := mustache.Render(emailTemplateSTR, map[string]string{"token": token})
		if err != nil {
			return err
		}

		m := gomail.NewMessage()
		m.SetHeader("From", m.FormatAddress("noreply@digitalungdom.se", "Digital Ungdom"))
		m.SetHeader("To", email)
		m.SetHeader("Subject", "Koppla ditt Discord konto till Digital Ungdom")
		m.SetBody("text/html", data)

		err = ctx.Server.Dialer.DialAndSend(m)
		if err != nil {
			return err
		}

		ctx.DM("Ett email har nu skickats till dig, följ instruktionerna där för att komma vidare.")

		return nil
	},
}
