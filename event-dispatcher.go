package event_dispatcher

import (
	"reflect"
	"runtime"
	"sort"
)

type Event interface{}
type StoppableEvent interface {
	IsPropagationStopped() bool
	StopPropagation()
}

type Name string

type Listener func(Event, Name)

type EventDispatcher struct {
	listeners  map[Name]map[int][]Listener
	priorities map[Name][]int
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		listeners:  map[Name]map[int][]Listener{},
		priorities: map[Name][]int{},
	}
}
func (e *EventDispatcher) Dispatch(event Event, name Name) {
	if prioritarizedListeners, ok := e.listeners[name]; ok {
		stoppableEvent, stoppable := event.(StoppableEvent)
		for _, priority := range e.priorities[name] {
			if listeners, ok := prioritarizedListeners[priority]; ok {
				for _, listener := range listeners {
					if stoppable && stoppableEvent.IsPropagationStopped() {
						return
					}
					listener(event, name)
				}
			}
		}
	}
}

func (e *EventDispatcher) HasListener(name Name) bool {
	_, ok := e.listeners[name]
	return ok
}


func (e *EventDispatcher) AddListener(name Name, listener Listener, priority int) bool {
	if _, ok := e.listeners[name]; !ok {
		// initialize listener with name
		e.listeners[name] = map[int][]Listener{}
		e.priorities[name] = []int{}
	}

	if _, ok := e.listeners[name][priority]; !ok {
		// initialize listener with this priority
		e.listeners[name][priority] = []Listener{listener}
		e.priorities[name] = append(e.priorities[name], priority)
		sort.Ints(e.priorities[name])
		return true
	}

	// simply append listener
	e.listeners[name][priority] = append(e.listeners[name][priority], listener)
	return true
}

func (e *EventDispatcher) RemoveListener(name Name, listener Listener, priority int) bool {
	listeners, ok := e.listeners[name][priority]
	if !ok {
		return false
	}
	deleted := false
	for i, l := range listeners {
		if runtime.FuncForPC(reflect.ValueOf(listener).Pointer()).Name() == runtime.FuncForPC(reflect.ValueOf(l).Pointer()).Name() {
			e.listeners[name][priority] = append(listeners[:i], listeners[i+1:]...)
			// TODO CLEAN EMPTY PRIORITY
			deleted = true
		}
	}
	return deleted
}
