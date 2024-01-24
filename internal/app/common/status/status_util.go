package status

import (
	"database/sql"
	"errors"
	"github.com/uptrace/bun/driver/pgdriver"
)

func Success(code ...uint) Object {
	var _code uint = SUCCESS
	if len(code) == 1 {
		if code[0] >= SUCCESS && code[0] <= DELETED {
			_code = code[0]
		} else {
			panic("Bad success code")
		}
	}

	return New(_code, nil)
}

func Updated() Object {
	return Success(UPDATED)
}

func Deleted() Object {
	return Success(DELETED)
}

func Created() Object {
	return Success(CREATED)
}

// Error Used to create common.Object
func Error(code int, msg ...string) Object {
	message, ok := Messages[code]
	if !ok {
		panic("Object message is not defined!")
	}

	return New(uint(code), errors.New(message), msg...)
}

func InternalError() Object {
	return Error(INTERNAL_SERVER_ERROR)
}

func NotFoundError() Object {
	return Error(OBJECT_NOT_FOUND)
}

func RepositoryError(err error) Object {

	if err == sql.ErrNoRows {
		return Error(OBJECT_NOT_FOUND)
	}
	var pgerror pgdriver.Error
	ok := errors.As(err, &pgerror)
	if !ok {
		return Error(INTERNAL_SERVER_ERROR)
	}
	if pgerror.IntegrityViolation() {
		return Error(BAD_BODY_REQUEST_ERROR)
	}
	str := pgerror.Field('C')
	switch str {
	case "22P02", "22P03":
		return Error(BAD_BODY_REQUEST_ERROR)
	}
	return InternalError()
}

func ConditionalRepository(err error, successCode uint) Object {
	if err != nil {
		return RepositoryError(err)
	}
	return Success(successCode)
}
