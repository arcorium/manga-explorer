package status

// StatusSuccess Use it when there is no error, it works like returning nil on error

func NewStatus(code uint, err error, details ...string) Object {
	var message string
	if len(details) == 1 {
		message = details[0]
	}
	return Object{
		Code:    code,
		err:     err,
		message: message,
		//details: details,
	}
}

// Object Wrapped error with the error code
type Object struct {
	Code    uint
	err     error
	message string
}

// IsError Check error existences, it works like err != nil
func (e Object) IsError() bool {
	return e.err != nil
}

// ErrorMessage Returning the error message, it works like err.Object()
func (e Object) ErrorMessage() string {
	return e.err.Error()
}

func (e Object) DetailMessage() string {
	return e.message
}

//func (e Object) Details() []any {
//	return e.details
//}
//
//func (e Object) AppendDetail(values ...any) {
//	for _, val := range values {
//		e.details = append(e.details, val)
//	}
//}
