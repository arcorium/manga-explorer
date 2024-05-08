package status

import (
	"database/sql"
	"errors"
	"github.com/uptrace/bun/driver/pgdriver"
	"manga-explorer/internal/util/opt"
)

func Conditional(code Code, errMessage ...string) Object {
	if code.IsError() {
		return Error(code, errMessage...)
	}
	return Success(code)
}

func Success(code ...Code) Object {
	var _code = SUCCESS
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
func Error(code Code, msg ...string) Object {
	message, ok := messages[code]
	if !ok {
		panic("Object message is not defined!")
	}

	return New(code, errors.New(message), msg...)
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

func RepositoryError(err error, notFoundStatus opt.Optional[Code]) Object {
	return RepositoryErrorE(err, notFoundStatus, opt.Null[Code]())
}

func RepositoryErrorE(err error, notFoundStatus opt.Optional[Code], violationStatus opt.Optional[Code]) Object {
	if errors.Is(err, sql.ErrNoRows) {
		return Conditional(notFoundStatus.ValueOr(OBJECT_NOT_FOUND))
	}
	var pgerror pgdriver.Error
	ok := errors.As(err, &pgerror)
	if !ok {
		return New(INTERNAL_SERVER_ERROR, err)
	}
	if pgerror.IntegrityViolation() {

		return Conditional(violationStatus.ValueOr(BAD_REQUEST_ERROR))
	}
	str := pgerror.Field('C')
	switch str {
	case "22P02", "22P03":
		return Conditional(violationStatus.ValueOr(BAD_REQUEST_ERROR))
	}
	return InternalError()
}

func ConditionalRepository(err error, successCode Code, notFoundStatus opt.Optional[Code]) Object {
	return ConditionalRepositoryE(err, successCode, notFoundStatus, opt.Null[Code]())
}

func ConditionalRepositoryE(err error, successCode Code, notFoundStatus opt.Optional[Code], violationStatus opt.Optional[Code]) Object {
	if err != nil {
		return RepositoryErrorE(err, notFoundStatus, violationStatus)
	}
	return Success(successCode)
}
