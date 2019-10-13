package security

//Enumerates event keys
const (
	EventLogin       = "login"
	EventFailedLogin = "failed.login"
	EventLogout      = "logout"
)

//NewLoginEventMetaData is the event meta data factory function for login related events.
func NewLoginEventMetaData() interface{} {
	return LoginEventMetaData{}
}

//LoginEventMetaData models event meta data for login related events
type LoginEventMetaData struct {
	EMail string
}
