package irc

import (
	"crypto/tls"
	"fmt"
	"net/url"
	"strings"

	irc "github.com/thoj/go-ircevent"
	"gitlab.com/mitaka8/shulker/registry"
)

type IrcTransmitter struct {
	Dest string `json:"dest"`
	Name string `json:"name"`

	channel string `json:"-"`
	conn    *irc.Connection

	from string

	genericMessages chan registry.GenericMessage
	chatMessages    chan registry.ChatMessage
}

func (i *IrcTransmitter) Start() error {
	conn, ok := clients[i.Dest]

	i.chatMessages = make(chan registry.ChatMessage)
	i.genericMessages = make(chan registry.GenericMessage)

	u, err := url.Parse(i.Dest)

	if err != nil {
		return fmt.Errorf("cannot read IRC destination url: %v", err)
	}

	i.from = i.Name
	if i.Name == "" {
		i.from = u.Hostname()
	}
	if !ok {
		pass, _ := u.User.Password()

		conn := irc.IRC(u.User.Username(), u.User.Username())

		if conn.UseTLS = u.Scheme == "ircs"; conn.UseTLS {
			conn.TLSConfig = &tls.Config{
				ServerName: u.Hostname(),
			}
		}

		if pass != "" {
			conn.SASLLogin = u.User.Username()
			conn.UseSASL = true
			conn.SASLPassword = pass
		}
		if err := conn.Connect(u.Host); err != nil {
			fmt.Println(err)
			return err
		}
		i.channel = "#" + u.Fragment
		conn.AddCallback("001", func(e *irc.Event) {
			conn.Join(i.channel)
		})
		i.conn = conn
		return nil
	}
	i.conn = conn
	i.conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		if len(e.Arguments) < 1 {
			return
		}
		i.chatMessages <- NewIrcChatMessage(i.from, nickFromIdent(e.Source), e.Message())

	})

	return nil

}

func (i *IrcTransmitter) ChatMessages() chan registry.ChatMessage {
	return i.chatMessages
}
func (i *IrcTransmitter) GenericMessages() chan registry.GenericMessage {
	return i.genericMessages
}
func NewIrcTransmitter() registry.Transmitter {
	return &IrcTransmitter{}
}

func nickFromIdent(ident string) string {
	parts := strings.Split(ident, "!")
	return parts[0]
}

func init() {
	registry.RegisterTransmitter("IRC", NewIrcTransmitter)
}
