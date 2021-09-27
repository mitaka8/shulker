package registry

import (
	"encoding/json"
)

var receiverRegistry = make(map[string](func() Receiver))
var transmitterRegistry = make(map[string]func() Transmitter)

func RegisterReceiver(name string, config func() Receiver) {
	receiverRegistry[name] = config
}
func RegisterTransmitter(name string, config func() Transmitter) {
	transmitterRegistry[name] = config
}

func IsTransmitter(typ string) bool {
	_, exists := transmitterRegistry[typ]
	return exists
}

func IsReceiver(typ string) bool {

	_, exists := receiverRegistry[typ]
	return exists
}

func ConstructTransmitter(typ string, mes json.RawMessage) (Transmitter, error) {
	txbridge := transmitterRegistry[typ]()
	err := json.Unmarshal(mes, txbridge)
	if err != nil {
		return nil, err
	}
	if err = txbridge.Start(); err != nil {
		return nil, err
	}
	return txbridge, err
}

func ConstructReceiver(typ string, mes json.RawMessage) (Receiver, error) {
	rxbridge := receiverRegistry[typ]()
	err := json.Unmarshal(mes, rxbridge)

	if err != nil {
		return nil, err
	}
	if err = rxbridge.Start(); err != nil {
		return nil, err
	}
	return rxbridge, err
}
