package misc

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/digitalungdom-se/dub/pkg"
)

var groupReactions = map[string]string{
	"info":          "ℹ",
	"digitalungdom": "🖥",
	"music":         "🎵",
	"misc":          "🛠",
	"admin":         "🚨",
	"close":         "🔥",
}

var groupReactionOrder = []string{"digitalungdom", "music", "misc"}

var Help = pkg.Command{
	Name:        "help",
	Description: "Listar alla tillgängliga kommandon",
	Aliases:     []string{"commands", "command", "hjälp", "kommando", "kommandon"},
	Group:       "misc",
	Usage:       "help <command>",
	Example:     "help",
	ServerOnly:  false,
	AdminOnly:   false,

	Execute: func(context *pkg.Context) error {
		context.Delete()
		if len(context.Args) != 0 {
			command, found := context.Server.CommandHandler.GetCommand(context.Args[0])

			if !found {
				return errors.New("Kommandot kunde inte hittas")
			}

			embed := pkg.NewEmbed().
				SetTitle(fmt.Sprintf("**%v**", command.Name)).
				SetDescription(fmt.Sprintf("*%v*", command.Description)).
				AddField("ANVÄNDNING", fmt.Sprintf(">`%v`", command.Usage)).
				AddField("EXEMPEL", fmt.Sprintf(">`%v`", command.Example)).
				SetColor(4086462)

			if len(command.Aliases) > 0 {
				embed.AddField("ALIAS", fmt.Sprintf("`%v`", strings.Join(command.Aliases[:], ", ")))
			}

			_, err := context.ReplyEmbed(embed.MessageEmbed)

			if err != nil {
				return err
			}

			return nil
		}

		commands := context.Server.CommandHandler.GetCommands("")
		groups := make(map[string][]*pkg.Command)

		for _, command := range commands {
			groups[command.Group] = append(groups[command.Group], command)
		}

		embeds := make(map[string]*discordgo.MessageEmbed)

		description := "__Tryck knapparna längst ned för att byta sida__.\n" +
			"Du kan få mer information om ett kommando genom att köra `>help <command>`.\n\n" +
			":information_source: **--** Denna sida\n" +
			":desktop: **--** Digital Ungdom kommandon\n" +
			":musical_note:  **--** Musik kommandon\n" +
			":tools: **--** Misc kommandon\n"

		admin := false

		for _, member := range context.Server.Guild.Members {
			if member.User.ID == context.Message.Author.ID {
				if pkg.StringInSlice(context.Server.Roles.Admin.ID, member.Roles) {
					description += "🚨 **--** Admin kommandon\n"
					admin = true
				}
			}
		}

		description += ":fire:  **--** Stäng hjälp sida\n"

		embeds["info"] = pkg.NewEmbed().
			SetTitle("**HJÄLP SIDA**").
			SetDescription(description).
			SetColor(4086462).
			MessageEmbed

		for group, groupCommands := range groups {
			embed := pkg.NewEmbed().
				SetTitle(fmt.Sprintf("**%v**", group)).
				SetDescription(fmt.Sprintf("Hjälp sida för kommandon i *%v* gruppen", group))

			for _, command := range groupCommands {
				embed.AddField(fmt.Sprintf("__**%v**__", command.Name),
					fmt.Sprintf("%v\n>`%v`", command.Description, command.Usage))
			}

			embed = embed.SetColor(4086462)

			embeds[group] = embed.MessageEmbed
		}

		privateDM, err := context.Discord.UserChannelCreate(context.Message.Author.ID)
		if err != nil {
			return err
		}

		reactionator := pkg.NewReactionator(privateDM.ID, context.Discord, context.Server.ReactionListener,
			true, pkg.ReactionatorTypeHelp, context.Message.Author)

		err = reactionator.AddDefaultPage(groupReactions["info"], embeds["info"])
		if err != nil {
			return err
		}

		for _, group := range groupReactionOrder {
			err = reactionator.Add(groupReactions[group], embeds[group])
			if err != nil {
				return err
			}
		}

		if admin {
			err = reactionator.Add(groupReactions["admin"], embeds["admin"])
			if err != nil {
				return err
			}
		}

		err = reactionator.CloseButton()
		if err != nil {
			return err
		}

		reactionator.CloseAfter(3 * time.Minute)

		if activeReactionators, ok := context.Server.ReactionListener.Users[context.Message.Author.ID]; ok {
			if activeReactionators.Help != nil {
				activeReactionators.Help.Close()
			}
		}

		err = reactionator.Initiate()
		if err != nil {
			return err
		}

		if context.Message.ChannelID != privateDM.ID {
			_, err = context.Reply("Ett direkt meddelande har skickats till dig " +
				"med alla kommandon. Du finner dem längst upp till höger")

			if err != nil {
				return err
			}
		}

		return nil
	},
}
