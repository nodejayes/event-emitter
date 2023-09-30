package event_emitter_test

import (
	"fmt"
	"github.com/nodejayes/event-emitter"
	"sync"
	"testing"
)

var TestEvent = event_emitter.Event[string, int]{
	Token: "test_event",
	Handler: func(event string, testEventParams string) (int, error) {
		println(fmt.Sprintf("Event handeled: %v %v", event, testEventParams))
		return 0, nil
	},
	OnError: func(event string, err error) {
		println(fmt.Sprintf("Event error: %v %v", event, err.Error()))
	},
}

func TestNewEventEmitter(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	sub1 := event_emitter.Subscribe(TestEvent, func(params string, result int) {
		println(fmt.Sprintf("sub1 event %v fired with result %v and params %v", TestEvent.Token, result, params))
		wg.Done()
	})

	sub2 := event_emitter.Subscribe(TestEvent, func(params string, result int) {
		println(fmt.Sprintf("sub2 event %v fired with result %v and params %v", TestEvent.Token, result, params))
		wg.Done()
	})

	go event_emitter.Emit(TestEvent, "Hello")

	wg.Wait()
	event_emitter.Unsubscribe(sub1, sub2)
}

func TestUnsubscribe(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	sub1 := event_emitter.Subscribe(TestEvent, func(params string, result int) {
		println(fmt.Sprintf("sub1 event %v fired with result %v and params %v", TestEvent.Token, result, params))
		wg.Done()
	})

	sub2 := event_emitter.Subscribe(TestEvent, func(params string, result int) {
		t.Errorf("dont want to trigger unsubscribed Event")
	})
	event_emitter.Unsubscribe(sub2)

	go event_emitter.Emit(TestEvent, "Hello")

	wg.Wait()
	event_emitter.Unsubscribe(sub1)
}
