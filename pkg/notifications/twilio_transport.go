package notifications

import (
	"github.com/production-grid/pgrid-core/pkg/applications"
	"github.com/production-grid/pgrid-core/pkg/config"
	"github.com/production-grid/pgrid-core/pkg/logging"

	"github.com/sfreiberg/gotwilio"
)

//TwilioTransport sends email using the send grid.
type TwilioTransport struct {
	Config config.TwilioConfig
}

//Type returns the transport protocol type
func (transport *TwilioTransport) Type() applications.TransportType {
	return applications.TransportSMS
}

//Description returns the narrative description
func (transport *TwilioTransport) Description() string {
	return "SendGrid Transport"
}

//Send pretends to send output to the recipient
func (transport *TwilioTransport) Send(def *applications.EventDef, event *applications.Event, rcpt *applications.Recipient) {

	smsTemplate, err := LoadTemplate(TemplateTypeSMS, def, event, rcpt)
	if err != nil {
		logging.Errorln(err)
		return
	}
	if smsTemplate == "" {
		logging.Errorln("unable to load sms template")
		return
	}

	twilio := gotwilio.NewTwilioClient(transport.Config.SID, transport.Config.AuthToken)

	from := "+" + transport.Config.Number
	to := "+1" + rcpt.SMSNumber
	twilio.SendSMS(from, to, smsTemplate, "", "")

}
