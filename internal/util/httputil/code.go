package httputil

import (
	"manga-explorer/internal/app/common"
	"manga-explorer/internal/app/common/status"
	"net/http"
)

func HttpCodeFromError(err common.Status) int {
	switch err.Code {
	case status.SUCCESS:
		return http.StatusOK
	case status.SUCCESS_CREATED:
		return http.StatusCreated
	case status.VERIFICATION_TOKEN_MISUSE, status.BAD_BODY_REQUEST_ERROR, status.BAD_PARAMETER_ERROR,
		status.JWT_TOKEN_MALFORMED, status.TOKEN_MALFORMED, status.TOKEN_MALTYPE, status.TOKEN_LOOKUP_MALFORMED:
		return http.StatusBadRequest
	case status.VERIFICATION_TOKEN_NOT_FOUND, status.USER_AGENT_UNKNOWN_ERROR, status.ACCESS_TOKEN_EXPIRED,
		status.ACCESS_TOKEN_WITHOUT_REFRESH_TOKEN, status.AUTH_UNAUTHORIZED:
		return http.StatusUnauthorized
	}

	return http.StatusInternalServerError
}
