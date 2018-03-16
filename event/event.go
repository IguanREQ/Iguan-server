package event

import (
	"time"
)

type EmitterType uint8

const (
	EmitterLangPhp EmitterType = iota + 1
)

type Body struct {
	Class       string      `json:"class"`
	Name        string      `json:"name"`
	Payload     interface{} `json:"payload"`
	PayloadType string      `json:"payloadType"`
}

func (b *Body) Valid() bool {
	return true
}

type Event struct {
	SourceTag  string        `json:"sourceTag"`
	EmittedAt  time.Time     `json:"emittedAt"`
	Delay      time.Duration `json:"delay"`
	Dispatcher EmitterType   `json:"dispatcher"`
	Body       *Body         `json:"body"`
}

func (e *Event) Valid() bool {
	if e.Body == nil {
		return false
	}
	return true
}
