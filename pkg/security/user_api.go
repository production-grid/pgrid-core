package security

import (
	"net/http"

	"github.com/production-grid/pgrid-core/pkg/applications"
	"github.com/production-grid/pgrid-core/pkg/crypto"
	"github.com/production-grid/pgrid-core/pkg/database/relational"
	"github.com/production-grid/pgrid-core/pkg/events"
	"github.com/production-grid/pgrid-core/pkg/httputils"
	"github.com/production-grid/pgrid-core/pkg/logging"
)

//LoginRequest models an API login request.
type LoginRequest struct {
	XSRF         string `json:"xsrf"`
	EmailAddress string `json:"email"`
	Password     string `json:"password"`
}

//SessionDTO models an interactive sesion
type SessionDTO struct {
	UserID               string   `json:"userId"`
	FirstName            string   `json:"firstName"`
	LastName             string   `json:"lastName"`
	TenantID             string   `json:"tenantId,omitempty"`
	TenantName           string   `json:"tenantName,omitempty"`
	AdminLogo            string   `json:"adminLogo,omitempty"`
	ApplicationName      string   `json:"applicationName"`
	TagLine              string   `json:"tagline"`
	EffectivePermissions []string `json:"effectivePermissions"`
}

//GetSession returns session meta data
func GetSession(session applications.Session, w http.ResponseWriter, req *http.Request) {

	dto := SessionDTO{
		UserID:          session.UserID,
		FirstName:       session.FirstName,
		LastName:        session.LastName,
		TenantID:        session.TenantID,
		ApplicationName: applications.CurrentApplication.Name,
		TagLine:         applications.CurrentApplication.TagLine,
	}

	if session.EffectivePermissions != nil {
		perms := make([]string, 0)

		for key, val := range session.EffectivePermissions {
			if val {
				perms = append(perms, key)
			}
		}

		dto.EffectivePermissions = perms
	}

	httputils.SendJSON(dto, w)

}

//PostLogout terminates an interactive session
func PostLogout(session applications.Session, w http.ResponseWriter, req *http.Request) {

	if session.SessionKey != "" {
		applications.CurrentApplication.SessionStore.Delete(session.SessionKey)
	}

	httputils.SendJSON(httputils.Acknowledgement{Success: true}, w)
}

//PostLogin processes an interactive login request
func PostLogin(session applications.Session, w http.ResponseWriter, req *http.Request) {

	request := LoginRequest{}

	if httputils.ConsumeRequestBody(&request, w, req) {

		logging.LogJSON(request)

		userFinder := UserFinder{}
		user, err := userFinder.FindByEmailAddress(relational.REPLICA, request.EmailAddress)
		if err != nil {
			httputils.SendError(err, w)
			return
		}

		if user == nil {
			crypto.SimulateDoubleHashTimeDelay()
			httputils.SendJSON(httputils.Acknowledgement{Success: false, Error: "Invalid Login"}, w)
			return
		}
		if user.IsLocked {
			httputils.SendJSON(httputils.Acknowledgement{Success: false, Error: "User account locked"}, w)
			return
		}
		expectedHash := crypto.DoubleHash{
			InnerSalt: user.InnerSalt,
			OuterHash: user.PasswordHash,
		}
		if !crypto.CompareDoubleHash(request.Password, expectedHash) {
			httputils.SendJSON(httputils.Acknowledgement{Success: false, Error: "Invalid Login"}, w)
			return
		}

		secureSession, err := user.InitSession(&session)
		if err != nil {
			httputils.SendError(err, w)
			return
		}
		cookie := http.Cookie{Name: applications.SessionCookieName, Value: secureSession.SessionKey, Path: "/"}
		if applications.CurrentApplication.CoreConfiguration.SecureCookies {
			cookie.Secure = true
		}
		http.SetCookie(w, &cookie)
		httputils.SendJSON(httputils.Acknowledgement{Success: true, ID: secureSession.SessionKey}, w)
		events.Dispatch(EventLogin, *secureSession, LoginEventMetaData{EMail: request.EmailAddress})
	}

}
