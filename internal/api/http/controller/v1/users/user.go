package users

import (
	"github.com/gin-gonic/gin"
	"manga-explorer/internal/app/common"
	"manga-explorer/internal/app/common/constant"
	"manga-explorer/internal/app/common/status"
	"manga-explorer/internal/domain/users/dto"
	"manga-explorer/internal/domain/users/service"
	"manga-explorer/internal/util"
	"manga-explorer/internal/util/httputil"
	"manga-explorer/internal/util/httputil/resp"
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
		claims, stat := common.GetClaims(ctx)
		if stat.IsError() {
			resp.Error(ctx, stat)
			return
		}
		id = claims.UserId
	}

	// Validate
	if !util.IsUUID(id) {
		stat := status.Error(status.BAD_PARAMETER_ERROR)
		resp.ErrorDetailed(ctx, stat, common.NewParameterError("id", "should be an UUID type"))
		return
	}

	usr, cerr := u.userService.FindUserProfileById(id)
	resp.Conditional(ctx, cerr, usr, nil)
}

func (u *UserController) GetUserProfiles(ctx *gin.Context) {
	res, stat := u.userService.GetAllUsers()
	resp.Conditional(ctx, stat, res, nil)
}

func (u *UserController) AddUser(ctx *gin.Context) {
	var input dto.AddUserInput
	stat, fieldsErr := httputil.BindUriJson(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	stat = u.userService.AddUser(&input)
	resp.Conditional(ctx, stat, nil, nil)
}

func (u *UserController) Register(ctx *gin.Context) {
	var input dto.UserRegisterInput
	stat, fieldsErr := httputil.BindJson(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	userResponse, stat := u.userService.RegisterUser(&input)
	if stat.IsError() {
		resp.Error(ctx, stat)
		return
	}

	resp.Success(ctx, stat, userResponse, nil)
}

func (u *UserController) DeleteUser(ctx *gin.Context) {
	id, present := ctx.GetQuery("id")
	if !present {
		resp.ErrorDetailed(ctx, status.Error(status.BAD_QUERY_ERROR), common.NewNotPresentParameter("id"))
		return
	}

	// Validate
	if !util.IsUUID(id) {
		stat := status.Error(status.BAD_PARAMETER_ERROR)
		resp.ErrorDetailed(ctx, stat, common.NewParameterError("id", "should be an UUID type"))
		return
	}

	stat := u.userService.DeleteUser(id)
	resp.Conditional(ctx, stat, nil, nil)
}

func (u *UserController) EditUser(ctx *gin.Context) {
	var input dto.UpdateUserInput
	stat, fieldsErr := httputil.BindAuthorizedJSON(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	stat = u.userService.UpdateUser(&input)
	resp.Conditional(ctx, stat, nil, nil)
}

func (u *UserController) UpdateUserExtended(ctx *gin.Context) {
	var input dto.UpdateUserExtendedInput
	stat, fieldsErr := httputil.BindUriJson(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	stat = u.userService.UpdateUserExtended(&input)
	resp.Conditional(ctx, stat, nil, nil)
}

func (u *UserController) EditUserProfile(ctx *gin.Context) {
	var input dto.ProfileUpdateInput
	stat, fieldsErr := httputil.BindAuthorizedJSON(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	stat = u.userService.UpdateProfile(&input)
	resp.Conditional(ctx, stat, nil, nil)
}

func (u *UserController) UpdateUserProfileExtended(ctx *gin.Context) {
	var input dto.UpdateProfileExtendedInput
	stat, fieldsErr := httputil.BindUriJson(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	stat = u.userService.UpdateProfileExtended(&input)
	resp.Conditional(ctx, stat, nil, nil)
}

func (u *UserController) ChangePassword(ctx *gin.Context) {
	var input dto.ChangePasswordInput
	stat, fieldsErr := httputil.BindAuthorizedJSON(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	stat = u.userService.ChangePassword(&input)
	resp.Conditional(ctx, stat, nil, nil)
}

func (u *UserController) RequestResetPassword(ctx *gin.Context) {
	var input dto.ResetPasswordRequestInput
	stat, fieldsErr := httputil.BindJson(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	stat = u.userService.RequestResetPassword(&input)
	if stat.IsError() {
		resp.Error(ctx, stat)
		return
	}
	resp.SuccessMessage(ctx, stat, constant.MSG_SUCCESS_REQUEST_RESET_PASSWORD)
}

func (u *UserController) ResetPassword(ctx *gin.Context) {
	var input dto.ResetPasswordInput
	stat, fieldsErr := httputil.BindUriJson(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	// Reset user password
	stat = u.userService.ResetPassword(&input)
	if stat.IsError() {
		resp.Error(ctx, stat)
		return
	}

	// SelfLogout user devices
	resp.Conditional(ctx, stat, nil, nil)
}
