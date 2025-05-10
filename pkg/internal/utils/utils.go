package utils

import (
	"regexp"
)

var OnlyNumberRegex = regexp.MustCompile(`[0-9]+`)

var AllowedTransportTypes = []string{"metro", "bus", "rail", "tram"}

// RequestError represents request-related errors that should return 400 Bad request
type RequestError struct {
	Message string
}

func (e *RequestError) Error() string {
	return e.Message
}
