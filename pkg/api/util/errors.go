package util

import "errors"

var (
	ErrInvalidUUID      = errors.New("invalid UUID")
	ErrForbidden        = errors.New("user is not authorized to perform this action")
	ErrMalformedRequest = errors.New("malformed request")
	ErrAlreadyExists    = errors.New("item already exists")
	ErrNoChanges        = errors.New("no changes made")
	ErrNotFound         = errors.New("item not found")
)
