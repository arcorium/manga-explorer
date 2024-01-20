package users

import (
	"errors"
)

var ErrHashPassword = errors.New("failed to hash password")
var ErrEmailValidation = errors.New("email is invalid")
var ErrUnknownRole = errors.New("role unknown")
var ErrUnknownVerificationUsage = errors.New("usage unknown")
