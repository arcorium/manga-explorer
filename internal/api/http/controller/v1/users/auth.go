package users

import (
  "github.com/gin-gonic/gin"
  "manga-explorer/internal/common"
  "manga-explorer/internal/common/constant"
  "manga-explorer/internal/common/status"
  "manga-explorer/internal/domain/users"
  "manga-explorer/internal/domain/users/dto"
  authService "manga-explorer/internal/domain/users/service"
  "manga-explorer/internal/util"
  "manga-explorer/internal/util/httputil"
  "manga-explorer/internal/util/httputil/resp"
  "strings"
)

func NewAuthController(credService authService.IAuthentication) AuthController {
  return AuthController{
    authService: credService,
  }
}

type AuthController struct {
  authService authService.IAuthentication
}

// Login sign in account
//
//	@Summary		Login
//	@Description	login to authenticate
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			input	body		dto.LoginInput	true	"user login input"
//	@Success		200		{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=dto.LoginResponse}}
//	@Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
//	@Failure		500		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
//	@Router			/auth/login [post]
func (a AuthController) Login(ctx *gin.Context) {
  input := dto.LoginInput{}
  stat, errFields := httputil.BindJson(ctx, &input)
  if stat.IsError() {
    resp.ErrorDetailed(ctx, stat, errFields)
    return
  }

  // Get device name from context
  usr, err := util.GetContextValue[*users.Device](ctx, constant.UserAgentKey)
  if err != nil {
    resp.Error(ctx, status.InternalError())
    return
  }
  input.DeviceName = usr.Name

  res, stat := a.authService.Authenticate(&input)
  resp.Conditional(ctx, stat, &res, nil)
}

// Logout Handle user logout, when the uri provide id parameter, those credential will be removed otherwise
// current credential will be removed
//
//	@Summary		Logout
//	@Description	logout credential id or current logged in credential
//	@Tags			auth
//	@Produce		json
//	@Param			cred_id	path		uuid.UUID	false	"credential id"
//	@Success		200		{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
//	@Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=common.ParameterError}}
//	@Failure		500		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
//	@Router			/auth/logout/{cred_id} [post]
func (a AuthController) Logout(ctx *gin.Context) {
  token, stat := common.GetClaims(ctx)
  if stat.IsError() {
    resp.Error(ctx, stat)
    return
  }
  // Check parameter
  credId := ctx.Param("id")
  credId = strings.Trim(credId, "/")
  if len(credId) == 0 {
    stat = a.authService.SelfLogout(token.UserId, token.Id)
  } else {
    if !util.IsUUID(credId) {
      stat = status.Error(status.BAD_PARAMETER_ERROR)
      resp.ErrorDetailed(ctx, stat, common.NewParameterError("id", "should be an UUID type"))
      return
    }
    stat = a.authService.Logout(token.UserId, credId)
  }
  resp.Conditional(ctx, stat, nil, nil)
}

// LogoutAllDevice Remove all credentials associated on logged-in user
//
//	@Summary		Logout All Device
//	@Description	Remove all current user logged in credentials
//	@Tags			auth
//	@Produce		json
//	@Success		200	{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
//	@Failure		400	{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
//	@Router			/auth/logouts [post]
func (a AuthController) LogoutAllDevice(ctx *gin.Context) {
  token, stat := common.GetClaims(ctx)
  if stat.IsError() {
    resp.Error(ctx, stat)
    return
  }
  stat = a.authService.LogoutDevices(token.UserId)
  resp.Conditional(ctx, stat, nil, nil)
}

// RefreshToken Get new access token and recreate refresh token
//
//	@Summary		Refresh Token
//	@Description	Recreate access token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			input	body		dto.RefreshTokenInput	true	"refresh token input, token type should be Bearer"
//	@Success		200		{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=dto.RefreshTokenResponse}}
//	@Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
//	@Failure		500		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
//	@Router			/auth/refresh-token [post]
func (a AuthController) RefreshToken(ctx *gin.Context) {
  input := dto.RefreshTokenInput{}
  stat, fieldErrors := httputil.BindJson(ctx, &input)
  if stat.IsError() {
    resp.ErrorDetailed(ctx, status.Error(status.BAD_REQUEST_ERROR), fieldErrors)
    return
  }

  token, stat := a.authService.RefreshToken(&input)
  resp.Conditional(ctx, stat, token, nil)
}

// GetCredentials Get all user credentials or session
//
//	@Summary		Get Credentials
//	@Description	Get all user credentials
//	@Tags			auth
//	@Produce		json
//	@Success		200	{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=[]dto.CredentialResponse}}
//	@Failure		400	{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
//	@Router			/auth/devices [get]
func (a AuthController) GetCredentials(ctx *gin.Context) {
  token, stat := common.GetClaims(ctx)
  if stat.IsError() {
    resp.Error(ctx, stat)
    return
  }
  creds, stat := a.authService.GetCredentials(token.UserId)
  resp.Conditional(ctx, stat, creds, nil)
}
