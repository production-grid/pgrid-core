package security

import (
	"github.com/production-grid/pgrid-core/pkg/applications"
	"github.com/production-grid/pgrid-core/pkg/ids"
)

//Session models an interactive sesion
type Session struct {
	SessionKey           string
	UserID               string
	FirstName            string
	LastName             string
	TenantID             string
	EffectivePermissions []string
}

//InitSession creates an interactive session
func InitSession(user User, tenantID *string) (*Session, error) {

	session := Session{
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
