package discord

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/andersfylling/disgord"
	"gitlab.com/mitaka8/shulker/registry"
)

type DiscordReceiver struct {
	BotToken  string `json:"bot_token"`
	Webhook   string `json:"webhook"`
	ChannelID string `json:"channel"`
	discord   *disgord.Client
}

func (d *DiscordReceiver) Start() error {
	client, err := getClient(d.BotToken)
	if err != nil {
		return err
	}
	d.discord = client
	return nil
}

func (d *DiscordReceiver) ChatMessage(message registry.ChatMessage) error {
	webhookPayload, err := json.Marshal(&disgord.ExecuteWebhookParams{
		Content:   message.Message(),
		Username:  message.Author().Name(),
		AvatarURL: message.Author().AvatarURL(),
	})

	if err != nil {
		return err
	}
	msg, err := http.Post(d.Webhook, "application/json", bytes.NewReader(webhookPayload))
	io.Copy(os.Stdout, msg.Body)
	return err
	// fmt.Printf("[%v] %s@%s: <%s> %s\n", message.Timestamp(), message.Medium(), message.Source(), message.Author().Name(), message.Message())
}

func (d *DiscordReceiver) GenericMessage(message registry.GenericMessage) error {
	_, err := d.discord.SendMsg(disgord.ParseSnowflakeString(d.ChannelID), message.Message())
	return err

}

func NewDiscordReceiver() registry.Receiver {
	return &DiscordReceiver{}
}

func init() {
	registry.RegisterReceiver("Discord", NewDiscordReceiver)
}
