package security

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/production-grid/pgrid-core/pkg/applications"
	"github.com/production-grid/pgrid-core/pkg/database/relational"
	"github.com/production-grid/pgrid-core/pkg/httputils"
)

//GetAdminPerms returns session meta data
func GetAdminPerms(session applications.Session, w http.ResponseWriter, req *http.Request) {

	perms := applications.CurrentApplication.AllPermissions(applications.PermScopeAll)

	lingo := applications.CurrentApplication.TenantLingo

	dtos := make([]interface{}, len(perms))

	for idx, perm := range perms {
		dto := PermissionDTO{
			Key: perm.Key,
		}
		if strings.Contains(perm.Description, "%v") {
			dto.Description = fmt.Sprintf(perm.Description, lingo.TenantSingular)
		} else {
			dto.Description = perm.Description
		}

		dtos[idx] = dto
	}

	//TODO add sorting

	results := applications.Results{
		VisibleResults:   dtos,
		TotalResultCount: len(dtos),
	}

	httputils.SendJSON(results, w)

}

//AdminGroupResource provides a Level 2 Rest API for dealing with admin groups
type AdminGroupResource struct {
}

// Path specifies the Rest resource base path
func (rc *AdminGroupResource) Path() string {
	return "/admin/groups"
}

// Permissions returns access permission settings for this resource
func (rc *AdminGroupResource) Permissions() applications.CrudResourcePermissions {

	permGroup := applications.PermissionGroup{Permission: PermAdmin}

	return applications.CrudResourcePermissions{
		ReadPermissions:   permGroup,
		UpdatePermissions: permGroup,
		DeletePermissions: permGroup,
	}

}

//ToDTO copies the domain values to the properly formatted display values
func (rc *AdminGroupResource) ToDTO(session *applications.Session, req *http.Request, from interface{}) (interface{}, error) {

	domain := from.(*AdminGroup)
	return domain.ToDTO(session), nil

}

//FromDTO reads values from a dto and copies them to the domain object, usually in response to form submission
func (rc *AdminGroupResource) FromDTO(session *applications.Session, req *http.Request, from interface{}, to interface{}) (interface{}, error) {

	dto := from.(*AdminGroupDTO)
	domain := to.(*AdminGroup)
	domain.Name = dto.Name
	domain.Permissions = dto.Permissions

	return domain, nil
}

//NewDTO creates a new instance of the dto struct.
func (rc *AdminGroupResource) NewDTO(session *applications.Session, req *http.Request) interface{} {
	return &AdminGroupDTO{}
}

//NewDomain creates a new instance of the domain struct.
func (rc *AdminGroupResource) NewDomain(session *applications.Session, req *http.Request) interface{} {
	return &AdminGroup{}
}

//MetaData returns metadata about the this resource.  Used for rendering the UI.
func (rc *AdminGroupResource) MetaData(session *applications.Session, req *http.Request) applications.CrudResourceMetaData {

	crud := applications.NewCrudMetaData()
	crud.ListPageTitle = "Admin Security Groups"
	crud.ListPageHelp = "These are all the admin security groups currently in the system."
	crud.NewFormTitle = "Add Group"
	crud.ResourceName = "Admin Security Group"
	crud.ResourceNamePlural = "Admin Security Groups"
	crud.EditFormTitle = "Manage Admin Security Group"
	crud.EditFormHelp = "Manage Admin Security Groups"
	crud.DeleteEnabled = true
	crud.WithColumn(applications.CrudField{Caption: "Name", ID: "name"})
	crud.WithColumn(applications.CrudField{Caption: "Created Date", ID: "createdDate"})

	crud.WithField(applications.CrudField{Caption: "Name", ID: "name", DataType: "string", Required: true, Mutable: true, Min: 2, Max: 64})
	crud.WithField(applications.CrudField{Caption: "Date Created", ID: "createdDate", DataType: "string", Mutable: false, HideIfEmpty: true})
	crud.WithField(applications.CrudField{
		Caption:     "Permissions",
		ID:          "permissions",
		DataType:    "admin-permissions",
		Required:    false,
		Mutable:     true,
		ColumnStyle: "col-md-12",
	})

	return crud

}

// All returns all visible domain objects with their default sorting
func (rc *AdminGroupResource) All(session *applications.Session, req *http.Request) ([]interface{}, error) {

	finder := AdminGroupFinder{}

	return finder.FindAll(relational.REPLICA)

}

//One returns a single domain object
func (rc *AdminGroupResource) One(session *applications.Session, req *http.Request, id string) (interface{}, error) {

	finder := AdminGroupFinder{}

	return finder.FindByID(relational.REPLICA, id)

}

//Update creates or updates an entity, returning the entity after it's been created or upudated
func (rc *AdminGroupResource) Update(session *applications.Session, req *http.Request, dto interface{}, domain interface{}) (interface{}, error) {

	fullDomain := domain.(*AdminGroup)
	fullDomain.Save()

	return fullDomain, nil
}

//Delete deletes an entity from the system, returning true if the operation succeeded
func (rc *AdminGroupResource) Delete(session *applications.Session, req *http.Request, domain interface{}) (bool, error) {

	fullDomain := domain.(*AdminGroup)
	err := fullDomain.Delete()

	return true, err
}
