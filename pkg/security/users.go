package security

import (
	"database/sql"
	"time"

	"github.com/production-grid/pgrid-core/pkg/applications"
	"github.com/production-grid/pgrid-core/pkg/database/relational"
	"github.com/production-grid/pgrid-core/pkg/ids"
	"github.com/production-grid/pgrid-core/pkg/logging"
	"github.com/production-grid/pgrid-core/pkg/util"
)

const tableUsers = "users"
const tableUserPerms = "user_perms"
const tableUserAdminGroups = "user_admin_groups"

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
	Groups       []string
}

// UserDTO models a user in a display friendly format.
type UserDTO struct {
	ID           string   `json:"id"`
	EMail        string   `json:"email"`
	MobileNumber string   `json:"mobileNumber"`
	FirstName    string   `json:"firstName"`
	MiddleName   string   `json:"middleName"`
	LastName     string   `json:"lastName"`
	IsLocked     bool     `json:"locked"`
	TimeZone     string   `json:"timeZone"`
	LastLogin    string   `json:"lastLogin"`
	RegDate      string   `json:"regDate"`
	Permissions  []string `json:"permissions"`
	Groups       []string `json:"groups"`
}

//Identifier returns the id of the dto
func (dto *UserDTO) Identifier() string {
	return dto.ID
}

// ToDTO converts this domain struct to a UI friendly version.
func (user *User) ToDTO(session *applications.Session) UserDTO {

	dto := UserDTO{}
	dto.ID = user.ID
	dto.EMail = user.EMail
	dto.FirstName = user.FirstName
	dto.MiddleName = user.MiddleName
	dto.LastName = user.LastName
	dto.MobileNumber = user.MobileNumber
	dto.IsLocked = user.IsLocked
	dto.RegDate = util.FormatShortTimeStamp(user.RegDate.In(session.Location()))
	if user.LastLogin != nil {
		dto.LastLogin = util.FormatShortTimeStamp(user.LastLogin.In(session.Location()))
	} else {
		dto.LastLogin = "N/A"
	}
	dto.Permissions = user.Permissions
	dto.Groups = user.Groups

	return dto

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

func (user *User) resolveEffectivePermissions(session *applications.Session) error {

	permMap := make(map[string]bool)

	//start with user level permissions
	if user.Permissions != nil {
		for _, code := range user.Permissions {
			permMap[code] = true
		}
	}

	//add group permission
	if user.Groups != nil {
		perms, err := user.LoadGroupPermissions(relational.REPLICA)
		if err != nil {
			return err
		}
		for _, code := range perms {
			permMap[code] = true
		}
	}

	session.EffectivePermissions = permMap

	return nil

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

//FindAll returns the number of users in the system
func (finder *UserFinder) FindAll(dbType string) ([]interface{}, error) {
	return relational.FindAll(dbType, NewUser, tableUsers, "order by last_name, first_name")
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
	user.Groups, err = user.LoadGroups(dbType)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

//FindByMobileNumber returns a user by their email address
func (finder *UserFinder) FindByMobileNumber(dbType string, mobileNumber string) (*User, error) {
	user := User{}
	err := relational.FindOneByWhereClause(dbType, tableUsers, "where mobile_number = $1", &user, mobileNumber)
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
	user.Groups, err = user.LoadGroups(dbType)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

//FindByLoginID returns a user by their mobile and email
func (finder *UserFinder) FindByLoginID(dbType string, loginID string) (*User, error) {
	user := User{}
	err := relational.FindOneByWhereClause(dbType, tableUsers, "where email = $1 or mobile_number = $1", &user, loginID)
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
	user.Groups, err = user.LoadGroups(dbType)
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

	user.Groups, err = user.LoadGroups(dbType)
	if err != nil {
		return nil, err
	}

	return user, nil

}

//LoadPermissions loads permissions from the database
func (user *User) LoadPermissions(dbType string) ([]string, error) {

	return relational.LoadChildren(dbType, tableUserPerms, "user_id", "permission", user.ID)

}

//LoadGroupPermissions loads permissions from the database
func (user *User) LoadGroupPermissions(dbType string) ([]string, error) {

	db := relational.ResolveDatabase(dbType)

	sql := "select distinct gp.permission from admin_group_perms as gp join user_admin_groups as ug on (gp.admin_group_id = ug.group_id) where ug.user_id = $1"

	logging.Traceln("SQL", sql)

	rows, err := db.Query(sql, user.ID)

	if err != nil {
		return nil, err
	}

	ids := make([]string, 0)

	for rows.Next() {
		var id *string
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, *id)
	}

	return ids, err

}

//LoadGroups loads groups from the database
func (user *User) LoadGroups(dbType string) ([]string, error) {

	return relational.LoadChildren(dbType, tableUserAdminGroups, "user_id", "group_id", user.ID)

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
	err = relational.SaveChildren(tx, tableUserAdminGroups, user, "user_id", "group_id", user.Groups)
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
