package notifications

import (
	"github.com/production-grid/pgrid-core/pkg/applications"
	"github.com/production-grid/pgrid-core/pkg/logging"
)

// Transport defines the contract for a notification transport
type Transport interface {
	Type() applications.TransportType
	Description() string
	Send(def *applications.EventDef, event *applications.Event, rcpt *applications.Recipient)
}

// MockTransport implements a simple noop notifiation transport for testing
type MockTransport struct {
}

//Type returns the transport protocol type
func (transport *MockTransport) Type() applications.TransportType {
	return applications.TransportEmail
}

//Description returns the narrative description
func (transport *MockTransport) Description() string {
	return "Mock Transport"
}

//Send pretends to send output to the recipient
func (transport *MockTransport) Send(def *applications.EventDef, event *applications.Event, rcpt *applications.Recipient) {

	logging.Infof("Mock Transport Received Event: %v for email: %v", event.Key, rcpt.EMailAddress)

}
