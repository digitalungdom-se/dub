package pkg

import (
	"errors"
	"time"

	"github.com/bwmarrin/discordgo"
)

type ReactionatorType int8

const (
	ReactionatorTypeHelp ReactionatorType = iota
	ReactionatorTypeController
)

type (
	ReactionListener struct {
		discord  *discordgo.Session
		Messages map[string]*reactionator
		Users    map[string]*activeReactionators
	}

	reactionator struct {
		discord          *discordgo.Session
		channelID        string
		Message          *discordgo.Message
		reactions        map[string]interface{}
		reactionOrder    []string
		defaultPage      interface{}
		ListenToRemove   bool
		removeReaction   bool
		reactionListener *ReactionListener
		messageType      ReactionatorType
		user             *discordgo.User
		isDM             bool
	}

	activeReactionators struct {
		Help *reactionator
	}
)

func NewReactionListener(discord *discordgo.Session) *ReactionListener {
	reactionListener := new(ReactionListener)
	reactionListener.discord = discord
	reactionListener.Messages = make(map[string]*reactionator)
	reactionListener.Users = make(map[string]*activeReactionators)

	return reactionListener
}

func (reactionListener *ReactionListener) listen(messageID string, reactionator *reactionator) {
	reactionListener.Messages[messageID] = reactionator
}

func (reactionListener *ReactionListener) React(message *discordgo.MessageReaction) error {
	var reactionator *reactionator
	var ok bool
	var err error

	if reactionator, ok = reactionListener.Messages[message.MessageID]; !ok {
		return nil
	}

	switch reactionator.reactions[message.Emoji.Name].(type) {
	case func(message *discordgo.MessageReaction):
		reactionator.reactions[message.Emoji.Name].(func(*discordgo.MessageReaction))(message)
	case func():
		reactionator.reactions[message.Emoji.Name].(func())()
	case *discordgo.MessageEmbed:
		_, err = reactionListener.discord.ChannelMessageEditEmbed(reactionator.channelID, reactionator.Message.ID,
			reactionator.reactions[message.Emoji.Name].(*discordgo.MessageEmbed))
	case string:
		reactionListener.discord.ChannelMessageEdit(reactionator.channelID, reactionator.Message.ID,
			reactionator.reactions[message.Emoji.Name].(string))
	default:
		if !reactionator.isDM {
			err = reactionListener.discord.MessageReactionRemove(message.ChannelID, message.MessageID, message.Emoji.Name, message.UserID)
			if err != nil {
				return err
			}
		}
	}

	if err != nil {
		return err
	}

	if !reactionator.isDM && reactionator.removeReaction {
		err = reactionListener.discord.MessageReactionRemove(message.ChannelID, message.MessageID, message.Emoji.Name, message.UserID)
	}

	return err
}

func NewReactionator(channelID string, discord *discordgo.Session, reactionListener *ReactionListener,
	listenToRemove bool, removeReaction bool, messageType ReactionatorType, user *discordgo.User) *reactionator {
	reactionator := new(reactionator)
	reactionator.channelID = channelID
	reactionator.discord = discord
	reactionator.reactions = make(map[string]interface{})
	reactionator.reactionListener = reactionListener
	reactionator.ListenToRemove = listenToRemove
	reactionator.removeReaction = removeReaction
	reactionator.messageType = messageType
	reactionator.user = user

	channel, err := discord.State.Channel(channelID)
	if err != nil {
		if channel, err = discord.Channel(channelID); err != nil {
			reactionator.isDM = false
		}
	}

	reactionator.isDM = channel.Type == discordgo.ChannelTypeDM

	return reactionator
}

func (reactionator *reactionator) AddDefaultPage(reaction string, content interface{}) error {
	switch content.(type) {
	case string:
	case *discordgo.MessageEmbed:
	default:
		return errors.New("invalid default message type")
	}

	reactionator.defaultPage = content

	var err error
	if reaction != "" {
		err = reactionator.Add(reaction, content)
	}

	return err
}

func (reactionator *reactionator) Add(reaction string, action interface{}) error {
	switch action.(type) {
	case func():
	case func(*discordgo.MessageReaction):
	case *discordgo.MessageEmbed:
	case string:
	default:
		return errors.New("invalid action type")
	}

	if _, ok := reactionator.reactions[reaction]; ok {
		return errors.New("reaction already exists")
	}

	reactionator.reactions[reaction] = action
	reactionator.reactionOrder = append(reactionator.reactionOrder, reaction)

	return nil
}

func (reactionator *reactionator) Initiate() error {
	var err error
	var msg *discordgo.Message

	switch reactionator.defaultPage.(type) {
	case string:
		content := reactionator.defaultPage.(string)
		msg, err = reactionator.discord.ChannelMessageSend(reactionator.channelID, content)
	case *discordgo.MessageEmbed:
		content := reactionator.defaultPage.(*discordgo.MessageEmbed)
		msg, err = reactionator.discord.ChannelMessageSendEmbed(reactionator.channelID, content)
	default:
		return errors.New("invalid default message type")
	}

	if err != nil {
		return err
	}

	time.Sleep(1 * time.Second)

	for _, value := range reactionator.reactionOrder {
		err = reactionator.discord.MessageReactionAdd(reactionator.channelID, msg.ID, value)
		if err != nil {
			return err
		}
	}

	reactionator.reactionListener.listen(msg.ID, reactionator)
	err = reactionator.discord.State.MessageAdd(msg)
	if err != nil {
		return err
	}

	reactionator.Message, err = reactionator.discord.State.Message(msg.ChannelID, msg.ID)

	if reactionator.messageType == ReactionatorTypeHelp {
		if reactionator.reactionListener.Users[reactionator.user.ID] != nil &&
			reactionator.reactionListener.Users[reactionator.user.ID].Help != nil {
			reactionator.reactionListener.Users[reactionator.user.ID].Help.Close()
		}

		if _, ok := reactionator.reactionListener.Users[reactionator.user.ID]; !ok {
			reactionator.reactionListener.Users[reactionator.user.ID] = new(activeReactionators)
		}

		reactionator.reactionListener.Users[reactionator.user.ID].Help = reactionator
	}

	return err
}

func (reactionator *reactionator) Close() {
	embed := NewEmbed().
		SetTitle("Denna sida är stängd").
		SetColor(16711680).
		MessageEmbed

	reactionator.discord.ChannelMessageEditEmbed(reactionator.channelID, reactionator.Message.ID, embed)
	reactionator.discord.State.MessageRemove(reactionator.Message)

	delete(reactionator.reactionListener.Messages, reactionator.Message.ID)
	if reactionator.messageType == ReactionatorTypeHelp {
		reactionator.reactionListener.Users[reactionator.user.ID].Help = nil
	}

	reactionator = nil
}

func (reactionator *reactionator) CloseButton() error {
	return reactionator.Add("🔥", reactionator.Close)
}

func (reactionator *reactionator) CloseAfter(duration time.Duration) {
	time.AfterFunc(duration, reactionator.Close)
}
