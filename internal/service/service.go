package service

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrAdminExists        = errors.New("admin already exists")
	ErrAppNotFound        = errors.New("app not found")
	ErrUserNotFound       = errors.New("user not found")
	ErrAdminNotFound      = errors.New("admin not found")
)
