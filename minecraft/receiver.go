package minecraft

import (
	"encoding/json"
	"strings"

	"github.com/james4k/rcon"
	"gitlab.com/mitaka8/shulker/registry"
)

var defaultName = "Minecraft"

type MinecraftChatReceiver struct {
	RconHost     string `json:"rcon_host"`
	RconPassword string `json:"rcon_pass"`

	rcon *rcon.RemoteConsole
}

func (m *MinecraftChatReceiver) connectToRcon() error {
	conn, err := rcon.Dial(m.RconHost, m.RconPassword)
	if err != nil {
		return err
	}

	m.rcon = conn

	return nil
}

func (m *MinecraftChatReceiver) Start() error {
	return nil
}

func (m *MinecraftChatReceiver) ChatMessage(message registry.ChatMessage) error {
	var sb strings.Builder
	sb.WriteString(`/tellraw @a ["",{"text":"[","color":"blue"},{"text":`)
	sb.WriteString(mustJsonMarshal(message.Medium()))
	sb.WriteString(`,"color":"blue"},{"text":"] ","color":"blue"},{"text":`)
	sb.WriteString(mustJsonMarshal(message.Source()))
	sb.WriteString(`,"color":"white"},{"text":": <","color":"white"},{"text":`)
	sb.WriteString(mustJsonMarshal(message.Author().Name()))
	sb.WriteString(`,"color":"white"},{"text":"> ","color":"white"},`)
	sb.WriteString(mustJsonMarshal(message.Message()))
	sb.WriteString(`]`)

	if err := m.connectToRcon(); err != nil {
		return err
	}
	_, err := m.rcon.Write(sb.String())
	return err
}
func (m *MinecraftChatReceiver) GenericMessage(message registry.GenericMessage) error {
	var sb strings.Builder
	sb.WriteString(`/tellraw @a ["",{"text":"[","color":"blue"},{"text":`)
	sb.WriteString(mustJsonMarshal(message.Medium()))
	sb.WriteString(`,"color":"blue"},{"text":"]","color":"blue"},{"text":`)
	sb.WriteString(mustJsonMarshal(message.Source()))
	sb.WriteString(`,"color":"white"},{"text":": ","color":"white"},`)
	sb.WriteString(mustJsonMarshal(message.Message()))
	sb.WriteString(`]`)

	if err := m.connectToRcon(); err != nil {
		return err
	}

	_, err := m.rcon.Write(sb.String())

	return err
}

func mustJsonMarshal(message string) string {
	marshalled, _ := json.Marshal(message)
	return string(marshalled)
}
func NewMinecraftReceiver() registry.Receiver {
	return &MinecraftChatReceiver{}
}

func init() {
	registry.RegisterReceiver("Minecraft", NewMinecraftReceiver)
}
