# event-emitter
a Go Event Emitter

## Example

```go
// define a Event
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

...

func main() {
    var wg sync.WaitGroup
    wg.Add(1)
	
	// subscribe to a Event
	sub := event_emitter.Subscribe(TestEvent, func(params string, result int) {
        println(fmt.Sprintf("sub1 event %v fired with result %v and params %v", TestEvent.Token, result, params))
        wg.Done()
    })
	
	// emit a Event async
    go event_emitter.Emit(TestEvent, "Hello")
	
	wg.Wait()
	
	// unsubscribe from event
    event_emitter.Unsubscribe(sub)
}
```
