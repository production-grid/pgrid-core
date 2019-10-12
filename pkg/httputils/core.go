package httputils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/production-grid/pgrid-core/pkg/applications"
	"github.com/production-grid/pgrid-core/pkg/logging"
)

// Send403 sends a 403 Access Denied Error to the client.
func Send403(w http.ResponseWriter) {
	SendError(NewError(403, "Access Denied"), w)
}

// Send404 sends a 403 Not Found Error to the client.
func Send404(w http.ResponseWriter) {
	SendError(NewError(404, "Not Found"), w)
}

// FindSessionCookie locates the production grid session cookie in the request and
// returns it.
func FindSessionCookie(req *http.Request) *http.Cookie {

	for _, cookie := range req.Cookies() {
		if cookie.Name == applications.SessionCookieName {
			return cookie
		}
	}

	return nil

}

/*
ConsumeRequestBody assumes the request body is JSON, parses it, and populates the
given result interface.
*/
func ConsumeRequestBody(result interface{}, w http.ResponseWriter, req *http.Request) bool {

	content, err := ioutil.ReadAll(req.Body)

	if err != nil {
		SendError(err, w)
		return false
	}

	err = json.Unmarshal(content, &result)

	if err != nil {
		SendError(err, w)
		return false
	}

	return true

}

/*
SendError sends a generic error to the response writer, including a 500 status
code if the given error does not implement api.Error.
*/
func SendError(err error, w http.ResponseWriter) {

	logging.Errorf("Sending backend error: %v", err)

	apiErr, ok := err.(Error)

	ack := Acknowledgement{}
	ack.Success = false
	ack.Error = err.Error()
	ack.Description = err.Error()

	if ok {
		(w).WriteHeader(apiErr.StatusCode)
	} else {
		(w).WriteHeader(500)
	}

	SendJSON(ack, w)

}

/*
SendJSON sends the given responsePayload to the response writer as JSON.
*/
func SendJSON(responsePayload interface{}, w http.ResponseWriter) {

	content, err := json.Marshal(responsePayload)

	if err != nil {
		logging.Errorf("Failed to JSON encode payload: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Lenth", strconv.Itoa(len(content)))
	w.Write(content)

}
