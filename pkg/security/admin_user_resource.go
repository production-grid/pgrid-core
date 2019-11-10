package security

import (
	"net/http"

	"github.com/production-grid/pgrid-core/pkg/applications"
)

//AdminUserResource provides a Level 2 Rest API for dealing with users
type AdminUserResource struct {
}

// Path specifies the Rest resource base path
func (rc *AdminUserResource) Path() string {
	return "/admin/users"
}

// Permissions returns access permission settings for this resource
func (rc *AdminUserResource) Permissions() applications.CrudResourcePermissions {

	permGroup := applications.PermissionGroup{Permission: PermAdmin}

	return applications.CrudResourcePermissions{
		ReadPermissions:   permGroup,
		UpdatePermissions: permGroup,
		DeletePermissions: permGroup,
	}

}

//MetaData returns metadata about the this resource.  Used for rendering the UI.
func (rc *AdminUserResource) MetaData(session applications.Session, req *http.Request) applications.CrudResourceMetaData {

	crud := applications.NewCrudMetaData()
	crud.ListPageTitle = "Admin Users"
	crud.ListPageHelp = "These are internal BlockChyp users."
	crud.NewFormTitle = "Invite User"
	crud.ResourceName = "Admin User"
	crud.ResourceNamePlural = "Admin Users"
	crud.EditFormTitle = "Admin User Settings"
	crud.EditFormHelp = "Manage admin user permissions and settings."
	crud.DeleteEnabled = true
	crud.WithColumn(applications.CrudField{Caption: "E-Mail", ID: "emailAddress"})
	crud.WithColumn(applications.CrudField{Caption: "First Name", ID: "givenName"})
	crud.WithColumn(applications.CrudField{Caption: "Last Name", ID: "surName"})
	crud.WithColumn(applications.CrudField{Caption: "Mobile", ID: "mobileNumber"})
	crud.WithColumn(applications.CrudField{Caption: "Registration Date", ID: "registrationDate"})
	crud.WithColumn(applications.CrudField{Caption: "Status", ID: "status"})

	crud.WithField(applications.CrudField{Caption: "First Name", ID: "givenName", DataType: "string", Required: true, Mutable: true, Min: 2, Max: 64})
	crud.WithField(applications.CrudField{Caption: "Last Name", ID: "surName", DataType: "string", Required: true, Mutable: true, Min: 2, Max: 64})
	crud.WithField(applications.CrudField{Caption: "E-Mail Address", ID: "emailAddress", DataType: "string", Required: true, Mutable: true, Min: 2, Max: 64})
	crud.WithField(applications.CrudField{Caption: "Mobile Number", ID: "mobileNumber", DataType: "string", Required: true, Mutable: true, Min: 2, Max: 64})
	crud.WithField(applications.CrudField{Caption: "Locked", ID: "locked", DataType: "boolean", Required: false, Mutable: true})
	crud.WithField(applications.CrudField{Caption: "Roles", ID: "roles", DataType: "adminRoles", Required: false, Mutable: true})
	crud.WithField(applications.CrudField{Caption: "Registration Date", ID: "registrationDate", DataType: "string", Mutable: false})
	crud.WithField(applications.CrudField{Caption: "Last Login Date", ID: "loginDate", DataType: "string", Mutable: false})
	crud.WithField(applications.CrudField{Caption: "Deleted", ID: "deleted", DataType: "boolean", Mutable: false})
	crud.WithField(applications.CrudField{Caption: "Password Locked", ID: "passwordLocked", DataType: "boolean", Mutable: true})

	return crud

}

// All returns all visible domain objects with their default sorting
func (rc *AdminUserResource) All(session applications.Session, req *http.Request) ([]interface{}, error) {

	return nil, nil

}

//One returns a single domain object
func (rc *AdminUserResource) One(session applications.Session, req *http.Request, id string) (interface{}, error) {
	return nil, nil
}

//Update creates or updates an entity, returning the entity after it's been created or upudated
func (rc *AdminUserResource) Update(session applications.Session, req *http.Request, dto interface{}, domain interface{}) (interface{}, error) {
	return domain, nil
}

//Delete deletes an entity from the system, returning true if the operation succeeded
func (rc *AdminUserResource) Delete(session applications.Session, req *http.Request, id string) (bool, error) {
	return false, nil
}
