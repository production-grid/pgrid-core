package security

import (
	"database/sql"
	"time"

	"github.com/production-grid/pgrid-core/pkg/applications"
	"github.com/production-grid/pgrid-core/pkg/database/relational"
	"github.com/production-grid/pgrid-core/pkg/ids"
)

const tableUsers = "users"
const tableUserPerms = "user_perms"

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
	Permissions  []string
}

//Recipient returns the notification recipient struct for this user
func (user *User) Recipient() applications.Recipient {

	return applications.Recipient{
		EMailAddress: user.EMail,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		SMSNumber:    user.MobileNumber,
	}

}

//InitSession creates an interactive session
func (user *User) InitSession(currentSession *applications.Session) (*applications.Session, error) {

	//TODO add user->session index to enforce session fixation rules

	session := applications.Session{
		SessionKey: ids.NewSecureID(),
		UserID:     user.ID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		TenantID:   currentSession.TenantID,
	}

	user.resolveEffectivePermissions(&session)

	err := applications.CurrentApplication.SessionStore.Put(session.SessionKey, session)

	if err != nil {
		return nil, err
	}

	return &session, nil

}

func (user *User) resolveEffectivePermissions(session *applications.Session) {

	permMap := make(map[string]bool)

	//start with user level permissions
	if user.Permissions != nil {
		for _, code := range user.Permissions {
			permMap[code] = true
		}
	}

	session.EffectivePermissions = permMap

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

//FindByPermission returns a user granted the given permission
func (finder *UserFinder) FindByPermission(dbType string, perm string) ([]User, error) {

	rawResults, err := relational.FindByWhereClause(dbType, NewUser, tableUsers, "where is_locked = false and id in (select user_id from user_perms where permission = $1) ", perm)
	if err != nil {
		return nil, err
	}
	results := make([]User, len(rawResults))
	for idx, rawResult := range rawResults {
		results[idx] = *(rawResult.(*User))
	}

	return results, nil
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
	user.Permissions, err = user.LoadPermissions(dbType)
	if err != nil {
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

	user.Permissions, err = user.LoadPermissions(dbType)
	if err != nil {
		return nil, err
	}
	return user, nil

}

//LoadPermissions loads permissions from the database
func (user *User) LoadPermissions(dbType string) ([]string, error) {

	return relational.LoadChildren(dbType, tableUserPerms, "user_id", "permission", user.ID)

}

// SaveWithTx saves a user to the database with a transaction context.
func (user *User) SaveWithTx(tx *sql.Tx) (string, error) {
	id, err := relational.SaveWithTx(tx, user, tableUsers)
	if err != nil {
		tx.Rollback()
		return "", err
	}
	err = relational.SaveChildren(tx, tableUserPerms, user, "user_id", "permission", user.Permissions)
	if err != nil {
		tx.Rollback()
		return "", err
	}
	return id, err
}

// Save saves a user to the database without a transaction.
func (user *User) Save() (string, error) {

	tx, err := relational.NewWritableTx()
	if err != nil {
		return "", err
	}
	id, err := user.SaveWithTx(tx)
	if err != nil {
		tx.Rollback()
	}
	tx.Commit()
	return id, err
}

// Delete deletes a user without a transaction.
func (user *User) Delete() error {
	return relational.SoftDelete(user, tableUsers)
}
