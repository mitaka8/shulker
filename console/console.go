package console

import (
	"fmt"

	"gitlab.com/mitaka8/shulker/registry"
)

type ConsoleReceiver struct {
}

func (c *ConsoleReceiver) ChatMessage(message registry.ChatMessage) error {
	fmt.Printf("[%v] %s: <%s> %s\n", message.Timestamp(), message.Source(), message.Author().Name(), message.Message())
	return nil
}

func (c *ConsoleReceiver) GenericMessage(message registry.GenericMessage) error {
	fmt.Println(message.Message())
	return nil
}

func (c *ConsoleReceiver) Start() error {
	return nil
}

func NewConsoleReceiver() registry.Receiver {
	return &ConsoleReceiver{}
}

func init() {
	registry.RegisterReceiver("Console", NewConsoleReceiver)
}
