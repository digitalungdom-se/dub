package pkg

import (
	"errors"
	"log"
	"regexp"
	"strings"

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

func (context *Context) ReplyNoMention(content string) (*discordgo.Message, error) {
	msg, err := context.Discord.ChannelMessageSend(context.TextChannel.ID, content)

	return msg, err
}

func (context *Context) ReplyEmbedNoMention(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	msg, err := context.Discord.ChannelMessageSendEmbed(context.TextChannel.ID, embed)

	return msg, err
}

func (context *Context) ReplyComplexNoMention(message *discordgo.MessageSend) (*discordgo.Message, error) {
	msg, err := context.Discord.ChannelMessageSendComplex(context.TextChannel.ID, message)

	return msg, err
}

func (context *Context) Reply(content string) (*discordgo.Message, error) {
	if !context.IsDM() {
		content = context.Message.Author.Mention() + ", " + strings.ToLower(string(content[0])) + content[1:]
	}

	msg, err := context.Discord.ChannelMessageSend(context.TextChannel.ID, content)

	return msg, err
}

func (context *Context) ReplyEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	complexMessage := &discordgo.MessageSend{
		Embed: embed,
	}

	if !context.IsDM() {
		complexMessage.Content = context.Message.Author.Mention()
	}

	msg, err := context.Discord.ChannelMessageSendComplex(context.TextChannel.ID, complexMessage)

	return msg, err
}

func (context *Context) ReplyComplex(message *discordgo.MessageSend) (*discordgo.Message, error) {
	if !context.IsDM() {
		originalContent := message.Content
		message.Content = context.Message.Author.Mention()

		if len(originalContent) > 3 {
			message.Content = message.Content + ", " + strings.ToLower(string(message.Content[0])) + message.Content[1:]
		} else {
			message.Content = message.Content + ", " + originalContent
		}
	}

	msg, err := context.Discord.ChannelMessageSendComplex(context.TextChannel.ID, message)

	return msg, err
}

func (context *Context) IsDM() bool {
	channel, err := context.Server.Discord.State.Channel(context.Message.ChannelID)
	if err != nil {
		if channel, err = context.Server.Discord.Channel(context.Message.ChannelID); err != nil {
			return false
		}
	}

	return channel.Type == discordgo.ChannelTypeDM
}

func (context *Context) DM(content string) error {
	privateDM, err := context.Discord.UserChannelCreate(context.Message.Author.ID)
	if err != nil {
		return err
	}

	_, err = context.Discord.ChannelMessageSend(privateDM.ID, content)

	return err
}

func (context *Context) GetMentions() ([]string, error) {
	if len(context.Message.Mentions) == 0 {
		return nil, errors.New("no mentions")
	}
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
	return nil, errors.New("could not find user's voice state")
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
