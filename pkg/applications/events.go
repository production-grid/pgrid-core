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
	MetaDataPrototype      interface{} //used as merge data for the template UI
	MetaDataFactory        EventMetaDataFactory

	//Determines how to generate a list of recipients
	//If null, the subscriber method is used
	RecipientListFunc RecipientListFunc
}

//TransportType is used to enumerate transport types
type TransportType string

//Enumerates notification transport types
const (
	TransportSMS   TransportType = "sms"
	TransportEmail TransportType = "email"
	TransportPush  TransportType = "push"
)

// Recipient models an event recipient
type Recipient struct {
	FirstName     string
	LastName      string
	TenantID      string
	EMailAddress  string
	SMSNumber     string
	TransportType TransportType //if known
	//TODO add ids for push notifications when mobile apps are added
}

// EventListener models the contract for different event disteners
type EventListener interface {
	Dispatched(def *EventDef, event *Event)
}

//RecipientListFunc returns a channel of user id's to be notified.
//The logic might be different for different events.  For example, an order
//notification would be sent to the customer as opposed to security events
//which would go to system administrators.
type RecipientListFunc func(recipientChan chan<- Recipient, def *EventDef, event *Event)

//Event models an event in the system
type Event struct {
	Key      string
	Session  Session
	MetaData interface{}
}
