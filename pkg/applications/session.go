package applications

import (
	"net/http"
	"time"

	"github.com/production-grid/pgrid-core/pkg/logging"
)

//default time zone is the theatre capitol of the world
const DefaultTimeZone = "America/New_York"

//Session models an interactive sesion
type Session struct {
	SessionKey           string
	UserID               string
	FirstName            string
	LastName             string
	TenantID             string
	TimeZone             string
	EffectivePermissions map[string]bool //because maps are faster lookups
}

//Location returns the timezone/location for a given session.
func (session *Session) Location() *time.Location {

	if session.TimeZone != "" {
		loc, err := time.LoadLocation(session.TimeZone)
		if err != nil {
			logging.Warnf("Unable to load time zone: %v\n", session.TimeZone)
		} else {
			return loc
		}
	}

	loc, err := time.LoadLocation(DefaultTimeZone)

	if err != nil {
		//this should never happen
		panic(err)
	}

	return loc

}

//HasPermission returns true if session has the effective permission
func (session *Session) HasPermission(permCode string) bool {

	if session.EffectivePermissions == nil {
		return false
	}

	_, granted := session.EffectivePermissions[permCode]

	return granted

}

func findSessionCookie(req *http.Request) *http.Cookie {

	for _, cookie := range req.Cookies() {
		if cookie.Name == SessionCookieName {
			return cookie
		}
	}

	return nil

}

func resolveSession(req *http.Request) (Session, error) {

	cookie := findSessionCookie(req)

	if cookie == nil {
		return newSession(), nil
	}

	sessionCand := Session{}
	exists, err := CurrentApplication.SessionStore.Get(cookie.Value, &sessionCand)

	//TODO check session timeout

	if err != nil {
		return newSession(), err
	}

	if exists {
		return sessionCand, nil
	}

	return newSession(), nil

}

func newSession() Session {
	return Session{}
}
