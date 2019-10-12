package httputils

// Acknowledgement models the response to API requests that don't really return any
// detailed information, just whether or not the operation succeeded.
type Acknowledgement struct {
	Success        bool   `json:"success"`
	ID             string `json:"id,omitempty"`
	Description    string `json:"responseDescription,omitempty"`
	Error          string `json:"error,omitempty"`
	LocalizedError string `json:"localized_error,omitempty"`
}
