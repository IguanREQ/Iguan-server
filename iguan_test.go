package main

import (
	"net"
	"net/http"
	"net/http/httptest"
	"net/rpc"
	"testing"
	"time"

	"github.com/amkulikov/extrpc"
	"iguan/auth"
	"iguan/dispatcher"
	"iguan/event"
	"iguan/listener"
	"iguan/subscriber"
)

func TestBasic(t *testing.T) {
	go dispatcher.RunDispatcher()

	subscriberAddr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		t.Log(*req)
	})).URL
	addr := "127.0.0.1:8080"
	go listener.RunTCP(addr)
	time.Sleep(1 * time.Second)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatal("dialing", err)
	}
	defer conn.Close()
	client := rpc.NewClientWithCodec(extrpc.NewGobClientCodec(conn))

	var reply interface{}

	conn.Write([]byte{auth.AuthTypeNone})
	err = client.Call("Event.Register", dispatcher.RegisterArgs{
		SourceTag: "test_tag",
		EventMask: "Ololo.Run",
		Subjects: []subscriber.SubjectNotifyInfo{
			{
				DestType: subscriber.DestTypeHttp,
				DestPath: subscriberAddr,
			},
		},
	}, &reply)
	if err != nil {
		t.Error(err)
	}
	t.Log(reply)

	err = client.Call("Event.Register", dispatcher.RegisterArgs{
		SourceTag: "test_tag",
		EventMask: "Ololo.*",
		Subjects: []subscriber.SubjectNotifyInfo{
			{
				DestType: subscriber.DestTypeHttp,
				DestPath: subscriberAddr,
			},
		},
	}, &reply)
	if err != nil {
		t.Error(err)
	}
	t.Log(reply)

	err = client.Call("Event.Register", dispatcher.RegisterArgs{
		SourceTag: "test_tag",
		EventMask: "Ololo.#",
		Subjects: []subscriber.SubjectNotifyInfo{
			{
				DestType: subscriber.DestTypeHttp,
				DestPath: subscriberAddr,
			},
		},
	}, &reply)
	if err != nil {
		t.Error(err)
	}
	t.Log(reply)

	for i := 0; i < 2; i++ {
		err = client.Call("Event.Fire", dispatcher.FireArgs{
			Event: &event.Event{
				SourceTag:  "test_tag",
				EmittedAt:  time.Time{},
				Delay:      0,
				Dispatcher: 1,
				Body: &event.Body{
					Class:       "YoloClass",
					Name:        "Ololo.Run",
					Payload:     "somepayload",
					PayloadType: "string",
				},
			},
			Caller: nil,
		}, &reply)
		if err != nil {
			t.Error(err)
		}
		t.Log(reply)
	}

	Wait()
}
