package applications

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/production-grid/pgrid-core/pkg/cache"
	"github.com/production-grid/pgrid-core/pkg/config"
	"github.com/production-grid/pgrid-core/pkg/database/relational"
	"github.com/production-grid/pgrid-core/pkg/loaders"
	"github.com/production-grid/pgrid-core/pkg/logging"
)

//SessionCookieName models the correct name for a session cookie
const SessionCookieName = "_pg_session"

//CurrentApplication is a reference to the current application
var CurrentApplication *Application

// Application is the main entry point for productiong grid app.
// Developers configure the application with services and modules,
// then starts it.
type Application struct {
	Name              string
	TagLine           string
	TenantLingo       *TenantLingo
	Modules           []FeatureModule
	SessionStore      cache.Provider
	Cache             cache.Provider
	ConfigLoader      loaders.ResourceLoader
	TemplateLoader    loaders.ResourceLoader
	CoreConfiguration config.CoreConfiguration
	AutoInit          bool //launch the automatic configuration tool if no user
	DefaultAdminUser  *DefaultAdminUser
	Router            *mux.Router
	Server            *http.Server
	EventListeners    []EventListener

	//private stuff
	schemaFiles   []string
	apiRoutes     []APIRoute
	contentRoutes []ContentRoute
	eventDefs     map[string]EventDef
	permissions   []Permission
}

//EventDef looks up an event definition by key
func (app *Application) EventDef(key string) *EventDef {
	def, exists := app.eventDefs[key]
	if exists {
		return &def
	}

	return nil
}

// Start starts the application.
func (app *Application) Start() {

	if app.TenantLingo == nil {
		app.TenantLingo = &TenantLingo{
			TenantSingularKey:     "tenant",
			TenantPluralKey:       "tenants",
			TenantSingularDefault: "Tenant",
			TenantPluralDefault:   "Tenants",
		}
	}

	logging.Infof("Starting %v with Lingo for %v...", app.Name, app.TenantLingo.TenantPluralDefault)

	err := app.initModules()

	if err != nil {
		app.handleStartupError(err)
	}

	err = app.initDatabase()

	if err != nil {
		app.handleStartupError(err)
	}

	err = initRouter(app)

	if err != nil {
		app.handleStartupError(err)
	}

	err = app.notifyAfterAppInit()

	if err != nil {
		app.handleStartupError(err)
	}

	CurrentApplication = app

}

//StartAndListen starts the server and listens for connections
func (app *Application) StartAndListen() {

	app.Start()
	app.Listen()

}

//Listen listens for connections to the http router
func (app *Application) Listen() {

	if app.Server == nil {
		panic("the application has not been started")
	} else {
		logging.Infoln("Starting HTTP Server...")
		logging.Fatal(app.Server.ListenAndServe())
	}
}

func (app *Application) notifyAfterAppInit() error {

	if app.Modules != nil {
		//first loop prenotifies all modules
		for _, mod := range app.Modules {
			err := mod.AfterAppInit(app)
			if err != nil {
				return err
			}
		}
	}
	return nil

}

func (app *Application) initDatabase() error {

	return relational.Init(app.CoreConfiguration.DatabaseConfiguration)

}

func (app *Application) handleStartupError(err error) {
	logging.Errorln("Application startup failed.")
	panic(err)
}

func (app *Application) initModules() error {

	if app.Modules != nil {
		//first loop prenotifies all modules
		for _, mod := range app.Modules {
			err := mod.BeforeAppInit(app)
			if err != nil {
				return err
			}
		}

		//this loop activates the modules
		for _, mod := range app.Modules {
			err := app.initModule(mod)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func (app *Application) addSchemaFilesFrom(mod FeatureModule) error {

	modSchema, err := mod.SchemaFiles(app)

	if err != nil {
		return err
	}

	if modSchema != nil && len(modSchema) > 0 {
		if app.schemaFiles == nil {
			app.schemaFiles = make([]string, 0)
		}
		app.schemaFiles = append(app.schemaFiles, modSchema...)
	}
	return nil

}

func (app *Application) addPermissionsFrom(mod FeatureModule) {

	perms := mod.Permissions(app)

	if perms != nil && len(perms) > 0 {
		if app.permissions == nil {
			app.permissions = make([]Permission, 0)
		}
		app.permissions = append(app.permissions, perms...)
	}

}

func (app *Application) addEventsFrom(mod FeatureModule) {

	events := mod.EventDefs(app)

	if events != nil && len(events) > 0 {
		if app.eventDefs == nil {
			app.eventDefs = make(map[string]EventDef)
		}
		for _, def := range events {
			app.eventDefs[def.Key] = def
		}
	}

}

func (app *Application) initModule(mod FeatureModule) error {
	err := mod.BeforeModuleInit(app)
	if err != nil {
		return err
	}

	logging.Infof("Initializing Module: %v", mod.ID())
	err = app.addSchemaFilesFrom(mod)
	if err != nil {
		return err
	}
	app.addPermissionsFrom(mod)

	err = app.addRoutesFrom(mod)

	if err != nil {
		return err
	}

	app.addEventsFrom(mod)

	return nil

}

func (app *Application) addRoutesFrom(mod FeatureModule) error {

	//API routes
	apiRoutes, err := mod.APIRoutes(app)
	if err != nil {
		return err
	}

	if app.apiRoutes == nil {
		app.apiRoutes = make([]APIRoute, 0)
	}
	for _, route := range apiRoutes {
		//translate module context
		route.Path = "/" + mod.ID() + route.Path
		app.apiRoutes = append(app.apiRoutes, route)
	}

	//TODO: Server rendered content routes

	return nil

}

// PreMigrate runs the post migration database schema changes, if any.
func (app *Application) PreMigrate() {

	logging.Infof("Pre migrating Database Schema for %v\n", app.Name)

	err := app.initModules()

	if err != nil {
		logging.Error(err)
		logging.Errorln("Pre migration failed due to previous errors.")
	}

	err = relational.PreMigrate(app.ConfigLoader, app.CoreConfiguration.DatabaseConfiguration, app.schemaFiles)

	if err != nil {
		logging.Error(err)
		logging.Errorln("Pre migration failed due to previous errors.")
	}

	app.clearTransientState()

}

// PostMigrate runs the post migration database schema changes, if any.
func (app *Application) PostMigrate() {

	logging.Infof("Post migrating Database Schema for %v\n", app.Name)

	err := app.initModules()

	if err != nil {
		logging.Error(err)
		logging.Errorln("Post migration failed due to previous errors.")
	}

	err = relational.PostMigrate(app.ConfigLoader, app.CoreConfiguration.DatabaseConfiguration, app.schemaFiles)

	if err != nil {
		logging.Error(err)
		logging.Errorln("Post migration failed due to previous errors.")
	}

	app.clearTransientState()

}

// AllPermKeys returns all permission codes for the given scope
func (app *Application) AllPermKeys(permScope PermScope) []string {

	results := make([]string, 0)

	for _, perm := range app.permissions {
		if perm.Scope == permScope {
			results = append(results, perm.Key)
		}
	}

	return results

}

func (app *Application) clearTransientState() {

	app.Router = nil
	app.Server = nil
	app.apiRoutes = nil
	app.schemaFiles = nil
	app.permissions = nil
	app.contentRoutes = nil

}
