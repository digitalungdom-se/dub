package digitalungdom

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/digitalungdom-se/dub/pkg"
)

var Whois = pkg.Command{
	Name:        "whois",
	Description: "Få Digital Ungdom kontot av en användare",
	Aliases:     []string{"verifiera"},
	Group:       "digitalungdom",
	Usage:       "whois @<user>",
	Example:     "whois @kelszo#6200",
	ServerOnly:  true,
	AdminOnly:   false,

	Execute: func(ctx *pkg.Context) error {
		mentions, err := ctx.GetMentions()
		if err != nil || len(mentions) == 0 {
			ctx.Reply("Du måste nämna en användare.")
			return nil
		}

		user, err := ctx.Server.Database.GetUserByDiscordID(mentions[0])
		if err != nil {
			ctx.Reply(fmt.Sprintf("Inget konto kunde hittas kopplat till den användare. %v borde verkligen koppla sitt konta genom `$verify <digital_ungdom_användarnamn>`!", ctx.Message.Mentions[0].Mention()))
			return nil
		}

		if user.Details.Username == "" {
			ctx.Reply(fmt.Sprintf("Inget konto kunde hittas kopplat till den användare. %v borde verkligen koppla sitt konta genom `$verify <digital_ungdom_användarnamn>`!", ctx.Message.Mentions[0].Mention()))
			return nil
		}

		firstName := strings.Split(user.Details.Name, " ")[0]

		userEmbed := pkg.NewEmbed()
		userEmbed.SetTitle(user.Details.Username)
		userEmbed.SetDescription(fmt.Sprintf("%v har varit medlem sedan %v. %v har lagt upp %v posts på Agora och fått %v stjärnor.",
			user.Details.Name, user.ID.Timestamp().Format("02/01/2006"), firstName, user.Agora.Score.Posts, user.Agora.Score.Stars))
		userEmbed.SetURL(fmt.Sprintf("https://digitalungdom.se/@%v", user.Details.Username))
		if user.Profile.Status != "" {
			userEmbed.AddField("__**STATUS**__", user.Profile.Status)
		}

		if user.Profile.Bio != "" {
			userEmbed.AddField("__**BIO**__", user.Profile.Bio)
		}

		if user.Profile.URL != "" {
			userEmbed.AddField("__**URL**__", user.Profile.URL)
		}

		var colour int64
		colour, err = strconv.ParseInt(user.Profile.Colour[1:], 16, 64)
		if err != nil {
			colour = 4086462
		}
		userEmbed.SetColor(int(colour))

		var resp *http.Response
		resp, err = http.Get(fmt.Sprintf("https://digitalungdom.se/api/agora/get/profile_picture?id=%v&size=128", user.ID.Hex()))
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		var profilePictureResult struct {
			ImageType struct {
				Ext  string
				Mime string
			}
			Image string
		}

		json.NewDecoder(resp.Body).Decode(&profilePictureResult)

		var filesToSend []*discordgo.File
		if profilePictureResult.Image != "" {
			stringReader := strings.NewReader(profilePictureResult.Image)

			reader := base64.NewDecoder(base64.StdEncoding, stringReader)

			discordFile := &discordgo.File{
				Name:        "pp.png",
				ContentType: profilePictureResult.ImageType.Mime,
				Reader:      reader,
			}

			filesToSend = []*discordgo.File{discordFile}

			userEmbed.SetThumbnail("attachment://pp.png")
		}

		message := &discordgo.MessageSend{
			Embed: userEmbed.MessageEmbed,
			Files: filesToSend,
		}

		_, err = ctx.ReplyComplex(message)

		return err
	},
}
