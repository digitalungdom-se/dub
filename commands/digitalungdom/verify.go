package digitalungdom

import (
	"github.com/asaskevich/govalidator"
	"github.com/digitalungdom-se/dub/internal"
	"github.com/digitalungdom-se/dub/pkg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Verify = internal.Command{
	Name:        "verify",
	Description: "Koppla ditt discord konto till ditt Digital Ungdom konto",
	Aliases:     []string{"verifiera"},
	Group:       "digitalungdom",
	Usage:       "verify <username|email|token>",
	Example:     "verify kelvin",
	ServerOnly:  false,
	AdminOnly:   false,

	Execute: func(ctx *pkg.Context, server *internal.Server) error {
		msg, err := ctx.Reply("**BEARBETAR...**")
		if err != nil {
			return err
		}

		if !ctx.IsDM() {
			err := ctx.Delete()
			if err != nil {
				return err
			}
		}

		for _, role := range server.Guild.Roles {
			if role.Name == "verified" {
				verifiedID := role.ID
				for _, member := range server.Guild.Members {
					if member.User.ID == ctx.Message.Author.ID {
						if pkg.StringInSlice(verifiedID, member.Roles) {
							_, err = ctx.Discord.ChannelMessageEdit(msg.ChannelID, msg.ID, "Du är redan verifierad.")
							return err
						}
					}
				}
			}
		}

		if len(ctx.Args) == 0 {
			_, err = ctx.Discord.ChannelMessageEdit(msg.ChannelID, msg.ID, "Du måste ange ditt användarnamn eller email.")
			return err
		}

		queryOptions := new(options.FindOneAndUpdateOptions)

		var returnDocument options.ReturnDocument = options.After
		queryOptions.ReturnDocument = &returnDocument

		queryOptions.Projection = bson.M{
			"details.email": true,
		}

		if len(ctx.Args[0]) == 16 {
			email, err := server.Database.ConnectDUAccount(ctx.Args[0], ctx.Message.Author.ID)
			if err != nil || email == "" {
				_, err = ctx.Discord.ChannelMessageEdit(msg.ChannelID, msg.ID, "Inget konto kunde hittas")

				return err
			}

			err = server.Mailer.SendVerifyDiscordConfirmation(email, ctx.Message.Author.Username)
			if err != nil {
				return err
			}

			var verifiedID string

			for _, role := range server.Guild.Roles {
				if role.Name == "verified" {
					verifiedID = role.ID
				}
			}

			err = ctx.Discord.GuildMemberRoleAdd(server.Guild.ID, ctx.Message.Author.ID, verifiedID)
			if err != nil {
				return err
			}

			_, err = ctx.Discord.ChannelMessageEdit(msg.ChannelID, msg.ID, "Grattis ditt Discord konto är nu kopplat till Digital Ungdom!")

			return err
		}

		var user internal.User

		if govalidator.IsEmail(ctx.Args[0]) {
			user, err = server.Database.GetUserByEmail(ctx.Args[0])

			if err != nil {
				_, err = ctx.Discord.ChannelMessageEdit(msg.ChannelID, msg.ID, "Inget konto kunde hittas")
				return err
			}

			if user.ConnectedApps.Discord != "" {
				_, err = ctx.Discord.ChannelMessageEdit(msg.ChannelID, msg.ID, "Ditt Digital Ungdom kontot är redan kopplat till ett Discord konto.")
				return err
			}

		} else {
			user, err = server.Database.GetUserByUsername(ctx.Args[0])
			if err != nil {
				_, err = ctx.Discord.ChannelMessageEdit(msg.ChannelID, msg.ID, "Inget konto kunde hittas")
				return err
			}

			if user.ConnectedApps.Discord != "" {
				_, err = ctx.Discord.ChannelMessageEdit(msg.ChannelID, msg.ID, "Ditt Digital Ungdom kontot är redan kopplat till ett Discord konto.")
				return err
			}
		}

		token, err := server.Database.AddDiscordVerificationToken(user.ID.Hex())
		if err != nil {
			return err
		}

		err = server.Mailer.SendVerifyDiscord(user.Details.Email, token)
		if err != nil {
			return err
		}

		_, err = ctx.Discord.ChannelMessageEdit(msg.ChannelID, msg.ID, "Ett email har nu skickats till dig, följ instruktionerna där för att komma vidare.")
		if err != nil {
			return err
		}

		return nil
	},
}
