package pkg

import (
	"errors"
	"log"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

type Context struct {
	Discord     *discordgo.Session
	Server      *Server
	TextChannel *discordgo.Channel
	Message     *discordgo.MessageCreate
	Args        []string
}

func NewContext(discord *discordgo.Session, server *Server, textChannel *discordgo.Channel,
	message *discordgo.MessageCreate, args []string) *Context {
	context := new(Context)

	context.Discord = discord
	context.Server = server
	context.TextChannel = textChannel
	context.Message = message
	context.Args = args

	return context
}

func (context *Context) Delete() error {
	err := context.Discord.ChannelMessageDelete(context.Message.ChannelID, context.Message.ID)
	return err
}

func (context *Context) Reply(content string) (*discordgo.Message, error) {
	msg, err := context.Discord.ChannelMessageSend(context.TextChannel.ID, content)

	return msg, err
}

func (context *Context) ReplyEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	msg, err := context.Discord.ChannelMessageSendEmbed(context.TextChannel.ID, embed)

	return msg, err
}

func (context *Context) GetMentions() ([]string, error) {
	r, err := regexp.Compile(`<@!(\d{18})>`)

	if err != nil {
		log.Println("Error in regex,", err)
		return []string{}, err
	}
	matches := r.FindAllStringSubmatch(context.Message.Content, -1)
	mentions := make([]string, len(matches)-1)

	for _, mention := range matches {
		mentions = append(mentions, mention[1])
	}

	return mentions, nil
}

func (context *Context) GetVoiceChannel() (*discordgo.VoiceState, error) {
	for _, guild := range context.Discord.State.Guilds {
		for _, vs := range guild.VoiceStates {
			if vs.UserID == context.Message.Author.ID {
				return vs, nil
			}
		}
	}
	return nil, errors.New("Could not find user's voice state")
}

func (context *Context) JoinUserVoiceChannel() (*discordgo.VoiceConnection, error) {
	vs, err := context.GetVoiceChannel()
	if err != nil {
		return nil, err
	}

	var voiceConnection *discordgo.VoiceConnection
	voiceConnection, err = context.Discord.ChannelVoiceJoin(vs.GuildID, vs.ChannelID, false, false)
	if err != nil {
		return nil, err
	}

	return voiceConnection, nil
}

func (context *Context) DM(content string) error {
	privateDM, err := context.Discord.UserChannelCreate(context.Message.Author.ID)
	if err != nil {
		return err
	}

	_, err = context.Discord.ChannelMessageSend(privateDM.ID, content)

	return err
}
