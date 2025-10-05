package shared_events

import (
	"fmt"
	"sync"
)

type Event interface {
	EventName() string
}

type EventHandler func(event Event)

type EventDispatcher struct {
	handlers map[string][]EventHandler
	mu       sync.RWMutex
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]EventHandler),
	}
}

func (d *EventDispatcher) Register(eventName string, handler EventHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.handlers[eventName] = append(d.handlers[eventName], handler)
}

func (d *EventDispatcher) Dispatch(eventName string, event Event) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if handlers, ok := d.handlers[eventName]; ok {
		for _, handler := range handlers {
			go handler(event)
		}
	} else {
		fmt.Printf("Event dispatched: %s (no handlers registered)\n", eventName)
	}
}
