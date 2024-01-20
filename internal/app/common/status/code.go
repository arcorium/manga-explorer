package status

const (
	// Success Code
	SUCCESS = iota
	CREATED
	UPDATED
	DELETED

	// Object Code
	BAD_BODY_REQUEST_ERROR
	USER_AGENT_UNKNOWN_ERROR
	BAD_PARAMETER_ERROR
	BAD_QUERY_ERROR
	INTERNAL_SERVER_ERROR
	CREDENTIALS_NOT_FOUND
	// General
	OBJECT_NOT_FOUND

	VERIFICATION_TOKEN_NOT_FOUND
	VERIFICATION_TOKEN_EXPIRED
	VERIFICATION_TOKEN_MISUSE

	// User
	USER_LOGIN_ERROR
	USER_NOT_FOUND
	USER_UPDATE_FAILED
	USER_CREATION_BAD_EMAIL
	USER_CREATION_ALREADY_EXIST

	JWT_TOKEN_MALFORMED
	ACCESS_TOKEN_EXPIRED
	ACCESS_TOKEN_WITHOUT_REFRESH_TOKEN
	AUTH_UNAUTHORIZED
	TOKEN_MALFORMED

	// Auth Middleware
	TOKEN_LOOKUP_MALFORMED
	TOKEN_MALTYPE
	TOKEN_NOT_VALID

	// File Upload
	FILE_UPLOAD_FAILED
)

var Messages = map[int]string{
	BAD_BODY_REQUEST_ERROR:       "Request body is malformed",
	USER_AGENT_UNKNOWN_ERROR:     "Unknown user-agent",
	BAD_QUERY_ERROR:              "URL query is have malformed type",
	BAD_PARAMETER_ERROR:          "URL parameter is missing",
	INTERNAL_SERVER_ERROR:        "Internal server error",
	AUTH_UNAUTHORIZED:            "You are not authorized to access this",
	VERIFICATION_TOKEN_NOT_FOUND: "Token you provide is wrong",
	VERIFICATION_TOKEN_EXPIRED:   "Token you provide is expired",
	VERIFICATION_TOKEN_MISUSE:    "Token you provide should not used here",
	OBJECT_NOT_FOUND:             "Object not found",

	// Token
	TOKEN_MALFORMED:                    "Token you provide is malformed",
	ACCESS_TOKEN_EXPIRED:               "Token you provide is expired",
	ACCESS_TOKEN_WITHOUT_REFRESH_TOKEN: "Access Token has no associated refresh token",
	TOKEN_MALTYPE:                      "Token has invalid type",

	// User
	USER_LOGIN_ERROR:            "Email/Password is not match!",
	USER_NOT_FOUND:              "User doesn't exists",
	USER_UPDATE_FAILED:          "Could not update user",
	USER_CREATION_BAD_EMAIL:     "Email is not valid. Please provide valid email",
	USER_CREATION_ALREADY_EXIST: "User with those email is already exist",

	// Auth Middleware
	TOKEN_LOOKUP_MALFORMED: "Token you provide is malformed",
	TOKEN_NOT_VALID:        "Token you provide is not valid",
	CREDENTIALS_NOT_FOUND:  "Credentials is empty",

	// File Upload
	FILE_UPLOAD_FAILED: "Failed to upload an image",
}
