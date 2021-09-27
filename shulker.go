package shulker

import (
	"log"
	"sync"

	"gitlab.com/mitaka8/shulker/registry"
)

func NewFromConfig(c []registry.Config) {
	var wg sync.WaitGroup
	for i, bridge := range c {
		if registry.IsTransmitter(bridge.Left) && registry.IsReceiver(bridge.Right) {

			from, err := registry.ConstructTransmitter(bridge.Left, bridge.RawLeftConfig)
			if err != nil {
				log.Fatalf("Cannot construct bridge %v %v (left) transmitter: %v\n", i, bridge.Left, err)
			}

			to, err := registry.ConstructReceiver(bridge.Right, bridge.RawRightConfig)

			if err != nil {
				log.Fatalf("Cannot construct bridge %v %v (right) receiver: %v\n", i, bridge.Right, err)
			}
			run(&wg, from, to, "left")
		}

		if registry.IsTransmitter(bridge.Right) && registry.IsReceiver(bridge.Left) {

			from, err := registry.ConstructTransmitter(bridge.Right, bridge.RawRightConfig)
			if err != nil {
				log.Fatalf("Cannot construct bridge %v %v (right) transmitter: %v\n", i, bridge.Right, err)
			}

			to, err := registry.ConstructReceiver(bridge.Left, bridge.RawLeftConfig)

			if err != nil {
				log.Fatalf("Cannot construct bridge %v %v (left) receiver: %v\n", i, bridge.Left, err)
			}

			run(&wg, from, to, "right")

		}
	}
	wg.Wait()
}

func run(wg *sync.WaitGroup, from registry.Transmitter, to registry.Receiver, which string) {
	wg.Add(1)
	go (func() {
		for genericMessage := range from.GenericMessages() {
			go to.GenericMessage(genericMessage)
		}
	})()

	go (func() {
		for chatMessage := range from.ChatMessages() {
			println("sending from remote " + which)
			go to.ChatMessage(chatMessage)
			println("sent from remote " + which)
		}
		println("closed")
	})()
}
