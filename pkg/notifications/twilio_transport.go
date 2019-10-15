package notifications

import "github.com/production-grid/pgrid-core/pkg/applications"

//TwilioTransport sends email using the send grid.
type TwilioTransport struct {
}

//Type returns the transport protocol type
func (transport *TwilioTransport) Type() applications.TransportType {
	return applications.TransportSMS
}
