package discord

import (
	"time"

	"github.com/andersfylling/disgord"
	"gitlab.com/mitaka8/shulker/registry"
)

type DiscordUser struct {
	nickname string
	avatar   string
}

func (d *DiscordUser) AvatarURL() string {
	return d.avatar
}
func (d *DiscordUser) Name() string {
	return d.nickname
}

type DiscordMessage struct {
	author    *DiscordUser
	timestamp time.Time
	content   string
	source    string
}

func (d DiscordMessage) Author() registry.Author {
	return d.author
}

func (d DiscordMessage) Medium() string {
	return "Discord"
}
func (d DiscordMessage) Message() string {
	return d.content
}

func (d DiscordMessage) Source() string {
	return d.source
}

func (d DiscordMessage) Timestamp() time.Time {
	return d.timestamp
}
func newMessage(message *disgord.MessageCreate, guildName string) registry.ChatMessage {
	nick := message.Message.Author.Username
	if message.Message.Member.Nick != "" {
		nick = message.Message.Member.Nick
	}

	return &DiscordMessage{
		author: &DiscordUser{
			nickname: nick,
			avatar:   message.Message.Author.Avatar,
		},

		timestamp: time.UnixMilli(message.Message.Timestamp.UnixMilli()), // Can't implicitely cast.  So copy.
		source:    guildName,
		content:   message.Message.Content,
	}
}
