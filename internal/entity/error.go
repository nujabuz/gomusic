package entity

import "errors"

var (
	// any errors
	ErrUnexpected       = errors.New("an unexpected error occurred on the server")
	ErrParseRequestBody = errors.New("an error occurred while parsing the request body")
	ErrNotFound         = errors.New("not found")

	// user errors
	ErrUserNotFound      = errors.New("user not found.")
	ErrUserAlreadyExists = errors.New("user already exists.")
)
