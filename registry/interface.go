package registry

import (
	"encoding/json"
	"time"
)

type Author interface {
	AvatarURL() string
	Name() string
}

type ChatMessage interface {
	Author() Author
	Message() string
	Medium() string
	Source() string
	Timestamp() time.Time
}

type Receiver interface {
	Start() error
	GenericMessage(GenericMessage) error
	ChatMessage(ChatMessage) error
}
type GenericMessage interface {
	Message() string
	Source() string
	Medium() string
	Timestamp() time.Time
}

type Transmitter interface {
	ChatMessages() chan ChatMessage
	GenericMessages() chan GenericMessage
	Start() error
}

type Config struct {
	Left  string `json:"left"`
	Right string `json:"right"`

	LeftConfig  interface{} `json:"-"`
	RightConfig interface{} `json:"-"`

	RawLeftConfig  json.RawMessage `json:"left_config"`
	RawRightConfig json.RawMessage `json:"right_config"`
}
