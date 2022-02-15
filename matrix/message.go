package matrix

import (
	"time"

	"gitlab.com/mitaka8/shulker/registry"
	mautrix_event "maunium.net/go/mautrix/event"
)

type MatrixAuthor struct {
	username  string
	avatarURL string
}

func (m *MatrixAuthor) Name() string {
	return m.username
}

func (m *MatrixAuthor) AvatarURL() string {
	return m.avatarURL
}

type MatrixMessage struct {
	author  *MatrixAuthor
	message string
	time    time.Time
	source  string
}

func (m *MatrixMessage) Medium() string {
	return "Matrix"
}
func (m *MatrixMessage) Author() registry.Author {
	return m.author
}
func (m *MatrixMessage) Source() string {
	return m.source
}
func (m *MatrixMessage) Message() string {
	return m.message
}
func (m *MatrixMessage) Timestamp() time.Time {
	return m.time
}
func NewMatrixMessage(message *mautrix_event.MessageEventContent, member *memberInfo) *MatrixMessage {
	return &MatrixMessage{
		author: &MatrixAuthor{
			username:  member.displayname,
			avatarURL: member.avatarUrl,
		},
		time:    time.Now(),
		message: message.Body,
	}
}
