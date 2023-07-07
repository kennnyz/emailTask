package models

import "errors"

var (
	UserAlreadyExistErr = errors.New("user already exists")
	NotFoundUserErr     = errors.New("user not found")
)
