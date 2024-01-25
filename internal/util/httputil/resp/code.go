package resp

import (
	"manga-explorer/internal/common/status"
	"net/http"
)

func HttpCodeFromError(err status.Object) int {
	switch err.Code {
	case status.SUCCESS:
		return http.StatusOK
	case status.CREATED:
		return http.StatusCreated
	case status.VERIFICATION_TOKEN_MISUSE, status.BAD_REQUEST_ERROR, status.BAD_PARAMETER_ERROR,
		status.JWT_TOKEN_MALFORMED, status.TOKEN_MALFORMED, status.TOKEN_MALTYPE, status.TOKEN_LOOKUP_MALFORMED:
		return http.StatusBadRequest
	case status.VERIFICATION_TOKEN_NOT_FOUND, status.USER_AGENT_UNKNOWN_ERROR, status.ACCESS_TOKEN_EXPIRED,
		status.ACCESS_TOKEN_WITHOUT_REFRESH_TOKEN, status.AUTH_UNAUTHORIZED:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
