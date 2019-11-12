package security

import (
	"time"

	"github.com/production-grid/pgrid-core/pkg/applications"
	"github.com/production-grid/pgrid-core/pkg/database/relational"
	"github.com/production-grid/pgrid-core/pkg/util"
)

const tableTenants = "tenants"

//NewTenant is the tenant factory function
func NewTenant() interface{} {
	return &Tenant{}
}

//Tenant models a system tenant
type Tenant struct {
	ID          string    `col:"id"`
	Name        string    `col:"name"`
	TenantType  string    `col:"tenant_type"`
	Tagline     string    `col:"tagline"`
	Subdomain   string    `col:"subdomain"`
	Description string    `col:"description"`
	TimeZone    string    `col:"time_zone"`
	CreatedDate time.Time `col:"created_date"`
}

//TenantDTO models a system tenant in human readable format
type TenantDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	TenantType  string `json:"tenantType"`
	Tagline     string `json:"tagline"`
	Subdomain   string `json:"subdomain"`
	Description string `json:"description"`
	TimeZone    string `json:"timeZone"`
	CreatedDate string `json:"createdDate"`
}

//Identifier returns the id of the dto
func (dto *TenantDTO) Identifier() string {
	return dto.ID
}

// ToDTO converts this domain struct to a UI friendly version.
func (domain *Tenant) ToDTO(session *applications.Session) TenantDTO {

	dto := TenantDTO{}
	dto.ID = domain.ID
	dto.Name = domain.Name
	dto.CreatedDate = util.FormatShortTimeStamp(domain.CreatedDate.In(session.Location()))
	dto.Description = domain.Description
	dto.Subdomain = domain.Subdomain
	dto.Tagline = domain.Tagline
	dto.TenantType = domain.TenantType
	dto.TimeZone = domain.TimeZone
	return dto

}

// TenantFinder is used to locate tenants in the system
type TenantFinder struct {
}

// FindInterfaceByID returns a user as an empty interface.
func (finder *TenantFinder) FindInterfaceByID(dbType string, id string) (interface{}, error) {
	return finder.FindByID(dbType, id)
}

//FindByID locates a tenant by id
func (finder *TenantFinder) FindByID(dbType string, id string) (*Tenant, error) {

	domain := &Tenant{}
	err := relational.FindByID(dbType, tableTenants, id, domain)

	if err != nil {
		return nil, err
	}

	return domain, nil

}

//FindCount returns the number of users in the system
func (finder *TenantFinder) FindCount(dbType string) (int, error) {
	return relational.FindCount(dbType, tableTenants, &Tenant{})
}

//FindAll returns the number of tenants in the system
func (finder *TenantFinder) FindAll(dbType string) ([]interface{}, error) {
	return relational.FindAll(dbType, NewTenant, tableTenants, "order by name")
}

// Save saves a user to the database without a transaction.
func (domain *Tenant) Save() (string, error) {

	if domain.ID == "" {
		domain.CreatedDate = time.Now()
	}

	return relational.Save(domain, tableTenants)

}

// Delete deletes a tenant without a transaction.
func (domain *Tenant) Delete() error {
	return relational.SoftDelete(domain, tableTenants)
}
