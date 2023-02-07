package go_hidrive

import (
	"encoding/json"
	"errors"
	"fmt"
)

type HiDriveError struct {
	Code    json.Number `json:"code"`
	Message string      `json:"msg"`
}

// HiDriveError returns a string for the error and satisfies the error interface.
func (e *HiDriveError) Error() string {
	out := fmt.Sprintf("HiDriveError %q", e.Code.String())
	if e.Message != "" {
		out += ": " + e.Message
	}
	return out
}

var (
	ErrAuthNoToken = errors.New("no refresh token available")
	//ErrNotImplemented = errors.New("function is missing implementation")
)
