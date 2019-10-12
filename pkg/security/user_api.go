package security

import (
	"net/http"

	"github.com/production-grid/pgrid-core/pkg/applications"
	"github.com/production-grid/pgrid-core/pkg/crypto"
	"github.com/production-grid/pgrid-core/pkg/database/relational"
	"github.com/production-grid/pgrid-core/pkg/httputils"
)

//LoginRequest models an API login request.
type LoginRequest struct {
	XSRF         string `json:"xsrf"`
	EmailAddress string `json:"email"`
	Password     string `json:"password"`
}

//PostLogin processes an interactive login request
func PostLogin(session applications.Session, w http.ResponseWriter, req *http.Request) {

	request := LoginRequest{}

	if httputils.ConsumeRequestBody(&request, w, req) {

		userFinder := UserFinder{}
		user, err := userFinder.FindByEmailAddress(relational.REPLICA, request.EmailAddress)
		if err != nil {
			httputils.SendError(err, w)
			return
		}

		if user == nil {
			crypto.SimulateDoubleHashTimeDelay()
			httputils.SendJSON(httputils.Acknowledgement{Success: false, Description: "Invalid Login"}, w)
			return
		}
		if user.IsLocked {
			httputils.SendJSON(httputils.Acknowledgement{Success: false, Description: "User account locked"}, w)
			return
		}
		expectedHash := crypto.DoubleHash{
			InnerSalt: user.InnerSalt,
			OuterHash: user.PasswordHash,
		}
		if !crypto.CompareDoubleHash(request.Password, expectedHash) {
			httputils.SendJSON(httputils.Acknowledgement{Success: false, Description: "Invalid Login"}, w)
			return
		}

		session, err := user.InitSession()
		if err != nil {
			httputils.SendError(err, w)
			return
		}
		cookie := http.Cookie{Name: applications.SessionCookieName, Value: session.SessionKey, Path: "/"}
		if !applications.CurrentApplication.CoreConfiguration.SecureCookies {
			cookie.Secure = true
		}
		http.SetCookie(w, &cookie)
		httputils.SendJSON(httputils.Acknowledgement{Success: true}, w)
	}

}

//PostUserReg processes an incoming user registration request
func PostUserReg(session applications.Session, w http.ResponseWriter, req *http.Request) {

}
