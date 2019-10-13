package events

import (
	"sync"

	"github.com/production-grid/pgrid-core/pkg/applications"
)

var initOnce sync.Once

var dispatcherChannel chan *applications.Event

//DispatchEvent dispatches an event
func DispatchEvent(event *applications.Event) {

	initOnce.Do(func() {

		startDispatcher()

	})

	go queueEvent(event)

}

func queueEvent(event *applications.Event) {

	dispatcherChannel <- event

}

func startDispatcher() {

	dispatcherChannel = make(chan *applications.Event, 1024)

	go startEventProcessor()

}

func startEventProcessor() {

	for {
		e := <-dispatcherChannel
		for _, listener := range applications.CurrentApplication.EventListeners {
			go listener.Dispatched(e)
		}
	}

}

//Dispatch dispatches an event based on the input parameters
func Dispatch(key string, session applications.Session, metaData interface{}) {

	DispatchEvent(&applications.Event{
		Key:      key,
		Session:  session,
		MetaData: metaData,
	})
}
