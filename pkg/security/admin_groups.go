package security

import (
	"database/sql"
	"time"

	"github.com/production-grid/pgrid-core/pkg/applications"
	"github.com/production-grid/pgrid-core/pkg/database/relational"
	"github.com/production-grid/pgrid-core/pkg/util"
)

const tableAdminGroups = "admin_groups"
const tableAdminGroupPerms = "admin_group_perms"

//NewAdminGroup is the admin group factory function
func NewAdminGroup() interface{} {
	return &AdminGroup{}
}

//AdminGroup models admin level user groups.
//AdminGroup permissions are granted across all tenant scopes.
type AdminGroup struct {
	ID          string    `col:"id"`
	Name        string    `col:"name"`
	CreatedDate time.Time `col:"created_date"`
	Permissions []string
}

//AdminGroupDTO models admin level user groups.
type AdminGroupDTO struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	CreatedDate string   `json:"createdDate"`
	Permissions []string `json:"permissions"`
}

//Identifier returns the id of the dto
func (dto *AdminGroupDTO) Identifier() string {
	return dto.ID
}

// ToDTO converts this domain struct to a UI friendly version.
func (domain *AdminGroup) ToDTO(session *applications.Session) AdminGroupDTO {

	dto := AdminGroupDTO{}
	dto.ID = domain.ID
	dto.Name = domain.Name
	dto.CreatedDate = util.FormatShortTimeStamp(domain.CreatedDate.In(session.Location()))
	dto.Permissions = domain.Permissions
	return dto

}

//AdminGroupFinder is used to find admin groups.
type AdminGroupFinder struct {
}

// FindInterfaceByID returns an admin group as an empty interface.
func (finder *AdminGroupFinder) FindInterfaceByID(dbType string, id string) (interface{}, error) {
	return finder.FindByID(dbType, id)
}

//FindByID locates a tenant by id
func (finder *AdminGroupFinder) FindByID(dbType string, id string) (*AdminGroup, error) {

	domain := &AdminGroup{}
	err := relational.FindByID(dbType, tableAdminGroups, id, domain)

	if err != nil {
		return nil, err
	}

	domain.Permissions, err = domain.LoadPermissions(dbType)
	if err != nil {
		return nil, err
	}

	return domain, nil

}

//FindCount returns the number of users in the system
func (finder *AdminGroupFinder) FindCount(dbType string) (int, error) {
	return relational.FindCount(dbType, tableAdminGroups, &AdminGroup{})
}

//FindAll returns the number of tenants in the system
func (finder *AdminGroupFinder) FindAll(dbType string) ([]interface{}, error) {
	return relational.FindAll(dbType, NewAdminGroup, tableAdminGroups, "order by name")
}

//LoadPermissions loads permissions from the database
func (domain *AdminGroup) LoadPermissions(dbType string) ([]string, error) {

	return relational.LoadChildren(dbType, tableAdminGroupPerms, "admin_group_id", "permission", domain.ID)

}

// SaveWithTx saves a user to the database with a transaction context.
func (domain *AdminGroup) SaveWithTx(tx *sql.Tx) (string, error) {

	if domain.ID == "" {
		domain.CreatedDate = time.Now()
	}

	id, err := relational.SaveWithTx(tx, domain, tableAdminGroups)
	if err != nil {
		tx.Rollback()
		return "", err
	}
	err = relational.SaveChildren(tx, tableAdminGroupPerms, domain, "admin_group_id", "permission", domain.Permissions)
	if err != nil {
		tx.Rollback()
		return "", err
	}
	return id, err
}

// Save saves a user to the database without a transaction.
func (domain *AdminGroup) Save() (string, error) {

	tx, err := relational.NewWritableTx()
	if err != nil {
		return "", err
	}
	id, err := domain.SaveWithTx(tx)
	if err != nil {
		tx.Rollback()
	}
	tx.Commit()
	return id, err

}

// Delete deletes a tenant without a transaction.
func (domain *AdminGroup) Delete() error {
	return relational.SoftDelete(domain, tableAdminGroups)
}
