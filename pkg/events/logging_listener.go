package events

import (
	"encoding/json"
	"fmt"

	"github.com/production-grid/pgrid-core/pkg/applications"
	"github.com/production-grid/pgrid-core/pkg/logging"
	"github.com/sirupsen/logrus"
)

//LoggingEventListener logs events
type LoggingEventListener struct {
}

//Dispatched is called whenever an event id processed
func (listener *LoggingEventListener) Dispatched(event *applications.Event) {

	def, ok := applications.CurrentApplication.EventDefs[event.Key]

	if !ok {
		logging.Warnf("unregistered event type: %v\n", event.Key)
		return
	}

	switch def.LogLevel {
	case logrus.DebugLevel:
		logging.Debugln(listener.formatLogLine(&def, event))
	case logrus.InfoLevel:
		logging.Infoln(listener.formatLogLine(&def, event))
	case logrus.WarnLevel:
		logging.Warnln(listener.formatLogLine(&def, event))
	case logrus.ErrorLevel:
		logging.Warnln(listener.formatLogLine(&def, event))
	case logrus.FatalLevel:
		logging.Fatalln(listener.formatLogLine(&def, event))
	default:
		logging.Traceln(listener.formatLogLine(&def, event))
	}

}

func (listener *LoggingEventListener) formatLogLine(def *applications.EventDef, event *applications.Event) string {

	md, err := json.Marshal(event.MetaData)

	if err != nil {
		return fmt.Sprintf("error deserializing event meta data for event: %v\n", event.Key)
	}

	return fmt.Sprintf(
		"Event: %v, Session: %v, UserID: %v TenantID: %v, MetaData: %v\n",
		event.Key,
		event.Session.SessionKey,
		event.Session.UserID,
		event.Session.TenantID,
		string(md),
	)

}
