package notifications

import (
	"sync"

	"github.com/production-grid/pgrid-core/pkg/applications"
	"github.com/production-grid/pgrid-core/pkg/logging"
)

var initOnce sync.Once

//NotifyingEventListener sends notifications based on incoming events
type NotifyingEventListener struct {
	DefaultRecipientListFunc applications.RecipientListFunc
	Transports               []Transport
	transportMap             map[applications.TransportType]Transport
}

//Dispatched is called whenever an event id processed
func (listener *NotifyingEventListener) Dispatched(def *applications.EventDef, event *applications.Event) {

	initOnce.Do(func() {

		listener.transportMap = make(map[applications.TransportType]Transport)

		for _, tx := range listener.Transports {
			listener.transportMap[tx.Type()] = tx
		}

	})

	if (def.RecipientListFunc != nil) || (listener.DefaultRecipientListFunc != nil) {

		recipientChannel := make(chan applications.Recipient, 100)

		if def.RecipientListFunc != nil {
			def.RecipientListFunc(recipientChannel, def, event)
		} else {
			listener.DefaultRecipientListFunc(recipientChannel, def, event)
		}

		for {
			rcpt, ok := <-recipientChannel
			if !ok {
				break
			}
			//TODO dupe checking
			logging.LogJSONWithName("RCPT", rcpt)
			if rcpt.TransportType == "" {
				listener.Transports[0].Send(def, event, &rcpt)
			} else {
				listener.transportMap[rcpt.TransportType].Send(def, event, &rcpt)
			}
		}

	}

}
