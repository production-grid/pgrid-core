package security

import (
	"database/sql"
	"time"

	"github.com/production-grid/pgrid-core/pkg/database/relational"
)

const tableUsers = "users"

// User models a system user.  Permissions on this entity are global
// and generally tenant users would have no top level permissions.
type User struct {
	ID           string     `col:"id"`
	EMail        string     `col:"email"`
	MobileNumber string     `col:"mobile_number"`
	FirstName    string     `col:"first_name"`
	MiddleName   string     `col:"middle_name"`
	LastName     string     `col:"last_name"`
	PasswordHash string     `col:"password_hash"`
	InnerSalt    string     `col:"inner_salt"`
	IsLocked     bool       `col:"is_locked"`
	TimeZone     string     `col:"time_zone"`
	LastLogin    *time.Time `col:"last_login"`
	RegDate      time.Time  `col:"reg_date"`
	Permissions  string     `col:"permissions"`
}

// UserFinder is used to locate users in das system
type UserFinder struct {
}

// FindInterfaceByID returns a user as an empty interface.
func (finder *UserFinder) FindInterfaceByID(dbType string, id string) (interface{}, error) {
	return finder.FindByID(dbType, id)
}

//FindByID locates a user by id
func (finder *UserFinder) FindByID(dbType string, id string) (*User, error) {

	user := &User{}
	err := relational.FindByID(dbType, tableUsers, id, user)

	if err != nil {
		return nil, err
	}

	return user, nil

}

// SaveWithTx saves a user to the database with a transaction context.
func (user *User) SaveWithTx(tx *sql.Tx) (string, error) {
	return relational.SaveWithTx(tx, user, tableUsers)
}

// Save saves a user to the database without a transaction.
func (user *User) Save() (string, error) {
	return relational.Save(user, tableUsers)
}

// Delete deletes a user without a transaction.
func (user *User) Delete() error {
	return relational.SoftDelete(user, tableUsers)
}
