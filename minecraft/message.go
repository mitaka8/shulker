package minecraft

import (
	"strings"
	"time"

	"gitlab.com/mitaka8/shulker/registry"
)

type MinecraftUser struct {
	username string
}

type MinecraftMessage struct {
	author    *MinecraftUser
	message   string
	timestamp time.Time
	source    string
}

func (u MinecraftUser) AvatarURL() string {
	var sb strings.Builder
	sb.WriteString("https://minotar.net/armor/bust/")
	sb.WriteString(u.username)
	sb.WriteString("/300.png")
	return sb.String()
}

func (u MinecraftUser) Name() string {
	return u.username
}

func (m MinecraftGenericMessage) Medium() string {
	return "Minecraft"
}
func (m MinecraftMessage) Medium() string {
	return "Minecraft"
}

func (m MinecraftMessage) Message() string {
	return m.message
}

func (m MinecraftMessage) Source() string {
	return m.source
}
func (m MinecraftMessage) Author() registry.Author {
	return m.author
}

func (m MinecraftMessage) Timestamp() time.Time {
	return m.timestamp
}

func NewMinecraftMessage(username, message, name string) *MinecraftMessage {
	return &MinecraftMessage{
		author: &MinecraftUser{
			username: username,
		},
		message:   message,
		timestamp: time.Now(),
		source:    name,
	}
}

type MinecraftGenericMessage struct {
	message   string
	timestamp time.Time
	source    string
}

func (m *MinecraftGenericMessage) Message() string {
	return m.message
}
func (m *MinecraftGenericMessage) Source() string {
	return m.source
}
func (m *MinecraftGenericMessage) Timestamp() time.Time {
	return m.timestamp
}

func newMinecraftGenericMessage(message string, name string) *MinecraftGenericMessage {
	return &MinecraftGenericMessage{
		timestamp: time.Now(),
		message:   message,
		source:    name,
	}
}
