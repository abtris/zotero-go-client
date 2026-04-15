package zotero

import (
	"fmt"
	"net/http"
)

// APIError represents an error response from the Zotero API.
type APIError struct {
	StatusCode int
	Status     string
	Message    string
	RetryAfter int // seconds; set on 429 responses
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("zotero: %s: %s", e.Status, e.Message)
	}
	return fmt.Sprintf("zotero: %s", e.Status)
}

// WriteResponse represents the result of a multi-object write request.
type WriteResponse struct {
	Success    map[string]string       `json:"success"`
	Unchanged  map[string]string       `json:"unchanged"`
	Failed     map[string]WriteFailure `json:"failed"`
	Successful map[string]*Item        `json:"successful"`
}

// WriteFailure represents a single failed write within a multi-object request.
type WriteFailure struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// IsNotFound reports whether the error is a 404 Not Found response.
func IsNotFound(err error) bool {
	return isStatus(err, http.StatusNotFound)
}

// IsConflict reports whether the error is a 409 Conflict response.
func IsConflict(err error) bool {
	return isStatus(err, http.StatusConflict)
}

// IsRateLimited reports whether the error is a 429 Too Many Requests response.
func IsRateLimited(err error) bool {
	return isStatus(err, http.StatusTooManyRequests)
}

// IsPreconditionFailed reports whether the error is a 412 Precondition Failed response.
func IsPreconditionFailed(err error) bool {
	return isStatus(err, http.StatusPreconditionFailed)
}

func isStatus(err error, code int) bool {
	if err == nil {
		return false
	}
	apiErr, ok := err.(*APIError)
	return ok && apiErr.StatusCode == code
}
