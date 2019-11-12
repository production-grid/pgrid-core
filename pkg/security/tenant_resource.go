package security

import (
	"net/http"

	"github.com/production-grid/pgrid-core/pkg/applications"
	"github.com/production-grid/pgrid-core/pkg/database/relational"
)

//TenantResource provides a Level 2 Rest API for dealing with tenants
type TenantResource struct {
}

// Path specifies the Rest resource base path
func (rc *TenantResource) Path() string {
	return "/tenants"
}

// Permissions returns access permission settings for this resource
func (rc *TenantResource) Permissions() applications.CrudResourcePermissions {

	permGroup := applications.PermissionGroup{Permission: PermAdmin}

	return applications.CrudResourcePermissions{
		ReadPermissions:   permGroup,
		UpdatePermissions: permGroup,
		DeletePermissions: permGroup,
	}

}

//ToDTO copies the domain values to the properly formatted display values
func (rc *TenantResource) ToDTO(session *applications.Session, req *http.Request, from interface{}) (interface{}, error) {

	domain := from.(*Tenant)
	return domain.ToDTO(session), nil

}

//FromDTO reads values from a dto and copies them to the domain object, usually in response to form submission
func (rc *TenantResource) FromDTO(session *applications.Session, req *http.Request, from interface{}, to interface{}) (interface{}, error) {

	dto := from.(*TenantDTO)
	domain := to.(*Tenant)
	domain.Name = dto.Name
	domain.Description = dto.Description
	domain.Subdomain = dto.Subdomain
	domain.Tagline = dto.Tagline
	domain.TenantType = dto.TenantType
	domain.TimeZone = dto.TimeZone

	return domain, nil
}

//NewDTO creates a new instance of the dto struct.
func (rc *TenantResource) NewDTO(session *applications.Session, req *http.Request) interface{} {
	return &TenantDTO{}
}

//NewDomain creates a new instance of the domain struct.
func (rc *TenantResource) NewDomain(session *applications.Session, req *http.Request) interface{} {
	return &Tenant{}
}

//MetaData returns metadata about the this resource.  Used for rendering the UI.
func (rc *TenantResource) MetaData(session *applications.Session, req *http.Request) applications.CrudResourceMetaData {

	lingo := applications.CurrentApplication.TenantLingo

	crud := applications.NewCrudMetaData()
	crud.ListPageTitle = lingo.TenantPlural
	crud.ListPageHelp = "These are all the " + lingo.TenantPlural + " currently in the system."
	crud.NewFormTitle = "Add " + lingo.TenantSingular
	crud.ResourceName = lingo.TenantSingular
	crud.ResourceNamePlural = lingo.TenantPlural
	crud.EditFormTitle = lingo.TenantSingular + " Settings"
	crud.EditFormHelp = "Manage " + lingo.TenantPlural
	crud.DeleteEnabled = true
	crud.WithColumn(applications.CrudField{Caption: "Name", ID: "name"})
	crud.WithColumn(applications.CrudField{Caption: "Type", ID: "tenantType"})
	crud.WithColumn(applications.CrudField{Caption: "Subdomain", ID: "subdomain"})
	crud.WithColumn(applications.CrudField{Caption: "Time Zone", ID: "timeZone"})
	crud.WithColumn(applications.CrudField{Caption: "Created Date", ID: "createdDate"})

	crud.WithField(applications.CrudField{Caption: "Name", ID: "name", DataType: "string", Required: true, Mutable: true, Min: 2, Max: 64})
	crud.WithField(applications.CrudField{Caption: "Date Created", ID: "createdDate", DataType: "string", Mutable: false, HideIfEmpty: true})
	crud.WithField(applications.CrudField{Caption: "Type", ID: "tenantType", DataType: "tenant-type", Required: true, Mutable: true})
	crud.WithField(applications.CrudField{
		Caption:  "Subdomain",
		ID:       "subdomain",
		DataType: "subdomain",
		Required: false,
		Mutable:  true,
		Min:      2,
		Max:      32,
		Help:     "Unique subdomain for web sites without their own domain names.",
	})
	crud.WithField(applications.CrudField{
		Caption:  "Time Zone",
		ID:       "timeZone",
		DataType: "timeZone",
		Required: false,
		Mutable:  true,
		Help:     "Your local time zone.",
	})
	crud.WithField(applications.CrudField{
		Caption:     "Tagline",
		ID:          "tagline",
		DataType:    "string",
		Required:    false,
		Mutable:     true,
		Min:         2,
		Max:         64,
		ColumnStyle: "col-md-12",
		Help:        "A slogan or tagline that will be displayed on the web site.",
	})
	crud.WithField(applications.CrudField{
		Caption:     "Description",
		ID:          "description",
		DataType:    "textarea",
		Required:    false,
		Mutable:     true,
		ColumnStyle: "col-md-12",
		Help:        "A longer description that will be used for search engine optimization and auto-generated web content.",
	})

	return crud

}

// All returns all visible domain objects with their default sorting
func (rc *TenantResource) All(session *applications.Session, req *http.Request) ([]interface{}, error) {

	finder := TenantFinder{}

	return finder.FindAll(relational.REPLICA)

}

//One returns a single domain object
func (rc *TenantResource) One(session *applications.Session, req *http.Request, id string) (interface{}, error) {

	finder := TenantFinder{}

	return finder.FindByID(relational.REPLICA, id)

}

//Update creates or updates an entity, returning the entity after it's been created or upudated
func (rc *TenantResource) Update(session *applications.Session, req *http.Request, dto interface{}, domain interface{}) (interface{}, error) {

	fullDomain := domain.(*Tenant)
	fullDomain.Save()

	return fullDomain, nil
}

//Delete deletes an entity from the system, returning true if the operation succeeded
func (rc *TenantResource) Delete(session *applications.Session, req *http.Request, domain interface{}) (bool, error) {

	fullDomain := domain.(*Tenant)
	err := fullDomain.Delete()

	return true, err
}
