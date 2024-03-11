package storage

import "errors"

var (
	ErrUserExists    = errors.New("user already exists")
	ErrAdminExists   = errors.New("admin already exists")
	ErrAdminNotFound = errors.New("admin not found")
	ErrUserNotFound  = errors.New("user not found")
	ErrAppNotFound   = errors.New("app not found")
	ErrAppExists     = errors.New("app already exists")
)
