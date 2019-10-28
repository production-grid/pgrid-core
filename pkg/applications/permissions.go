package applications

// PermScope identifies the enforcement and grant scope of a permission
type PermScope int

// current list of permission scopes
const (
	PermScopeAdmin    PermScope = iota //admin level permission mapped at the user level or user group level
	PermScopePersonal                  //non admin permisions that are still scoped at the user level (personal permissions)
	PermScopeTenant                    //permissions that only make sense in the context of a tenant
)

// Permission models a permission that might be contributed to the application by a module
type Permission struct {
	//identifier and resource bundle key (once we get to that)
	Key string

	//description in the event that it can't be found in a rc bundle
	Description string

	//scope of the permission
	Scope PermScope

	//reference to another permissions, one of which must be granted in order to grant this permision
	GrantedBy []string
}

//DefaultAdminUser models an optional method of adding a default user during a cold bootstrap
//should not be used for production builds
type DefaultAdminUser struct {
	FirstName string
	LastName  string
	EMail     string
	Password  string
}
