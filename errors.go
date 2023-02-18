package go_hidrive

import (
	"encoding/json"
	"errors"
	"fmt"
)

/*
Error - represents HiDrive JSON-encoded errors.

Every time an API call receives non-OK code HiDrive also provides explanation in the response body.
This response is converted into this type and returned as error on each method.
*/
type Error struct {
	Code    json.Number `json:"code"`
	Message string      `json:"msg"`
}

// Error returns a string for the error and satisfies the error interface.
func (e *Error) Error() string {
	out := fmt.Sprintf("Error %q", e.Code.String())
	if e.Message != "" {
		out += ": " + e.Message
	}
	return out
}

var (
	ErrShouldNotBeEmpty = errors.New("value should not be empty")
)
