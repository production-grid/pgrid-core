package security

import (
	"os"

	"github.com/production-grid/pgrid-core/pkg/applications"
)

//Module declares the security features for the security module.
type Module struct {
}

//Name returns the name of the module
func (mod *Module) Name() string {
	return "security"
}

//BeforeAppInit is called on the module before appliation startup.
func (mod *Module) BeforeAppInit(app *applications.Application) error {

	return nil
}

//AfterAppInit is called on the module after application startup.
func (mod *Module) AfterAppInit(app *applications.Application) error {

	return nil
}

//BeforeModuleInit is called on the module before module startup.
func (mod *Module) BeforeModuleInit(app *applications.Application) error {

	return nil
}

//AfterModuleInit is called on the module after module startup.
func (mod *Module) AfterModuleInit(app *applications.Application) error {

	return nil
}

//SchemaFiles returns the database schema configuration files for this module.
func (mod *Module) SchemaFiles(app *applications.Application) ([]os.File, error) {

	return nil, nil
}
