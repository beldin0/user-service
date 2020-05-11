package userservice

import "errors"

var (
	// ErrDuplicate is the error returned when an Add request is sent with an email that is already in use
	ErrDuplicate = errors.New("key already exists")
)
