package event_emitter

import (
	"sync"
)

var eventMutex sync.RWMutex
var subs = make(map[string][]func(params any, result any))

type (
	// Event descibes a Event that can emit/subscribe with event_emitter
	Event[TParam any, TResult any] struct {
		Token   string
		Handler func(event string, params TParam) (TResult, error)
		OnError func(event string, err error)
	}
	// Subscription for a Event fires when Event was emitted and executed
	Subscription struct {
		Id    int
		Token string
	}
)

// Emit executes a Event and informs all Subscribers with the Params and the Result
func Emit[TParam any, TResult any](event Event[TParam, TResult], params TParam) {
	eventMutex.RLock()
	defer eventMutex.RUnlock()

	result, err := event.Handler(event.Token, params)
	if err != nil {
		event.OnError(event.Token, err)
		return
	}
	subscriptions := getSubscriptions(event.Token)
	for _, subscription := range subscriptions {
		subscription(params, result)
	}
}

// Subscribe to a Event execution and get the Subscription Information to unsubscribe
func Subscribe[TParams any, TResult any](event Event[TParams, TResult], subscriber func(params TParams, result TResult)) Subscription {
	eventMutex.Lock()
	defer eventMutex.Unlock()

	sub := Subscription{
		Token: event.Token,
	}
	actions := getSubscriptions(event.Token)
	if len(actions) < 1 {
		subs[event.Token] = make([]func(params any, result any), 0)
	}
	subs[event.Token] = append(subs[event.Token], func(params any, result any) {
		p, okParams := params.(TParams)
		if !okParams {
			return
		}
		r, okResult := result.(TResult)
		if !okResult {
			return
		}
		subscriber(p, r)
	})
	sub.Id = len(subs[event.Token]) - 1
	return sub
}

// Unsubscribe from a Event execution
func Unsubscribe(subscriptions ...Subscription) {
	eventMutex.Lock()
	defer eventMutex.Unlock()

	for _, subscription := range subscriptions {
		actions := getSubscriptions(subscription.Token)
		if (len(actions) - 1) < subscription.Id {
			return
		}
		subs[subscription.Token] = append(subs[subscription.Token][:subscription.Id], subs[subscription.Token][subscription.Id+1:]...)
	}
}

func getSubscriptions(token string) []func(params any, result any) {
	action, hasAction := subs[token]
	if !hasAction {
		return make([]func(params any, result any), 0)
	}
	return action
}
