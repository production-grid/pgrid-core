package security

import (
	"time"

	"github.com/production-grid/pgrid-core/pkg/database/relational"
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

	return relational.Save(domain, tableTenants)

}

// Delete deletes a tenant without a transaction.
func (domain *Tenant) Delete() error {
	return relational.SoftDelete(domain, tableTenants)
}
