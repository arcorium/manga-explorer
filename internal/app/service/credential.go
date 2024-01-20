package service

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"manga-explorer/internal/app/common"
	"manga-explorer/internal/app/common/status"
	"manga-explorer/internal/domain/auth"
	"manga-explorer/internal/domain/auth/dto"
	authMapper "manga-explorer/internal/domain/auth/mapper"
	authRepo "manga-explorer/internal/domain/auth/repository"
	authService "manga-explorer/internal/domain/auth/service"
	"manga-explorer/internal/domain/users"
	userRepo "manga-explorer/internal/domain/users/repository"
	"manga-explorer/internal/util"
	"manga-explorer/internal/util/containers"
)

func NewCredential(config *common.Config, credRepo authRepo.IAuthentication, userRepo userRepo.IUser) authService.IAuthentication {
	return &credentialService{config: config, authRepo: credRepo, userRepo: userRepo}
}

type credentialService struct {
	config *common.Config

	authRepo authRepo.IAuthentication
	userRepo userRepo.IUser
}

func (c credentialService) Authenticate(request *dto.LoginInput) (dto.LoginResponse, common.Status) {
	// Validate email
	if !request.Validate() {
		return dto.NoLoginResponse, common.StatusError(status.BAD_BODY_REQUEST_ERROR)
	}

	usr, err := c.userRepo.FindUserByEmail(request.Email)
	if err != nil {
		return dto.NoLoginResponse, common.StatusError(status.USER_NOT_FOUND)
	}

	if !usr.ValidatePassword(request.Password) {
		return dto.NoLoginResponse, common.StatusError(status.USER_LOGIN_ERROR)
	}

	// Refresh Token Creation
	refreshToken, err := util.GenerateJWTToken(users.DefaultClaims(c.config.RefreshTokenDuration, common.IssuerName), c.config.SigningMethod(), []byte(c.config.JWTSecretKey))
	if err != nil {
		return dto.NoLoginResponse, common.StatusError(status.INTERNAL_SERVER_ERROR)
	}

	// Access Token Creation
	accessTokenClaims := usr.GenerateAccessTokenClaims(c.config.AccessTokenDuration)
	accessToken, err := util.GenerateJWTToken(accessTokenClaims, c.config.SigningMethod(), []byte(c.config.JWTSecretKey))
	if err != nil {
		return dto.NoLoginResponse, common.StatusError(status.INTERNAL_SERVER_ERROR)
	}

	// Save Credential
	cred := auth.NewCredential(usr, "Test", accessTokenClaims["id"].(string), refreshToken)
	err = c.authRepo.Create(&cred)
	if err != nil {
		return dto.NoLoginResponse, common.StatusError(status.INTERNAL_SERVER_ERROR)
	}

	return dto.NewLoginResponse(accessToken), common.StatusSuccess()
}

func (c credentialService) RefreshToken(request *dto.RefreshTokenInput) (dto.RefreshTokenResponse, common.Status) {

	if !request.Validate() {
		return dto.NoRefreshTokenResponse, common.StatusError(status.BAD_PARAMETER_ERROR)
	}

	// Parse and validate access token
	keyfunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(c.config.JWTSecretKey), nil
	}
	accessToken, err := jwt.ParseWithClaims(request.AccessToken, &common.AccessTokenClaims{}, keyfunc)

	if !accessToken.Valid {
		var ve *jwt.ValidationError
		ok := errors.As(err, &ve)
		if ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return dto.NoRefreshTokenResponse, common.StatusError(status.JWT_TOKEN_MALFORMED)
			}
		}
	}

	accessTokenClaims, ok := accessToken.Claims.(*common.AccessTokenClaims)
	if !ok {
		return dto.NoRefreshTokenResponse, common.StatusError(status.JWT_TOKEN_MALFORMED)
	}

	// Find credential
	cred, err := c.authRepo.FindByAccessTokenId(accessTokenClaims.Id)
	if err != nil {
		return dto.NoRefreshTokenResponse, common.StatusError(status.ACCESS_TOKEN_WITHOUT_REFRESH_TOKEN)
	}

	// Check user existences
	usr, err := c.userRepo.FindUserById(cred.UserId)
	if err != nil {
		return dto.NoRefreshTokenResponse, common.StatusError(status.USER_NOT_FOUND)
	}

	// Check either the credential token is expired (when the credential token is expired remove it and users should relog)
	refreshToken, err := jwt.Parse(cred.Token, keyfunc)
	if !refreshToken.Valid {
		var ve *jwt.ValidationError
		ok := errors.As(err, &ve)
		cerr := common.StatusError(status.INTERNAL_SERVER_ERROR)
		if ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				cerr = common.StatusError(status.JWT_TOKEN_MALFORMED)
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Remove token
				err = c.authRepo.Remove(cred.UserId, cred.Id)
				if err == nil {
					cerr = common.StatusError(status.ACCESS_TOKEN_EXPIRED)
				}
			}
		}
		return dto.NoRefreshTokenResponse, cerr
	}

	// Generate access token
	newAccessTokenClaims := usr.GenerateAccessTokenClaims(c.config.AccessTokenDuration)
	newAccessToken, err := util.GenerateJWTToken(newAccessTokenClaims, c.config.SigningMethod(), []byte(c.config.JWTSecretKey))
	if err != nil {
		return dto.NoRefreshTokenResponse, common.StatusError(status.USER_NOT_FOUND)
	}

	// Prevent creating many access token from old access token
	err = c.authRepo.UpdateAccessTokenId(cred.Id, newAccessTokenClaims["id"].(string))
	if err != nil {
		return dto.NoRefreshTokenResponse, common.NewRepositoryStatus(err)
	}

	return dto.NewRefreshTokenResponse(newAccessToken), common.StatusSuccess()
}

func (c credentialService) GetCredentials(userId string) ([]dto.CredentialResponse, common.Status) {
	creds, err := c.authRepo.FindUserCredentials(userId)
	if err != nil {
		return nil, common.NewRepositoryStatus(err)
	}

	responses := containers.CastSlicePtr(creds, authMapper.ToCredentialResponse)
	return responses, common.StatusSuccess()
}

func (c credentialService) SelfLogout(userId, accessTokenId string) common.Status {
	err := c.authRepo.RemoveByAccessTokenId(userId, accessTokenId)
	return common.NewRepositoryStatus(err)
}

func (c credentialService) Logout(userId, credId string) common.Status {
	err := c.authRepo.Remove(userId, credId)
	return common.NewRepositoryStatus(err)
}

func (c credentialService) LogoutDevices(userId string) common.Status {
	err := c.authRepo.RemoveUserCredentials(userId)
	return common.NewRepositoryStatus(err)
}
