package security

import (
	"net/http"
	"strings"
	"time"

	"github.com/production-grid/pgrid-core/pkg/applications"
	"github.com/production-grid/pgrid-core/pkg/crypto"
	"github.com/production-grid/pgrid-core/pkg/database/relational"
	"github.com/production-grid/pgrid-core/pkg/logging"
)

//permission key constants
const (
	PermAdmin       = "security.perm.admin"        //super user - root
	PermImpersonate = "security.perm.impersonate"  //has the ability to impersonate a tenant user
	PermTenantAdmin = "security.perm.tenant.admin" //highest level of access within a tenant scope
	PermTenantUser  = "security.perm.tenant.user"  //basic tenant user
	PermUser        = "security.perm.user"         //basic system user
)

//Module declares the security features for the security module.
type Module struct {
}

//ID returns the unique id of the module
func (mod *Module) ID() string {
	return "security"
}

//APIRoutes declares the API routes provided by this module
func (mod *Module) APIRoutes(app *applications.Application) ([]applications.APIRoute, error) {

	routes := []applications.APIRoute{
		applications.APIRoute{
			Route: applications.Route{
				Path:        "/heartbeat",
				HandlerFunc: GetHeartbeat,
			},
			Method: http.MethodGet,
		},
		applications.APIRoute{
			Route: applications.Route{
				Path:        "/reg",
				HandlerFunc: PostUserReg,
			},
			Method: http.MethodPost,
		},
		applications.APIRoute{
			Route: applications.Route{
				Path:        "/login",
				HandlerFunc: PostLogin,
			},
			Method: http.MethodPost,
		},
	}

	return routes, nil
}

//Permissions returns all permissions defined by this module
func (mod *Module) Permissions(app *applications.Application) []applications.Permission {

	return []applications.Permission{
		applications.Permission{
			Key:         PermAdmin,
			Description: "System Administrator",
			Scope:       applications.PermScopeAdmin,
			GrantedBy:   []string{PermAdmin},
		},
		applications.Permission{
			Key:         PermImpersonate,
			Description: "Impersonate %v User",
			Scope:       applications.PermScopeAdmin,
			GrantedBy:   []string{PermAdmin},
		},
		applications.Permission{
			Key:         PermTenantAdmin,
			Description: "%v Administrator",
			Scope:       applications.PermScopeTenant,
			GrantedBy:   []string{PermAdmin, PermTenantAdmin},
		},
		applications.Permission{
			Key:         PermTenantUser,
			Description: "%v User",
			Scope:       applications.PermScopeTenant,
			GrantedBy:   []string{PermAdmin, PermTenantAdmin},
		},
		applications.Permission{
			Key:         PermUser,
			Description: "User",
			Scope:       applications.PermScopePersonal,
			GrantedBy:   []string{PermAdmin},
		},
	}

}

//BeforeAppInit is called on the module before appliation startup.
func (mod *Module) BeforeAppInit(app *applications.Application) error {

	return nil
}

//AfterAppInit is called on the module after application startup.
func (mod *Module) AfterAppInit(app *applications.Application) error {

	//create admin user if needed
	if app.DefaultAdminUser != nil {

		userFinder := UserFinder{}
		count, err := userFinder.FindCount(relational.PRIMARY)

		if err != nil {
			return err
		}
		if count == 0 {
			logging.Warnln("Creating Default Admin User:", app.DefaultAdminUser.EMail)
			user := User{}
			user.FirstName = app.DefaultAdminUser.FirstName
			user.LastName = app.DefaultAdminUser.LastName
			user.EMail = app.DefaultAdminUser.EMail
			user.RegDate = time.Now()
			hash := crypto.ComputeDoubleHash(app.DefaultAdminUser.Password)
			user.PasswordHash = hash.OuterHash
			user.InnerSalt = hash.InnerSalt
			permCodes := make([]string, 0)
			permCodes = append(permCodes, app.AllPermKeys(applications.PermScopeAdmin)...)
			permCodes = append(permCodes, app.AllPermKeys(applications.PermScopePersonal)...)
			user.Permissions = strings.Join(permCodes, ",")
			user.Save()
		}

	}

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
func (mod *Module) SchemaFiles(app *applications.Application) ([]string, error) {

	return []string{"schema/security.json"}, nil

}
