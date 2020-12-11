package api_errors

import "errors"

//Common Api Errors
var (
	ErrInvalidData = errors.New("invalid_data")

	ErrAlreadyExists = errors.New("already_exist")

	ErrNotFound = errors.New("record_not_found")

	ErrInvalidFilename = errors.New("invalid_filename")

	ErrInternalError = errors.New("internal_error")

	ErrDoesNotOwn = errors.New("does_not_own")

	ErrCouldNotDelete = errors.New("could_not_delete")

	ErrCouldNotCreate = errors.New("could_not_create")

	ErrCouldNotUpdate = errors.New("could_not_update")

	ErrAuth = errors.New("invalid_credentials")

	ErrSessionExpired = errors.New("session_expired")
)
