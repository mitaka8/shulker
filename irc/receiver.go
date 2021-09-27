package irc

import (
	"crypto/tls"
	"fmt"
	"net/url"
	"strings"

	irc "github.com/thoj/go-ircevent"
	"gitlab.com/mitaka8/shulker/registry"
)

var clients = make(map[string]*irc.Connection)

type IrcReceiver struct {
	Dest string `json:"dest"`

	channel string `json:"-"`
	conn    *irc.Connection
}

func (i *IrcReceiver) Start() error {
	conn, ok := clients[i.Dest]
	if !ok {
		u, err := url.Parse(i.Dest)
		if err != nil {
			return fmt.Errorf("cannot read IRC destination url: %v", err)
		}
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
		clients[i.Dest] = conn
		return nil
	}
	i.conn = conn
	return nil
}

func (i *IrcReceiver) ChatMessage(message registry.ChatMessage) error {
	var sb strings.Builder
	sb.WriteString("[")
	sb.WriteString(message.Medium())
	sb.WriteString("] ")
	sb.WriteString(message.Source())
	sb.WriteString(": <")
	sb.WriteString(message.Author().Name())
	sb.WriteString("> ")
	sb.WriteString(message.Message())
	i.conn.Privmsg(i.channel, sb.String())
	return nil
}
func (i *IrcReceiver) GenericMessage(message registry.GenericMessage) error {
	var sb strings.Builder
	sb.WriteString("[")
	sb.WriteString(message.Medium())
	sb.WriteString("] ")
	sb.WriteString(message.Source())
	sb.WriteString(": ")
	sb.WriteString(message.Message())
	i.conn.Privmsg(i.channel, sb.String())
	return nil
}
func NewIrcReceiver() registry.Receiver {
	return &IrcReceiver{}
}

func init() {
	registry.RegisterReceiver("IRC", NewIrcReceiver)
}
