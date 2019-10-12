package httputils

import "net/http"

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
