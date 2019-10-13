package applications

import (
	"github.com/sirupsen/logrus"
)

//EventMetaDataFactory allows the event system to create a new instance of the
//event meta data struct if it has to be re-consituted in a different process
type EventMetaDataFactory func() interface{}

// EventDef models a module event definition
type EventDef struct {
	Key          string
	TenantScoped bool

	// log level for log output
	LogLevel logrus.Level

	//Notify all users with these perms by default
	DefaultSubscriberPerms []string

	//Users with these perms are permitted to subscribe to event
	AllowedSubscriberPerms []string
	MetaDataPrototype      interface{}
	MetaDataFactory        EventMetaDataFactory
}

// EventListener models the contract for different event disteners
type EventListener interface {
	Dispatched(event *Event)
}

//Event models an event in the system
type Event struct {
	Key      string
	Session  Session
	MetaData interface{}
}
