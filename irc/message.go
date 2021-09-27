package irc

import (
	"net/url"
	"time"

	"gitlab.com/mitaka8/shulker/registry"
)

type IrcUser struct {
	username string
}

func (u *IrcUser) Name() string {
	return u.username
}

func (u *IrcUser) AvatarURL() string {
	return "https://robohash.org/" + url.PathEscape(u.username) + ".png?set=4"
}

type IrcChatMessage struct {
	hostname  string
	timestamp time.Time
	author    *IrcUser
	message   string
}

func (m *IrcChatMessage) Author() registry.Author {
	return m.author
}
func (m *IrcChatMessage) Medium() string {
	return "IRC"
}
func (m *IrcChatMessage) Source() string {
	return m.hostname
}
func (m *IrcChatMessage) Message() string {
	return m.message
}
func (m *IrcChatMessage) Timestamp() time.Time {
	return m.timestamp
}

func NewIrcChatMessage(hostname, username, message string) *IrcChatMessage {
	return &IrcChatMessage{
		hostname:  hostname,
		timestamp: time.Now(),
		author:    &IrcUser{username: username},
		message:   message,
	}
}
