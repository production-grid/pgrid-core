package notifications

import (
	"fmt"

	"github.com/production-grid/pgrid-core/pkg/applications"
	"github.com/production-grid/pgrid-core/pkg/config"
	"github.com/production-grid/pgrid-core/pkg/logging"

	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

//SendGridTransport sends email using the send grid.
type SendGridTransport struct {
	Config config.SendGridConfig
}

//Type returns the transport protocol type
func (transport *SendGridTransport) Type() applications.TransportType {
	return applications.TransportEmail
}

//Description returns the narrative description
func (transport *SendGridTransport) Description() string {
	return "SendGrid Transport"
}

//Send pretends to send output to the recipient
func (transport *SendGridTransport) Send(def *applications.EventDef, event *applications.Event, rcpt *applications.Recipient) {

	subjectTemplate, err := LoadTemplate(TemplateTypeSubject, def, event, rcpt)
	if err != nil {
		logging.Errorln(err)
		return
	}
	if subjectTemplate == "" {
		logging.Errorln("unable to load subject template")
		return
	}

	plainBodyTemplate, err := LoadTemplate(TemplateTypePlainBody, def, event, rcpt)
	if err != nil {
		logging.Errorln(err)
		return
	}
	if plainBodyTemplate == "" {
		logging.Errorln("unable to load plain body template")
		return
	}

	htmlBodyTemplate, err := LoadTemplate(TemplateTypeHTMLBody, def, event, rcpt)
	if err != nil {
		logging.Errorln(err)
		return
	}
	if htmlBodyTemplate == "" {
		logging.Errorln("unable to load html body template")
		return
	}

	client := sendgrid.NewSendClient(transport.Config.APIKey)
	from := mail.NewEmail(transport.Config.SenderName, transport.Config.SenderEMail)
	to := mail.NewEmail(rcpt.FirstName+" "+rcpt.LastName, rcpt.EMailAddress)
	message := mail.NewSingleEmail(from, subjectTemplate, to, plainBodyTemplate, htmlBodyTemplate)
	response, err := client.Send(message)
	if err != nil {
		logging.Errorln("sendgrid transport: ", err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
		logging.Infof("SendGrid Transport Processed Event: %v for email: %v", event.Key, rcpt.EMailAddress)
	}

}
