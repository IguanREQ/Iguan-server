package dispatcher

import (
	"errors"

	"github.com/amkulikov/extrpc"
	"iguan/auth"
	"iguan/event"
	"iguan/listener"
	"iguan/logs"
	"iguan/subscriber"
)

var ErrInvalidArgs = errors.New("Dispatcher: Invalid arguments")

func init() {
	listener.RegisterMethods(new(Event))
}

type Event bool

type FireArgs struct {
	Event *event.Event `json:"event"`

	Caller extrpc.Caller
}

func (a *FireArgs) Valid() bool {
	if a.Caller == nil || a.Event == nil || !a.Event.Valid() {
		return false
	}
	return true
}

type RegisterArgs struct {
	SourceTag string                         `json:"sourceTag"`
	EventMask string                         `json:"eventMask"`
	Subjects  []subscriber.SubjectNotifyInfo `json:"subjects"`

	Caller extrpc.Caller
}

func (a *RegisterArgs) Valid() bool {
	if a.Subjects == nil || a.Caller == nil {
		return false
	}
	return true
}

type UnregisterAllArgs struct {
	sourceTag string

	Caller *auth.Caller
}

func (a *UnregisterAllArgs) Valid() bool {
	if a.Caller == nil {
		return false
	}
	return true
}

func (e *Event) Fire(args *FireArgs, reply *interface{}) error {
	logs.Info("RPC: Fire called")
	if !args.Valid() {
		logs.Info("RPC: Invalid args")
		return ErrInvalidArgs
	}

	return AddEvent(args.Event)
}

func (e *Event) Register(args *RegisterArgs, reply *interface{}) error {
	logs.Info("RPC: Register called")
	if !args.Valid() {
		logs.Info("RPC: Invalid args")
		return ErrInvalidArgs
	}
	for _, subj := range args.Subjects {
		err := subscriber.Register(args.SourceTag, args.EventMask, &subj)
		if err != nil {
			logs.Error("RPC: Register failed: %s (%v)", err, args)
		}
	}
	return nil
}

/*
func (e *Event) Unregister(args *RegisterArgs, reply struct{}) error {
	if !args.Valid() {
		return ErrInvalidArgs
	}
	for i, subj := range args.subjects {
		subscriber.Unregister(&subj)
	}
	return nil
}

func (e *Event) UnregisterAll(args *UnregisterAllArgs, reply struct{}) error {
	if !args.Valid() {
		return ErrInvalidArgs
	}
	return subscriber.UnregisterAll()
}
*/
