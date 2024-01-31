package status

type Code int

func (c Code) Underlying() int {
	return int(c)
}

func (c Code) IsError() bool {
	return c.Underlying() > 4
}

const (
	// Success Code
	SUCCESS Code = iota
	CREATED
	UPDATED
	DELETED
	INTERNAL

	ERROR_WITH_MESSAGE

	// Object Code
	BAD_REQUEST_ERROR
	USER_AGENT_UNKNOWN_ERROR
	BAD_PARAMETER_ERROR
	EMPTY_BODY_REQUEST
	BAD_QUERY_ERROR
	INTERNAL_SERVER_ERROR
	CREDENTIALS_NOT_FOUND
	// General
	OBJECT_NOT_FOUND

	// Verification
	VERIFICATION_REQUEST_FAILED
	VERIFICATION_USER_NOT_EXISTS
	VERIFICATION_TOKEN_NOT_FOUND
	VERIFICATION_TOKEN_EXPIRED
	VERIFICATION_TOKEN_MISUSE

	// User
	USER_LOGIN_ERROR
	USER_NOT_FOUND
	USER_UPDATE_FAILED
	USER_CREATION_BAD_EMAIL
	USER_CREATION_ALREADY_EXIST
	USER_CHANGE_PASSWORD_WRONG_PASSWORD

	// Profile
	PROFILE_NOT_FOUND
	PROFILE_UPDATE_FAILED

	JWT_TOKEN_MALFORMED
	ACCESS_TOKEN_EXPIRED
	ACCESS_TOKEN_WITHOUT_REFRESH_TOKEN
	AUTH_UNAUTHORIZED
	TOKEN_MALFORMED

	// Auth Middleware
	TOKEN_LOOKUP_MALFORMED
	TOKEN_MALTYPE
	TOKEN_NOT_VALID
	LOGOUT_CREDENTIAL_NOT_FOUND

	// File
	FILE_UPLOAD_FAILED

	// Email
	MAIL_SEND_FAILED

	// Manga
	MANGA_CREATE_ALREADY_EXIST
	MANGA_NOT_FOUND
	MANGA_UPDATE_FAILED

	// Translation
	MANGA_TRANSLATION_ALREADY_EXIST
	MANGA_HAS_NO_TRANSLATIONS
	MANGA_TRANSLATION_NOT_FOUND
	MANGA_TRANSLATION_UPDATE_FAILED
	MANGA_TRANSLATION_CREATE_FAILED

	VOLUME_ALREADY_EXISTS
	VOLUME_CREATE_FAILED
	VOLUME_DELETE_FAILED

	CHAPTER_UPDATE_FAILED
	CHAPTER_NOT_FOUND
	CHAPTER_ALREADY_EXIST

	PAGE_INSERT_FAILED
	PAGE_NOT_FOUND

	GENRE_ALREADY_EXIST
	GENRE_NOT_FOUND

	// Rating
	RATING_NOT_FOUND

	// Comment
	COMMENT_PARENT_NOT_FOUND
	COMMENT_PARENT_DIFFERENT_SCOPE
	COMMENT_CREATE_FAILED
)

var messages = map[Code]string{
	SUCCESS: "Success",
	CREATED: "Created",
	UPDATED: "Updated",
	DELETED: "Deleted",

	ERROR_WITH_MESSAGE:       "Error",
	BAD_REQUEST_ERROR:        "Request is malformed or invalid",
	EMPTY_BODY_REQUEST:       "Expected request body",
	USER_AGENT_UNKNOWN_ERROR: "Unknown user-agent",
	BAD_QUERY_ERROR:          "URL query is have malformed type",
	BAD_PARAMETER_ERROR:      "URL parameter is missing",
	INTERNAL_SERVER_ERROR:    "Internal server error",
	AUTH_UNAUTHORIZED:        "You are not authorized to access this",

	// Verification
	VERIFICATION_REQUEST_FAILED:  "Failed to request verification",
	VERIFICATION_USER_NOT_EXISTS: "Trying to request verification to user that doesn't exists",
	VERIFICATION_TOKEN_NOT_FOUND: "Verification token you provide is wrong",
	VERIFICATION_TOKEN_EXPIRED:   "Verification token you provide is expired",
	VERIFICATION_TOKEN_MISUSE:    "Verification token you provide should not used here",
	OBJECT_NOT_FOUND:             "Object not found",

	// Token
	TOKEN_MALFORMED:                    "Token you provide is malformed",
	ACCESS_TOKEN_EXPIRED:               "Token you provide is expired",
	ACCESS_TOKEN_WITHOUT_REFRESH_TOKEN: "Access Token has no associated refresh token",
	TOKEN_MALTYPE:                      "Token has invalid type",

	// User
	USER_LOGIN_ERROR:                    "Email/Password is not match!",
	USER_NOT_FOUND:                      "User doesn't exists",
	USER_UPDATE_FAILED:                  "Could not update user",
	USER_CREATION_BAD_EMAIL:             "Email is not valid. Please provide valid email",
	USER_CREATION_ALREADY_EXIST:         "Email is already used",
	USER_CHANGE_PASSWORD_WRONG_PASSWORD: "Password you provided doesn't match",

	// Profile
	PROFILE_NOT_FOUND:     "User profile doesn't exists",
	PROFILE_UPDATE_FAILED: "Failed to update user profile",

	// Auth Middleware
	TOKEN_LOOKUP_MALFORMED:      "Token you provide is malformed",
	TOKEN_NOT_VALID:             "Token you provide is not valid",
	CREDENTIALS_NOT_FOUND:       "You have no devices currently logged in",
	LOGOUT_CREDENTIAL_NOT_FOUND: "Logout error due to no credential match found",

	// File
	FILE_UPLOAD_FAILED: "Failed to upload an image",

	// Mail
	MAIL_SEND_FAILED: "Email could not be sent",

	// Manga
	MANGA_CREATE_ALREADY_EXIST: "Manga you try to create is already exist",
	MANGA_NOT_FOUND:            "Manga not found",
	MANGA_UPDATE_FAILED:        "Could not update those manga",

	MANGA_TRANSLATION_ALREADY_EXIST: "Manga translation is already exist",
	MANGA_HAS_NO_TRANSLATIONS:       "Those manga has no translations",
	MANGA_TRANSLATION_NOT_FOUND:     "Manga translation not found",
	MANGA_TRANSLATION_UPDATE_FAILED: "Could not update manga translation",
	MANGA_TRANSLATION_CREATE_FAILED: "Could not create manga translation",

	VOLUME_ALREADY_EXISTS: "Volume you try to create is already exist",
	VOLUME_CREATE_FAILED:  "Failed to create volume",
	VOLUME_DELETE_FAILED:  "Volume is not found",

	CHAPTER_UPDATE_FAILED: "Could not update manga chapter",
	CHAPTER_NOT_FOUND:     "Manga chapter is not found",
	CHAPTER_ALREADY_EXIST: "Manga chapter is already exist",

	PAGE_INSERT_FAILED: "Could not insert pages on manga chapter",
	PAGE_NOT_FOUND:     "Chapter pages not found",

	GENRE_ALREADY_EXIST: "Magna genre already exist",
	GENRE_NOT_FOUND:     "Manga genre doesn't exist",

	RATING_NOT_FOUND: "Manga rating doesn't exist",

	COMMENT_PARENT_NOT_FOUND:       "Parent comment is not found",
	COMMENT_PARENT_DIFFERENT_SCOPE: "You are trying to reply comment from different scope",
	COMMENT_CREATE_FAILED:          "Failed to create comment",
}
