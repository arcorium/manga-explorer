package users

import (
	"github.com/gin-gonic/gin"
	"manga-explorer/internal/app/common"
	"manga-explorer/internal/app/common/status"
	dto2 "manga-explorer/internal/domain/users/dto"
	"manga-explorer/internal/domain/users/service"
	"manga-explorer/internal/util"
	"manga-explorer/internal/util/httputil"
)

func NewUserController(userService service.IUser) UserController {
	return UserController{userService: userService}
}

type UserController struct {
	userService service.IUser
}

func (u *UserController) GetUserProfile(ctx *gin.Context) {
	id := ctx.Param("id")
	if len(id) == 0 {
		claims, err := util.GetContextValue[*common.AccessTokenClaims](ctx, common.ClaimsKey)
		if err != nil {
			httputil.ErrorResponse(ctx, common.StatusError(status.AUTH_UNAUTHORIZED))
			return
		}
		id = claims.UserId
	}

	usr, cerr := u.userService.FindUserProfileById(id)
	httputil.Response(ctx, cerr, usr)
}

func (u *UserController) GetUserProfiles(ctx *gin.Context) {
	res, status := u.userService.GetAllUsers()
	httputil.Response(ctx, status, res)
}

func (u *UserController) AddUser(ctx *gin.Context) {
	input, status := httputil.BindUriJson[dto2.AddUserInput](ctx)
	if status.IsError() {
		httputil.ErrorResponse(ctx, status)
		return
	}

	status = u.userService.AddUser(&input)
	httputil.Response(ctx, status, nil)
}

func (u *UserController) DeleteUser(ctx *gin.Context) {
	id, present := ctx.GetQuery("id")
	if !present {
		httputil.ErrorResponse(ctx, common.StatusError(status.BAD_PARAMETER_ERROR))
		return
	}

	status := u.userService.DeleteUser(id)
	httputil.Response(ctx, status, nil)
}

func (u *UserController) EditUser(ctx *gin.Context) {
	var input dto2.UpdateUserInput
	status := httputil.BindAuthorizedJSON(ctx, &input)
	if status.IsError() {
		httputil.ErrorResponse(ctx, status)
		return
	}

	status = u.userService.UpdateUser(&input)
	httputil.Response(ctx, status, nil)
}

func (u *UserController) UpdateUserExtended(ctx *gin.Context) {
	input, status := httputil.BindUriJson[dto2.UpdateUserExtendedInput](ctx)
	if status.IsError() {
		httputil.ErrorResponse(ctx, status)
		return
	}

	status = u.userService.UpdateUserExtended(&input)
	httputil.Response(ctx, status, nil)
}

func (u *UserController) EditUserProfile(ctx *gin.Context) {
	var input dto2.ProfileUpdateInput
	status := httputil.BindAuthorizedJSON(ctx, &input)
	if status.IsError() {
		httputil.ErrorResponse(ctx, status)
		return
	}

	status = u.userService.UpdateProfile(&input)
	httputil.Response(ctx, status, nil)
}

func (u *UserController) UpdateUserProfileExtended(ctx *gin.Context) {
	input, status := httputil.BindUriJson[dto2.UpdateProfileExtendedInput](ctx)
	if status.IsError() {
		httputil.ErrorResponse(ctx, status)
		return
	}

	status = u.userService.UpdateProfileExtended(&input)
	httputil.Response(ctx, status, nil)
}

func (u *UserController) ChangePassword(ctx *gin.Context) {
	var input dto2.ChangePasswordInput
	status := httputil.BindAuthorizedJSON(ctx, &input)
	if status.IsError() {
		httputil.ErrorResponse(ctx, status)
		return
	}

	status = u.userService.ChangePassword(&input)
	httputil.Response(ctx, status, nil)
}
