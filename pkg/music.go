package pkg

import (
	"fmt"
	"io"
	"net/url"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/rylio/ytdl"
)

var musicEmojis = map[string]string{
	"stop": "❎",
	"pr":   "⏯",
	"skip": "⏭",
}

type (
	Controller struct {
		channel        *discordgo.Channel
		discord        *discordgo.Session
		guild          *discordgo.Guild
		previousStatus *discordgo.UpdateStatusData

		embed           *discordgo.MessageEmbed
		Message         *discordgo.Message
		voiceConnection *discordgo.VoiceConnection

		encoder  *dca.EncodeSession
		streamer *dca.StreamingSession

		musicQueue []Song
		playing    bool
		skip       bool
		volume     int
	}

	Song struct {
		downloadURL string
		metadata    *Metadata
		requesterID string
		url         string
	}

	Metadata struct {
		author         string
		authorImageURL string
		authorURL      string
		duration       string
		thumbnailURL   string
		title          string
	}
)

func NewController(discord *discordgo.Session, channel *discordgo.Channel, previousStatus *discordgo.UpdateStatusData, reactionListener *ReactionListener) (Controller, error) {
	var err error
	var controller Controller

	controller.discord = discord
	controller.channel = channel
	controller.previousStatus = previousStatus

	controller.NewControllerMessage(reactionListener)

	return controller, err
}

func (controller *Controller) Play() error {
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Volume = 128
	options.Application = "lowdelay"

	song := controller.musicQueue[0]

	var err error
	controller.encoder, err = dca.EncodeFile(song.downloadURL, options)
	if err != nil {
		return err
	}

	discordStatus := new(discordgo.UpdateStatusData)
	discordStatus.Game = new(discordgo.Game)
	discordStatus.Status = song.metadata.title
	discordStatus.Game.Name = song.metadata.title
	discordStatus.Game.Type = discordgo.GameTypeListening

	controller.discord.UpdateStatusComplex(*discordStatus)

	done := make(chan error)

	controller.streamer = dca.NewStream(controller.encoder, controller.voiceConnection, done)
	err = <-done

	if err == io.EOF {
		if !controller.skip {
			controller.end()
		}
	}

	return nil
}

func (controller *Controller) end() error {
	if controller.encoder != nil {
		controller.encoder.Cleanup()
	}

	if len(controller.musicQueue) <= 1 {
		controller.musicQueue = []Song{}
	} else {
		controller.musicQueue = controller.musicQueue[1:]
	}

	controller.updateControllerMessage()

	if len(controller.musicQueue) == 0 {
		controller.discord.UpdateStatusComplex(*controller.previousStatus)

		if controller.voiceConnection != nil {
			err := controller.voiceConnection.Disconnect()
			controller.voiceConnection = nil
			if err != nil {
				return err
			}
		}

		return nil
	}

	err := controller.Play()

	return err
}

func (controller *Controller) AddToQueue(context *Context) error {
	songURL := context.Args[0]

	u, err := url.ParseRequestURI(songURL)
	if err != nil || (u.Hostname() != "www.youtube.com" && u.Hostname() != "youtu.be" && u.Hostname() != "m.youtube.com") {
		context.Reply("Du måste skicka en youtube-länk.")
		return nil
	}

	videoInfo, err := ytdl.GetVideoInfo(songURL)
	if err != nil {
		return err
	}

	metadata := new(Metadata)
	metadata.title = videoInfo.Title

	info, err := GetYoutubeInfo(songURL)
	if err != nil {
		return err
	}

	metadata.thumbnailURL = info.ThumbnailURL
	metadata.author = info.Author.Name
	metadata.authorURL = info.Author.URL
	metadata.authorImageURL = info.Author.AvatarURL

	format := videoInfo.Formats.Extremes(ytdl.FormatAudioBitrateKey, true)[0]
	downloadURL, err := videoInfo.GetDownloadURL(format)
	if err != nil {
		return err
	}

	song := Song{url: songURL, downloadURL: downloadURL.String(), requesterID: context.Message.Author.ID}
	song.metadata = metadata

	controller.musicQueue = append(controller.musicQueue, song)

	if len(controller.musicQueue) <= 6 {
		controller.updateControllerMessage()
	}

	if controller.voiceConnection == nil {
		guild, err := controller.discord.Guild(controller.channel.GuildID)
		if err != nil {
			return err
		}

		for _, vs := range guild.VoiceStates {
			if vs.UserID == context.Message.Author.ID {
				controller.voiceConnection, err = controller.discord.ChannelVoiceJoin(
					vs.GuildID, vs.ChannelID, false, true)
				if err != nil {
					return err
				}

				break
			}
		}

		err = controller.Play()
		if err != nil {
			return err
		}
	}

	return nil
}

func (controller *Controller) Skip() {
	controller.skip = true
	controller.end()
	controller.skip = false
}

func (controller *Controller) Stop() {
	controller.musicQueue = []Song{}

	controller.end()
}

func (controller *Controller) PauseResume() {
	if controller.streamer != nil {
		if controller.streamer.Paused() {
			controller.streamer.SetPaused(false)
			return
		}

		controller.streamer.SetPaused(true)
	}
}

func (controller *Controller) NewControllerMessage(reactionListener *ReactionListener) error {
	for {
		messages, err := controller.discord.ChannelMessages(controller.channel.ID, 100, "", "", "")
		if err != nil {
			return err
		}
		var messagesID []string

		for _, message := range messages {
			messagesID = append(messagesID, message.ID)
		}

		err = controller.discord.ChannelMessagesBulkDelete(controller.channel.ID, messagesID)
		if err != nil {
			return err
		}

		if len(messages) <= 100 {
			break
		}
	}

	var embed *discordgo.MessageEmbed

	if len(controller.musicQueue) == 0 {
		embed = NewEmbed().
			SetTitle("Inget spelas").
			SetColor(4086462).
			MessageEmbed
	} else {
		embed = controller.newControllerEmbed()
	}

	reactionator := NewReactionator(controller.channel.ID,
		controller.discord, reactionListener, false, true, ReactionatorTypeController, nil)

	reactionator.AddDefaultPage("", embed)

	reactionator.Add(musicEmojis["stop"], controller.Stop)
	reactionator.Add(musicEmojis["pr"], controller.PauseResume)
	reactionator.Add(musicEmojis["skip"], controller.Skip)

	err := reactionator.Initiate()

	controller.Message = reactionator.Message

	return err
}

func (controller *Controller) newControllerEmbed() *discordgo.MessageEmbed {
	description := fmt.Sprintf("Spelas nu på begäran av <@%v>.\n\n",
		controller.musicQueue[0].requesterID)

	if len(controller.musicQueue) > 1 {
		description += fmt.Sprintf("*Visar de %v första låtar i kön.*", len(controller.musicQueue)-1)
	}

	embed := NewEmbed().
		SetTitle(controller.musicQueue[0].metadata.title).
		SetDescription(description).
		SetURL(controller.musicQueue[0].url).
		SetThumbnail(controller.musicQueue[0].metadata.thumbnailURL).
		SetAuthor(controller.musicQueue[0].metadata.author, controller.musicQueue[0].metadata.authorImageURL,
			controller.musicQueue[0].metadata.authorURL).
		SetColor(4086462)

	if len(controller.musicQueue) > 1 {
		musicQueue := controller.musicQueue[1:len(controller.musicQueue)]
		queue := ""
		for index, song := range musicQueue {
			queue += fmt.Sprintf("**%v.** %v | %v\n", index, song.metadata.author, song.metadata.title)
		}
		embed.AddField("__**Kö**__", queue)
	}

	return embed.MessageEmbed
}

func (controller *Controller) updateControllerMessage() error {
	var embed *discordgo.MessageEmbed

	if len(controller.musicQueue) > 0 {
		embed = controller.newControllerEmbed()
	} else {
		embed = NewEmbed().
			SetTitle("Inget spelas").
			SetColor(4086462).
			MessageEmbed
	}

	_, err := controller.discord.ChannelMessageEditEmbed(controller.channel.ID,
		controller.Message.ID, embed)
	if err != nil {
		return err
	}

	return nil
}
