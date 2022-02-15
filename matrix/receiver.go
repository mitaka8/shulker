package matrix

import (
	"gitlab.com/mitaka8/shulker/registry"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type MatrixReceiver struct {
	Homeserver string          `json:"homeserver"`
	Username   string          `json:"username"`
	Password   string          `json:"password"`
	RoomId     string          `json:"room_id"`
	Database   string          `json:"database"`
	Key        string          `json:"database"`
	DeviceId   string          `json:"device_id"`
	client     *mautrix.Client `json:"-"`
}

func (m *MatrixReceiver) Start() error {
	client, _, err := getOrMakeClient(m.Homeserver, m.Username, m.Password, m.Database, m.Key, m.DeviceId)
	if err != nil {
		return err
	}
	m.client = client
	return err
}

func (m *MatrixReceiver) ChatMessage(message registry.ChatMessage) error {
	_, err := m.client.SendMessageEvent(id.RoomID(m.RoomId), event.EventMessage, formatMessage(message))
	return err
}
func (m *MatrixReceiver) GenericMessage(message registry.GenericMessage) error {

	_, err := m.client.SendText(id.RoomID(m.RoomId), formatGenericMessage(message))
	return err
}

func init() {
	registry.RegisterReceiver("Matrix", func() registry.Receiver { return &MatrixReceiver{} })
}
