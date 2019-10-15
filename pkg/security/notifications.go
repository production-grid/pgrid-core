package security

import (
	"github.com/production-grid/pgrid-core/pkg/applications"
	"github.com/production-grid/pgrid-core/pkg/database/relational"
	"github.com/production-grid/pgrid-core/pkg/logging"
)

//SubscriberRecipientListFunc returns a channel of subscribing user id's for each event
func SubscriberRecipientListFunc(recipientChan chan<- applications.Recipient, def *applications.EventDef, event *applications.Event) {

	go fillList(recipientChan, def, event)

}

func fillList(recipientChan chan<- applications.Recipient, def *applications.EventDef, event *applications.Event) {

	userFinder := UserFinder{}

	if def.DefaultSubscriberPerms != nil {
		for _, perm := range def.DefaultSubscriberPerms {
			users, err := userFinder.FindByPermission(relational.REPLICA, perm)
			if err != nil {
				logging.Errorln(err)
			}
			for _, user := range users {
				recipientChan <- user.Recipient()
			}
		}
	}

	close(recipientChan)
}
