package discord

import (
	"os"

	"github.com/andersfylling/disgord"
	"github.com/sirupsen/logrus"
	"gitlab.com/mitaka8/shulker/registry"
)

var clients = make(map[string]*disgord.Client)

var log = &logrus.Logger{
	Out:       os.Stderr,
	Formatter: new(logrus.TextFormatter),
	Hooks:     make(logrus.LevelHooks),
	Level:     logrus.WarnLevel,
}

func getClient(botToken string) (*disgord.Client, error) {
	client, ok := clients[botToken]
	if !ok {
		client = disgord.New(disgord.Config{
			BotToken: botToken,
			Logger:   log,
		})
		clients[botToken] = client
		if err := client.Gateway().Connect(); err != nil {
			return nil, err
		}
	}
	return client, nil
}

type DiscordTransmitter struct {
	BotToken  string `json:"bot_token"`
	ChannelID string `json:"channel"`

	discord        *disgord.Client
	messageChannel chan registry.ChatMessage
	genericChannel chan registry.GenericMessage
	selfId         disgord.Snowflake
}

func (d *DiscordTransmitter) Start() error {
	client, err := getClient(d.BotToken)
	if err != nil {
		return err
	}
	d.discord = client
	d.genericChannel = make(chan registry.GenericMessage, 1)
	d.messageChannel = make(chan registry.ChatMessage, 1)

	d.discord.Gateway().BotReady(func() {
		user, err := d.discord.CurrentUser().Get()
		if err == nil {
			d.selfId = user.ID
		}
		d.discord.Gateway().MessageCreate(func(s disgord.Session, evt *disgord.MessageCreate) {
			if evt.Message.ChannelID.String() != d.ChannelID {
				return
			}
			if evt.Message.WebhookID != 0 {
				return
			}
			if evt.Message.Author.ID == d.selfId {
				return
			}
			guild, err := s.Guild(evt.Message.GuildID).Get()
			if err != nil {
				return
			}
			d.messageChannel <- newMessage(evt, guild.Name)
		})

	})

	return nil
}

func (d *DiscordTransmitter) ChatMessages() chan registry.ChatMessage {
	return d.messageChannel
}

func (d *DiscordTransmitter) GenericMessages() chan registry.GenericMessage {
	return d.genericChannel
}

func NewDiscordTransmitter() registry.Transmitter {
	return &DiscordTransmitter{}
}

func init() {
	registry.RegisterTransmitter("Discord", NewDiscordTransmitter)
}
