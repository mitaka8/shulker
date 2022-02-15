package matrix

import (
	"fmt"
	"log"
	"time"

	"gitlab.com/mitaka8/shulker/registry"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type MatrixTransmitter struct {
	Homeserver string `json:"homeserver"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	RoomId     string `json:"room_id"`
	Database   string `json:"database"`
	DeviceId   string `json:"device_id"`
	Key        string `json:"key"`

	client     *mautrix.Client    `json:"-"`
	olmMachine *crypto.OlmMachine `json:"-"`
	startedAt  int64
}

func (m *MatrixTransmitter) Start() error {
	client, olmMachine, err := getOrMakeClient(m.Homeserver, m.Username, m.Password, m.Database, m.Key, m.DeviceId)
	if err != nil {
		return err
	}
	m.client = client
	m.olmMachine = olmMachine
	m.startedAt = time.Now().Unix() * 1000
	return err
}

func (m *MatrixTransmitter) ChatMessages() chan registry.ChatMessage {
	chatMessages := make(chan registry.ChatMessage)

	syncer := m.client.Syncer.(*mautrix.DefaultSyncer)

	var onMessage = func(evt *event.Event) {
		if evt.Sender == m.client.UserID {
			return
		}
		if evt.Timestamp < m.startedAt {
			return
		}
		member, err := getMember(evt.Sender, m.client, id.RoomID(m.RoomId))
		if err != nil {
			fmt.Println(err)
		}
		chatMessages <- NewMatrixMessage(evt.Content.AsMessage(), member)

	}
	syncer.OnEventType(event.EventMessage, func(source mautrix.EventSource, evt *event.Event) {
		onMessage(evt)
	})
	syncer.OnEventType(event.EventEncrypted, func(source mautrix.EventSource, evt *event.Event) {
		decryptedEvent, err := m.olmMachine.DecryptMegolmEvent(evt)
		if evt.Timestamp < m.startedAt {
			return
		}
		if err != nil {
			log.Printf("failed to decrypt unknown event: %v\n", err)
		}
		if decryptedEvent.Type != event.EventMessage {
			return
		}
		onMessage(evt)
	})

	return chatMessages
}
func (m *MatrixTransmitter) GenericMessages() chan registry.GenericMessage {
	genericMessages := make(chan registry.GenericMessage)
	return genericMessages
}

func init() {
	registry.RegisterTransmitter("Matrix", func() registry.Transmitter { return &MatrixTransmitter{} })
}
