package applications

import "net/http"

//Session models an interactive sesion
type Session struct {
	SessionKey           string
	UserID               string
	FirstName            string
	LastName             string
	TenantID             string
	EffectivePermissions map[string]bool //because maps are faster lookups
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
