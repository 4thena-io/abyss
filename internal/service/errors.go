package service

import "errors"

var ErrNotFound = errors.New("not found")

var ErrConflict = errors.New("conflict")

var ErrUnauthorized = errors.New("unauthorized")

var ErrForbidden = errors.New("forbidden")
