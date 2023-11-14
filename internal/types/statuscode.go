package types

const (
	StatusOk         StatusCode = "OK"          // StatusOk is a status for successful request
	StatusBadRequest StatusCode = "Bad request" // StatusBadRequest is a status for bad request
	StatusForbidden  StatusCode = "Forbidden"   // StatusForbidden is a status for forbidden request
)

// Status is a type for status of the request
type StatusCode string

func (StatusCode) sealed() {}
