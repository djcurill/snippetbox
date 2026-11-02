package models

import "errors"

var (
	ErrNoRecord           error = errors.New("models: no matching record")
	ErrInavlidCredentials error = errors.New("models: wrong username or password")
	ErrDpuplicateEmail    error = errors.New("models: email already exists")
)
