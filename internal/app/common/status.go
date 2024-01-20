package common

import (
	"database/sql"
	"errors"
	"github.com/uptrace/bun/driver/pgdriver"
	"manga-explorer/internal/app/common/status"
)

// ConditionalStatus Used to create common.Status based on error
func ConditionalStatus(err error, badCode int) Status {
	if err != nil {
		return StatusError(badCode)
	}
	return StatusSuccess()
}

// StatusError Used to create common.Status
func StatusError(code int) Status {
	message, ok := status.Messages[code]
	if !ok {
		panic("Status message is not defined!")
	}

	return Status{
		Code: uint(code),
		err:  errors.New(message),
	}
}

func NewRepositoryStatus(err error, successCode ...uint) Status {
	if err != nil {
		if err == sql.ErrNoRows {
			return StatusError(status.OBJECT_NOT_FOUND)
		}
		var pgerror pgdriver.Error
		ok := errors.As(err, &pgerror)
		if !ok {
			return StatusError(status.INTERNAL_SERVER_ERROR)
		}
		if pgerror.IntegrityViolation() {
			return StatusError(status.BAD_BODY_REQUEST_ERROR)
		}
		str := pgerror.Field('C')
		switch str {
		case "22P02", "22P03":
			return StatusError(status.BAD_BODY_REQUEST_ERROR)
		}
		return StatusError(status.INTERNAL_SERVER_ERROR)
	}
	return StatusSuccess(successCode...)
}

// StatusSuccess Use it when there is no error, it works like returning nil on error
func StatusSuccess(code ...uint) Status {
	var _code uint = status.SUCCESS
	if len(code) == 1 {
		_code = code[0]
	}
	return Status{
		Code: _code,
		err:  nil,
	}
}

// Status Wrapped error with the error code
type Status struct {
	Code uint
	err  error
}

// IsError Check error existences, it works like err != nil
func (e Status) IsError() bool {
	return e.err != nil
}

// ErrorMessage Returning the error message, it works like err.Status()
func (e Status) ErrorMessage() string {
	return e.err.Error()
}
