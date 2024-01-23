package users

import (
	"github.com/gin-gonic/gin"
	"manga-explorer/internal/app/common"
	"manga-explorer/internal/app/common/constant"
	"manga-explorer/internal/app/common/status"
	"manga-explorer/internal/domain/users"
	"manga-explorer/internal/domain/users/dto"
	authService "manga-explorer/internal/domain/users/service"
	"manga-explorer/internal/util"
	"manga-explorer/internal/util/httputil"
	"manga-explorer/internal/util/httputil/resp"
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
// @Summary login
// @Description login to get token
// @Tags auth, account
// @Accept json
// @Produce json
// @Success 200 {object} common.SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Router /login [post]
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
	resp.Conditional(ctx, stat, res, nil)
}

// Logout Handle user logout, when the uri provide id parameter, those credential will be removed otherwise
// current credential will be removed
func (a AuthController) Logout(ctx *gin.Context) {
	token, stat := common.GetClaims(ctx)
	if stat.IsError() {
		resp.Error(ctx, stat)
		return
	}
	// Check parameter
	credId := ctx.Param("id")
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

func (a AuthController) LogoutAllDevice(ctx *gin.Context) {
	token, stat := common.GetClaims(ctx)
	if stat.IsError() {
		resp.Error(ctx, stat)
		return
	}
	stat = a.authService.LogoutDevices(token.UserId)
	resp.Conditional(ctx, stat, nil, nil)
}

func (a AuthController) RefreshToken(ctx *gin.Context) {
	input := dto.RefreshTokenInput{}
	stat, fieldErrors := httputil.BindJson(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, status.Error(status.BAD_BODY_REQUEST_ERROR), fieldErrors)
		return
	}

	token, stat := a.authService.RefreshToken(&input)
	resp.Conditional(ctx, stat, token, nil)
}

func (a AuthController) GetCredentials(ctx *gin.Context) {
	token, stat := common.GetClaims(ctx)
	if stat.IsError() {
		resp.Error(ctx, stat)
		return
	}
	creds, stat := a.authService.GetCredentials(token.UserId)
	resp.Conditional(ctx, stat, creds, nil)
}