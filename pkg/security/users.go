package security

import (
	"database/sql"
	"time"

	"github.com/production-grid/pgrid-core/pkg/applications"
	"github.com/production-grid/pgrid-core/pkg/database/relational"
	"github.com/production-grid/pgrid-core/pkg/ids"
)

const tableUsers = "users"

//NewUser is the user factory function
func NewUser() interface{} {
	return &User{}
}

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

//InitSession creates an interactive session
func (user *User) InitSession() (*applications.Session, error) {

	//TODO add user->session index to enforce session fixation rules

	session := applications.Session{
		SessionKey: ids.NewSecureID(),
		UserID:     user.ID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
	}

	err := applications.CurrentApplication.SessionStore.Put(session.SessionKey, session)

	if err != nil {
		return nil, err
	}

	return &session, nil

}

// UserFinder is used to locate users in das system
type UserFinder struct {
}

// FindInterfaceByID returns a user as an empty interface.
func (finder *UserFinder) FindInterfaceByID(dbType string, id string) (interface{}, error) {
	return finder.FindByID(dbType, id)
}

//FindCount returns the number of users in the system
func (finder *UserFinder) FindCount(dbType string) (int, error) {
	return relational.FindCount(dbType, tableUsers, &User{})
}

//FindByEmailAddress returns a user by their email address
func (finder *UserFinder) FindByEmailAddress(dbType string, email string) (*User, error) {
	user := User{}
	err := relational.FindOneByWhereClause(dbType, tableUsers, "where email = $1", &user, email)
	if err != nil {
		if err == relational.ErrNoResults {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
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
