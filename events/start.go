package events

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/digitalungdom-se/dub/pkg"
)

func StartHandler(server *pkg.Server) func(*discordgo.Session, *discordgo.Ready) {
	return func(discord *discordgo.Session, ready *discordgo.Ready) {
		err := server.Init()
		if err != nil {
			log.Fatal(err)
		}

		for {
			messages, err := server.Discord.ChannelMessages(server.Channels.Regler.ID, 100, "", "", "")
			if err != nil {
				log.Println(err)
			}
			var messagesID []string

			for _, message := range messages {
				messagesID = append(messagesID, message.ID)
			}

			err = server.Discord.ChannelMessagesBulkDelete(server.Channels.Regler.ID, messagesID)
			if err != nil {
				log.Println(err)
			}

			if len(messages) <= 100 {
				break
			}
		}

		rulesEmbed := pkg.NewEmbed()
		rulesEmbed.SetTitle("Regler för Discord kanalen.")
		rulesEmbed.SetColor(16711680)
		rulesEmbed.AddField("__**Du måste vara medlem i Digital Ungdom**__", "Bli medlem här: https://digitalungdom.se/bli-medlem.\n")
		rulesEmbed.AddField("__**Använd ditt förnamn**__", "Vi vill gärna skapa en personlig miljö. Använd `/nick`-kommandot. Du kan behålla ditt alias genom ex. `/nick Nikolaus (Jultomten)`.\n")
		rulesEmbed.AddField("__**Gyllene regeln**__", "Gör mot andra som du själv vill bli behandlad.\n")
		rulesEmbed.AddField("__**Inget stötande beteende.**__", "- Inga rasistiska, sexistiska, eller aggressiva kommentarer.\n- Inget spam.\n- Undvik gärna svordomar.")

		thinkAbout := pkg.NewEmbed()
		thinkAbout.SetTitle("Tänk på...")
		thinkAbout.SetColor(589568)
		thinkAbout.AddField("__**Undvik elitism**__", " Var upplyftande gentemot nybörjare! Alla vill lära sig mer och det blir roligare om vi bidrar till varandras utveckling och framsteg. Försök att inte överväldiga nybörjare med en onödigt teknisk jargong.\n")
		thinkAbout.AddField("__**Använd rätt kanal**__", " Vi har försökt strukturera kanalerna efter ämne, men om du inte känner att något passar kan du antingen skriva till oss eller skriva i #general.\n")
		thinkAbout.AddField("__**Var gärna tydlig med din nivå**__", "Som nybörjare kan en få mycket hjälp genom att förtydliga att en är ny till vissa områden. Som erfaren kan en känna sig osäker om ens korrespondent egentligen kan mycket mer än en tror!\n")

		reactionator := pkg.NewReactionator(server.Channels.Regler.ID, discord, server.ReactionListener, false, false, pkg.ReactionatorTypeController, nil)
		err = reactionator.AddDefaultPage("", "__**TRYCK PÅ DEN GRÖNA KNAPPEN UNDER MEDDELANDET FÖR ATT GODKÄNNA REGLERNA OCH DÄRMED KUNNA SKRIVA OCH DELTA I SAMTAL**__")
		if err != nil {
			log.Fatal(err)
		}

		err = reactionator.Add("✅", func(message *discordgo.MessageReaction) {
			var medlemRoleID string

			for _, role := range server.Guild.Roles {
				if role.Name == "accepterat_reglerna" {
					medlemRoleID = role.ID
				}
			}

			discord.GuildMemberRoleAdd(server.Guild.ID, message.UserID, medlemRoleID)
		})
		if err != nil {
			log.Fatal(err)
		}

		server.Discord.ChannelMessageSendEmbed(server.Channels.Regler.ID, rulesEmbed.MessageEmbed)
		server.Discord.ChannelMessageSendEmbed(server.Channels.Regler.ID, thinkAbout.MessageEmbed)
		err = reactionator.Initiate()
		if err != nil {
			log.Fatal(err)
		}
	}
}
