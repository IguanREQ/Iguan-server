package dispatcher

import (
	"container/list"
	"sync"
	"time"

	"iguan/event"
	"iguan/logs"
	"iguan/subscriber"
)

type EventsQueue struct {
	list *list.List
	mu   sync.RWMutex
}

var queue *EventsQueue

func init() {
	queue = &EventsQueue{
		list: list.New(),
	}
}

func AddEvent(e *event.Event) error {
	// TODO: add sync.Pool for memory reusing
	queue.mu.Lock()
	defer queue.mu.Unlock()

	var ep *event.Event

	// add event in queue by time
	for el := queue.list.Back(); el != nil; el = el.Prev() {
		ep = el.Value.(*event.Event)
		if ep.EmittedAt.Before(e.EmittedAt) {
			queue.list.InsertAfter(e, el)
			return nil
		}
	}
	// push to front if earlier events not found
	el := queue.list.PushFront(e)
	_ = el
	return nil
}

func SyncQueueToDisk() error {
	// TODO: implement disk-caching for queue
	return nil
}

func fireFirst() error {
	queue.mu.Lock()
	defer queue.mu.Unlock()

	fr := queue.list.Front()
	if fr == nil {
		return nil
	}

	fre := fr.Value.(*event.Event)
	queue.list.Remove(fr)
	return subscriber.Fire(fre)
}

func RunDispatcher() error {
	var err error
	for {
		err = fireFirst()
		if err != nil {
			logs.Error("Fire error: %s", err)
		}
		time.Sleep(1 * time.Second)
	}
}
