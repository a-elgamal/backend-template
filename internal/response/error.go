// Package response contains utilities that aid returning API responses from handlers
package response

import "fmt"

// ErrorResponse models Errors returned from the APIs
type ErrorResponse struct {
	Err ErrorDetail `json:"error"`
}

// ErrorDetail contains the details of the error that occur
type ErrorDetail struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
}

// Error returns a description of the error response
func (e ErrorResponse) Error() string {
	return fmt.Sprintf("%v: %v", e.Err.Code, e.Err.Msg)
}
