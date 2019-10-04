package applications

import (
	"github.com/production-grid/pgrid-core/pkg/config"
	"github.com/production-grid/pgrid-core/pkg/database/relational"
	"github.com/production-grid/pgrid-core/pkg/loaders"
	"github.com/production-grid/pgrid-core/pkg/logging"
)

// Application is the main entry point for productiong grid app.
// Developers configure the application with services and modules,
// then starts it.
type Application struct {
	Name              string
	Modules           []FeatureModule
	SchemaFiles       []string
	ConfigLoader      loaders.ResourceLoader
	CoreConfiguration config.CoreConfiguration
}

// Start starts the application.
func (app *Application) Start() {

	logging.Infof("Starting %v...", app.Name)

	err := app.initModules()

	if err != nil {
		app.handleStartupError(err)
	}

	err = app.initDatabase()

	if err != nil {
		app.handleStartupError(err)
	}

}

func (app *Application) initDatabase() error {
	return nil
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

func (app *Application) initModule(mod FeatureModule) error {
	err := mod.BeforeModuleInit(app)
	if err != nil {
		return err
	}

	logging.Infof("Initializing Module: %v", mod.Name())

	modSchema, err := mod.SchemaFiles(app)

	if err != nil {
		return err
	}

	if modSchema != nil && len(modSchema) > 0 {
		if app.SchemaFiles == nil {
			app.SchemaFiles = make([]string, 0)
		}
		app.SchemaFiles = append(app.SchemaFiles, modSchema...)
	}

	err = mod.AfterModuleInit(app)
	if err != nil {
		return err
	}

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

	err = relational.PreMigrate(app.ConfigLoader, app.CoreConfiguration.DatabaseConfiguration, app.SchemaFiles)

	if err != nil {
		logging.Error(err)
		logging.Errorln("Pre migration failed due to previous errors.")
	}

}

// PostMigrate runs the post migration database schema changes, if any.
func (app *Application) PostMigrate() {

	logging.Infof("Post migrating Database Schema for %v\n", app.Name)

	err := app.initModules()

	if err != nil {
		logging.Error(err)
		logging.Errorln("Post migration failed due to previous errors.")
	}

	err = relational.PostMigrate(app.ConfigLoader, app.CoreConfiguration.DatabaseConfiguration, app.SchemaFiles)

	if err != nil {
		logging.Error(err)
		logging.Errorln("Post migration failed due to previous errors.")
	}

}
