package notifications

import "github.com/production-grid/pgrid-core/pkg/applications"

//SendGridTransport sends email using the send grid.
type SendGridTransport struct {
}

//Type returns the transport protocol type
func (transport *SendGridTransport) Type() applications.TransportType {
	return applications.TransportEmail
}
