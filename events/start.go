package events

import (
	"log"

	"golang.org/x/sync/errgroup"

	"github.com/bwmarrin/discordgo"
	"github.com/digitalungdom-se/dub/internal"
)

func StartHandler(server *internal.Server) func(*discordgo.Session, *discordgo.Ready) {
	return func(discord *discordgo.Session, ready *discordgo.Ready) {
		err := server.Init()
		if err != nil {
			log.Fatal(err)
		}

		errGroup := errgroup.Group{}
		errGroup.Go(server.InitController)
		errGroup.Go(server.InitRules)

		err = errGroup.Wait()
		if err != nil {
			log.Fatal(err)
		}
	}
}
