package main

import (
	"github.com/production-grid/pgrid-core/pkg/applications"
	"github.com/production-grid/pgrid-core/pkg/cache"
	"github.com/production-grid/pgrid-core/pkg/config"
	"github.com/production-grid/pgrid-core/pkg/events"
	"github.com/production-grid/pgrid-core/pkg/loaders"
	"github.com/production-grid/pgrid-core/pkg/notifications"
	"github.com/production-grid/pgrid-core/pkg/security"
)

func main() {

	rcLoader := &loaders.FileResourceLoader{}
	configLoader := &loaders.LocalConfigLoader{}

	coreConfig, err := config.LoadCore(configLoader, "pgrid/pgrid_dev.yml")

	if err != nil {
		panic(err)
	}

	lingo := applications.TenantLingo{
		TenantSingular: "Arts Organization",
		TenantPlural:   "Arts Organizations",
		Types: []applications.TenantType{
			applications.TenantType{
				Singular: "Theatre Company",
				Plural:   "Theatre Companies",
			},
			applications.TenantType{
				Singular: "School Drama Program",
				Plural:   "School Drama Programs",
			},
			applications.TenantType{
				Singular: "School Music Program",
				Plural:   "School Music Programs",
			},
			applications.TenantType{
				Singular: "Dance Company",
				Plural:   "Dance Companies",
			},
			applications.TenantType{
				Singular: "Opera Company",
				Plural:   "Opera Companies",
			},
			applications.TenantType{
				Singular: "Choir",
				Plural:   "Choirs",
			},
			applications.TenantType{
				Singular: "Symphony",
				Plural:   "Symphonies",
			},
			applications.TenantType{
				Singular: "Other Arts Group",
				Plural:   "Other Arts Groups",
			},
		},
	}

	localCache := cache.LocalCache{}

	defaultUser := applications.DefaultAdminUser{
		FirstName: "System",
		LastName:  "Administrator",
		EMail:     "jpayne@productiongrid.com",
		Password:  "test123",
	}

	app := applications.Application{
		CoreConfiguration: *coreConfig,
		TenantLingo:       &lingo,
		SessionStore:      &localCache,
		Cache:             &localCache,
		DefaultAdminUser:  &defaultUser,
		Name:              "Production Grid",
		TagLine:           "Connecting the Performing Arts",
		ConfigLoader:      rcLoader,
		TemplateLoader:    rcLoader,
		EventListeners: []applications.EventListener{
			&events.LoggingEventListener{},
			&notifications.NotifyingEventListener{
				DefaultRecipientListFunc: security.SubscriberRecipientListFunc,
				Transports: []notifications.Transport{
					&notifications.TwilioTransport{
						Config: coreConfig.Notifications.Twilio,
					},
					&notifications.SendGridTransport{
						Config: coreConfig.Notifications.SendGrid,
					},
				},
			},
		},
		Modules: []applications.FeatureModule{
			&security.Module{},
		},
	}

	app.PreMigrate()
	app.PostMigrate()
	app.StartAndListen()

}
