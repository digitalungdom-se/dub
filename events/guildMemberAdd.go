package events

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/digitalungdom-se/dub/internal"
	"github.com/digitalungdom-se/dub/pkg"
)

func GuildMemberAddHandler(server *internal.Server) func(*discordgo.Session, *discordgo.GuildMemberAdd) {
	return func(discord *discordgo.Session, member *discordgo.GuildMemberAdd) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("%v | RECOVERED ANNA FROM %v\n", time.Now().Format("2006-01-02 15:04:05"), r)
			}
		}()

		gifbuff, err := pkg.NameToGif(member.User.Username, member.User.AvatarURL("128"))
		if err != nil {
			log.Print("Error creating user gif:", err)
			return
		}
		reader := bytes.NewReader((*gifbuff).Bytes())

		_, err = discord.ChannelFileSendWithMessage(server.Channels.General.ID,
			fmt.Sprintf("Välkommen till Digital Ungdoms Discord server <@%v>!", member.User.ID),
			"welcome.gif",
			reader)
		if err != nil {
			log.Print("Error sending welcome message:", err)
			return
		}

		var privateDM *discordgo.Channel

		privateDM, err = discord.UserChannelCreate(member.User.ID)
		if err != nil {
			log.Print("Error creating private channel:", err)
			return
		}

		content := fmt.Sprintf("Hej **%v** och välkommen till *Digital Ungdoms* Discord server."+
			" Jag är boten som hjälper till med kanalen. För att se alla mina funktioner skriv `$help` till mig.\n\n"+
			" __**Du måste godkänna reglerna genom att trycka på den gröna knappen längst ner i `#regler` i Digital Ungdoms Discord kanal innan du får börja skriva och delta i samtal.**__\n\n"+
			"Om du inte redan har ett *Digital Ungdom* konto så rekommenderar jag starkt att du skaffar ett."+
			" Som medlem kan du bland annat skriva på vårt forum (https://digitalungdom.se/agora)."+
			" Du kan enkelt bli medlem genom följande länk: https://digitalungdom.se/bli-medlem\n\n"+
			"Om du redan är medlem så kan du koppla ditt *Digital Ungdom* konto till ditt Discord konto"+
			" genom att skriva `$verify` och sedan ditt användarnamn eller epost. Till exempel `$verify username`\n\n"+
			"**Vi synns där inne!**",
			member.User.Username)

		_, err = discord.ChannelMessageSend(privateDM.ID, content)
		if err != nil {
			log.Print("Error sending private welcome message:", err)
			return
		}
	}
}
