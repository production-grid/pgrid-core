package httputils

import "net/http"

// ErrTestModeDisabled is thrown when test credentials are used on
// the live API.
var ErrTestModeDisabled = NewError(http.StatusBadRequest, "test credentials cannot be used on the live API: use the test flag for testing")

// ErrTestModeEnabled is thrown when live credentials are used on
// the test API.
var ErrTestModeEnabled = NewError(http.StatusBadRequest, "live credentials cannot be used on the test API: disable the test flag for live transactions")

// ErrBlockchainPurgeDisabled is thrown purges are attempted on the production blockchain
var ErrBlockchainPurgeDisabled = NewError(http.StatusBadRequest, "blockchain purge not enabled")

/*
Error models a special api error type.
*/
type Error struct {
	StatusCode  int
	Description string
}

func (err Error) Error() string {
	return err.Description
}

/*
NewError returns a new api error.
*/
func NewError(status int, desc string) Error {
	return Error{StatusCode: status, Description: desc}
}

// NewStatusError returns a generic HTTP error for a status.
func NewStatusError(status int) Error {
	return Error{StatusCode: status, Description: http.StatusText(status)}
}
