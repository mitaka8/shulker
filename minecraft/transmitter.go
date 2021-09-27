package minecraft

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/hpcloud/tail"
	"gitlab.com/mitaka8/shulker/registry"
)

type MinecraftChatTransmitter struct {
	LogFile string  `json:"log_file"`
	Name    *string `json:"name"`

	messageChannel chan registry.ChatMessage
	genericChannel chan registry.GenericMessage
}

func (m *MinecraftChatTransmitter) Start() error {

	if m.Name == nil {
		m.Name = &defaultName
	}
	if m.LogFile != "" {
		logFileLocation, err := url.Parse(m.LogFile)
		if err != nil {
			return fmt.Errorf("cannot read log file location: %v", err)
		}

		var reader *tail.Tail
		switch logFileLocation.Scheme {
		case "sftp":
			password, ok := logFileLocation.User.Password()
			if !ok {
				password = ""
			}
			reader, err = tailSftpFile(logFileLocation.Host, logFileLocation.User.Username(), password, logFileLocation.Path)
			if err != nil {
				return err
			}
		case "file":
			reader, err = tailLocalFile(logFileLocation.Path)
		}

		if err != nil {
			return err
		}

		chatMessage := regexp.MustCompile(`\[\d\d:\d\d:\d\d\] \[Server thread\/INFO\]: \<([a-zA-Z0-9_]+)\> (.+)`)
		joinMessage := regexp.MustCompile(`\[\d\d:\d\d:\d\d\] \[Server thread\/INFO\]: ([a-zA-Z0-9_]+) joined the game`)
		leaveMessage := regexp.MustCompile(`\[\d\d:\d\d:\d\d\] \[Server thread\/INFO\]: ([a-zA-Z0-9_]+) left the game`)

		m.messageChannel = make(chan registry.ChatMessage)
		m.genericChannel = make(chan registry.GenericMessage)

		go (func() {
			for line := range reader.Lines {

				if chatMessage.MatchString(line.Text) {
					matches := chatMessage.FindStringSubmatch(line.Text)
					if (len(matches) < 2) || (matches[1] == "") || (matches[2] == "") {
						continue
					}
					println("send to channel ", matches[2])
					m.messageChannel <- NewMinecraftMessage(matches[1], matches[2], *m.Name)
					continue
				}
				if joinMessage.MatchString(line.Text) {
					matches := joinMessage.FindStringSubmatch(line.Text)
					if len(matches) < 2 {
						continue
					}
					m.genericChannel <- m.joinedTheGame(matches[1])
				}
				if leaveMessage.MatchString(line.Text) {
					matches := leaveMessage.FindStringSubmatch(line.Text)
					if len(matches) < 2 {
						continue
					}
					m.genericChannel <- m.leftTheGame(matches[1])
				}
			}
		})()
	}

	return nil
}

func (m *MinecraftChatTransmitter) ChatMessages() chan registry.ChatMessage {
	return m.messageChannel
}

func (m *MinecraftChatTransmitter) GenericMessages() chan registry.GenericMessage {
	return m.genericChannel
}

func (m *MinecraftChatTransmitter) joinedTheGame(username string) registry.GenericMessage {
	var sb strings.Builder
	sb.WriteString(username)
	sb.WriteString(" joined the game")
	return newMinecraftGenericMessage(sb.String(), *m.Name)
}

func (m *MinecraftChatTransmitter) leftTheGame(username string) registry.GenericMessage {
	var sb strings.Builder
	sb.WriteString(username)
	sb.WriteString(" left the game")
	return newMinecraftGenericMessage(sb.String(), *m.Name)
}

func NewMinecraftTransmitter() registry.Transmitter {
	return &MinecraftChatTransmitter{}
}

func init() {
	registry.RegisterTransmitter("Minecraft", NewMinecraftTransmitter)
}
