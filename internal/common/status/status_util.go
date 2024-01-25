package status

import (
	"database/sql"
	"errors"
	"github.com/uptrace/bun/driver/pgdriver"
	"manga-explorer/internal/util/opt"
)

func Success(code ...uint) Object {
	var _code uint = SUCCESS
	if len(code) == 1 {
		if code[0] >= SUCCESS && code[0] <= INTERNAL {
			_code = code[0]
		} else {
			panic("Bad success code")
		}
	}

	return New(_code, nil)
}

// InternalSuccess Used for communication and should not be used for response
func InternalSuccess() Object {
	return Success(INTERNAL)
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
	message, ok := messages[code]
	if !ok {
		panic("Object message is not defined!")
	}

	return New(uint(code), errors.New(message), msg...)
}

func ErrorMessage(message string) Object {
	return New(ERROR_WITH_MESSAGE, errors.New(message))
}

func InternalError() Object {
	return Error(INTERNAL_SERVER_ERROR)
}

func NotFoundError() Object {
	return Error(OBJECT_NOT_FOUND)
}

func BadRequestError() Object {
	return Error(BAD_REQUEST_ERROR)
}

func RepositoryError(err error, notFoundStatus opt.Optional[int]) Object {
	return RepositoryErrorE(err, notFoundStatus, opt.NullInt)
}

func RepositoryErrorE(err error, notFoundStatus opt.Optional[int], violationStatus opt.Optional[int]) Object {
	if errors.Is(err, sql.ErrNoRows) {
		return Error(notFoundStatus.ValueOr(OBJECT_NOT_FOUND))
	}
	var pgerror pgdriver.Error
	ok := errors.As(err, &pgerror)
	if !ok {
		return InternalError()
	}
	if pgerror.IntegrityViolation() {
		return Error(violationStatus.ValueOr(BAD_REQUEST_ERROR))
	}
	str := pgerror.Field('C')
	switch str {
	case "22P02", "22P03":
		return Error(violationStatus.ValueOr(BAD_REQUEST_ERROR))
	}
	return InternalError()
}

func ConditionalRepository(err error, successCode uint, notFoundStatus opt.Optional[int]) Object {
	return ConditionalRepositoryE(err, successCode, notFoundStatus, opt.NullInt)
}

func ConditionalRepositoryE(err error, successCode uint, notFoundStatus opt.Optional[int], violationStatus opt.Optional[int]) Object {
	if err != nil {
		return RepositoryErrorE(err, notFoundStatus, violationStatus)
	}
	return Success(successCode)
}
