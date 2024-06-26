package resp

import (
  "manga-explorer/internal/common/status"
  "net/http"
)

func HttpCodeFromError(err status.Object) int {
  switch status.Code(err.Code) {
  case status.SUCCESS, status.UPDATED, status.DELETED:
    return http.StatusOK
  case status.CREATED:
    return http.StatusCreated

  case status.ERROR_WITH_MESSAGE, status.BAD_REQUEST_ERROR, status.BAD_QUERY_ERROR, status.BAD_PARAMETER_ERROR,
    status.OBJECT_NOT_FOUND, status.VERIFICATION_REQUEST_FAILED, status.VERIFICATION_USER_NOT_EXISTS,
    status.VERIFICATION_TOKEN_NOT_FOUND, status.VERIFICATION_TOKEN_EXPIRED, status.VERIFICATION_TOKEN_MISUSE,
    status.USER_NOT_FOUND, status.USER_LOGIN_ERROR, status.USER_UPDATE_FAILED, status.USER_CREATION_BAD_EMAIL,
    status.USER_CREATION_ALREADY_EXIST, status.USER_CHANGE_PASSWORD_WRONG_PASSWORD, status.PROFILE_NOT_FOUND,
    status.PROFILE_UPDATE_FAILED, status.TOKEN_LOOKUP_MALFORMED, status.TOKEN_MALTYPE, status.TOKEN_NOT_VALID,
    status.FILE_UPLOAD_FAILED, status.MAIL_SEND_FAILED, status.MANGA_CREATE_ALREADY_EXIST, status.MANGA_NOT_FOUND,
    status.MANGA_UPDATE_FAILED, status.MANGA_TRANSLATION_ALREADY_EXIST, status.MANGA_HAS_NO_TRANSLATIONS,
    status.MANGA_TRANSLATION_NOT_FOUND, status.MANGA_TRANSLATION_UPDATE_FAILED, status.VOLUME_ALREADY_EXISTS,
    status.VOLUME_DELETE_FAILED, status.CHAPTER_UPDATE_FAILED, status.CHAPTER_NOT_FOUND, status.CHAPTER_ALREADY_EXIST,
    status.PAGE_INSERT_FAILED, status.PAGE_NOT_FOUND, status.GENRE_ALREADY_EXIST, status.GENRE_NOT_FOUND,
    status.RATING_NOT_FOUND, status.COMMENT_PARENT_NOT_FOUND, status.COMMENT_PARENT_DIFFERENT_SCOPE,
    status.COMMENT_CREATE_FAILED, status.VOLUME_CREATE_FAILED, status.MANGA_TRANSLATION_CREATE_FAILED,
    status.EMPTY_BODY_REQUEST:
    return http.StatusBadRequest
  case status.USER_AGENT_UNKNOWN_ERROR, status.CREDENTIALS_NOT_FOUND, status.JWT_TOKEN_MALFORMED,
    status.ACCESS_TOKEN_EXPIRED, status.ACCESS_TOKEN_WITHOUT_REFRESH_TOKEN, status.AUTH_UNAUTHORIZED,
    status.TOKEN_MALFORMED, status.LOGOUT_CREDENTIAL_NOT_FOUND:

    return http.StatusUnauthorized
  case status.INTERNAL_SERVER_ERROR:
    return http.StatusInternalServerError
  }
  panic("Status code is not implemented for this message")
}
