package service

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"manga-explorer/internal/common"
	"manga-explorer/internal/common/constant"
	"manga-explorer/internal/common/status"
	"manga-explorer/internal/domain/users"
	"manga-explorer/internal/domain/users/dto"
	authMapper "manga-explorer/internal/domain/users/mapper"
	userRepo "manga-explorer/internal/domain/users/repository"
	authService "manga-explorer/internal/domain/users/service"
	"manga-explorer/internal/util"
	"manga-explorer/internal/util/containers"
	"manga-explorer/internal/util/opt"
)

func NewCredential(config *common.Config, credRepo userRepo.IAuthentication, userRepo userRepo.IUser) authService.IAuthentication {
	return &credentialService{config: config, authRepo: credRepo, userRepo: userRepo}
}

type credentialService struct {
	config *common.Config

	authRepo userRepo.IAuthentication
	userRepo userRepo.IUser
}

func (c credentialService) Authenticate(input *dto.LoginInput) (dto.LoginResponse, status.Object) {
	usr, err := c.userRepo.FindUserByEmail(input.Email)
	if err != nil {
		return dto.NoLoginResponse, status.RepositoryError(err, opt.New(status.USER_LOGIN_ERROR))
	}

	if !usr.ValidatePassword(input.Password) {
		return dto.NoLoginResponse, status.Error(status.USER_LOGIN_ERROR)
	}

	// Refresh Token Creation
	refreshToken, err := util.GenerateJWTToken(users.DefaultClaims(c.config.RefreshTokenDuration, constant.IssuerName), c.config.SigningMethod(), []byte(c.config.JWTSecretKey))
	if err != nil {
		return dto.NoLoginResponse, status.InternalError()
	}

	// Access Token Creation
	accessTokenClaims := usr.GenerateAccessTokenClaims(c.config.AccessTokenDuration)
	accessToken, err := util.GenerateJWTToken(accessTokenClaims, c.config.SigningMethod(), []byte(c.config.JWTSecretKey))
	if err != nil {
		return dto.NoLoginResponse, status.InternalError()
	}

	// Save Credential
	cred := users.NewCredential(usr, input.DeviceName, accessTokenClaims["id"].(string), refreshToken)
	err = c.authRepo.Create(&cred)
	if err != nil {
		return dto.NoLoginResponse, status.RepositoryError(err, opt.NullInt)
	}

	return dto.NewLoginResponse(accessToken), status.Success()
}

func (c credentialService) RefreshToken(request *dto.RefreshTokenInput) (dto.RefreshTokenResponse, status.Object) {
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
				return dto.NoRefreshTokenResponse, status.Error(status.JWT_TOKEN_MALFORMED)
			}
		}
	}

	accessTokenClaims, ok := accessToken.Claims.(*common.AccessTokenClaims)
	if !ok {
		return dto.NoRefreshTokenResponse, status.Error(status.JWT_TOKEN_MALFORMED)
	}

	// Find credential
	cred, err := c.authRepo.FindByAccessTokenId(accessTokenClaims.Id)
	if err != nil {
		return dto.NoRefreshTokenResponse, status.RepositoryError(err, opt.New(status.ACCESS_TOKEN_WITHOUT_REFRESH_TOKEN))
	}

	// Check user existences
	usr, err := c.userRepo.FindUserById(cred.UserId)
	if err != nil {
		return dto.NoRefreshTokenResponse, status.RepositoryError(err, opt.New(status.USER_NOT_FOUND))
	}

	// Check either the credential token is expired (when the credential token is expired remove it and users should relog)
	refreshToken, err := jwt.Parse(cred.Token, keyfunc)
	if !refreshToken.Valid {
		var ve *jwt.ValidationError
		ok := errors.As(err, &ve)
		cerr := status.Error(status.INTERNAL_SERVER_ERROR)
		if ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				cerr = status.Error(status.JWT_TOKEN_MALFORMED)
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Remove token
				err = c.authRepo.Remove(cred.UserId, cred.Id)
				if err == nil {
					cerr = status.Error(status.ACCESS_TOKEN_EXPIRED)
				}
			}
		}
		return dto.NoRefreshTokenResponse, cerr
	}

	// Generate access token
	newAccessTokenClaims := usr.GenerateAccessTokenClaims(c.config.AccessTokenDuration)
	newAccessToken, err := util.GenerateJWTToken(newAccessTokenClaims, c.config.SigningMethod(), []byte(c.config.JWTSecretKey))
	if err != nil {
		return dto.NoRefreshTokenResponse, status.InternalError()
	}

	// Prevent creating many access token from old access token
	err = c.authRepo.UpdateAccessTokenId(cred.Id, newAccessTokenClaims["id"].(string))
	if err != nil {
		return dto.NoRefreshTokenResponse, status.RepositoryError(err, opt.New(status.ACCESS_TOKEN_WITHOUT_REFRESH_TOKEN))
	}

	return dto.NewRefreshTokenResponse(newAccessToken), status.Updated()
}

func (c credentialService) GetCredentials(userId string) ([]dto.CredentialResponse, status.Object) {
	creds, err := c.authRepo.FindUserCredentials(userId)
	if err != nil {
		return nil, status.RepositoryError(err, opt.New(status.CREDENTIALS_NOT_FOUND))
	}

	responses := containers.CastSlicePtr(creds, authMapper.ToCredentialResponse)
	return responses, status.Success()
}

func (c credentialService) SelfLogout(userId, accessTokenId string) status.Object {
	err := c.authRepo.RemoveByAccessTokenId(userId, accessTokenId)
	return status.ConditionalRepository(err, status.DELETED, opt.New(status.LOGOUT_CREDENTIAL_NOT_FOUND))
}

func (c credentialService) Logout(userId, credId string) status.Object {
	err := c.authRepo.Remove(userId, credId)
	return status.ConditionalRepository(err, status.DELETED, opt.New(status.LOGOUT_CREDENTIAL_NOT_FOUND))
}

func (c credentialService) LogoutDevices(userId string) status.Object {
	err := c.authRepo.RemoveUserCredentials(userId)
	return status.ConditionalRepository(err, status.DELETED, opt.New(status.LOGOUT_CREDENTIAL_NOT_FOUND))
}
