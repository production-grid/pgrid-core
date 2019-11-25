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
	XSRF     string `json:"xsrf"`
	LoginID  string `json:"loginId"`
	Password string `json:"password"`
}

//SessionDTO models an interactive sesion
type SessionDTO struct {
	UserID               string          `json:"userId"`
	FirstName            string          `json:"firstName"`
	LastName             string          `json:"lastName"`
	TenantID             string          `json:"tenantId,omitempty"`
	TenantName           string          `json:"tenantName,omitempty"`
	AdminLogo            string          `json:"adminLogo,omitempty"`
	ApplicationName      string          `json:"applicationName"`
	TagLine              string          `json:"tagline"`
	TenantPlural         string          `json:"tenantPlural"`
	TenantSingular       string          `json:"tenantSingular"`
	RootTenantHost       string          `json:"rootTenantHost"`
	EffectivePermissions []string        `json:"effectivePermissions"`
	TenantTypes          []TenantTypeDTO `json:"tenantTypes"`
}

//TenantTypeDTO models human readable information about tenant types
type TenantTypeDTO struct {
	Singular string `json:"singular"`
	Plural   string `json:"plural"`
}

//GetSession returns session meta data
func GetSession(session applications.Session, w http.ResponseWriter, req *http.Request) {

	lingo := applications.CurrentApplication.TenantLingo

	dto := SessionDTO{
		UserID:          session.UserID,
		FirstName:       session.FirstName,
		LastName:        session.LastName,
		TenantID:        session.TenantID,
		ApplicationName: applications.CurrentApplication.Name,
		TagLine:         applications.CurrentApplication.TagLine,
		TenantPlural:    lingo.TenantPlural,
		TenantSingular:  lingo.TenantSingular,
		RootTenantHost:  applications.CurrentApplication.CoreConfiguration.RootTenantHost,
	}

	tenantTypes := make([]TenantTypeDTO, len(lingo.Types))
	for idx, tenantType := range lingo.Types {
		tenantTypes[idx] = TenantTypeDTO{
			Singular: tenantType.Singular,
			Plural:   tenantType.Plural,
		}
	}
	dto.TenantTypes = tenantTypes

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
		user, err := userFinder.FindByLoginID(relational.REPLICA, request.LoginID)
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
			httputils.SendJSON(httputils.Acknowledgement{Success: false, Error: "User Account Locked"}, w)
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
		events.Dispatch(EventLogin, *secureSession, LoginEventMetaData{EMail: user.EMail})
	}

}
