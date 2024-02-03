package users

import (
	"github.com/gin-gonic/gin"
	"manga-explorer/internal/common"
	"manga-explorer/internal/common/constant"
	"manga-explorer/internal/common/status"
	"manga-explorer/internal/domain/users/dto"
	"manga-explorer/internal/domain/users/service"
	"manga-explorer/internal/util"
	"manga-explorer/internal/util/httputil"
	"manga-explorer/internal/util/httputil/resp"
	"net/http"
)

func NewUserController(userService service.IUser) UserController {
	return UserController{userService: userService}
}

type UserController struct {
	userService service.IUser
}

// GetUserProfile Get user profile
//
//	@Summary		Get User Profile
//	@Description	Get user's details
//	@Tags			users
//	@Produce		json
//	@Param			user_id	path		string	true	"user id"
//	@Success		200		{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=dto.ProfileResponse}}
//	@Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=common.ParameterError}}
//	@Router			/users/{user_id}/profiles [get]
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

	// Response
	if !util.IsUUID(id) {
		stat := status.Error(status.BAD_PARAMETER_ERROR)
		resp.ErrorDetailed(ctx, stat, common.NewParameterError("id", "should be an UUID type"))
		return
	}

	usr, cerr := u.userService.FindUserProfileById(id)
	resp.Conditional(ctx, cerr, usr, nil)
}

// GetUsers Get all registered users
//
//	@Summary		Get Users
//	@Description	Get all registered users
//	@Tags			users
//	@Produce		json
//	@Success		200	{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=[]dto.UserResponse}}
//	@Failure		400	{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
//	@Router			/users [get]
func (u *UserController) GetUsers(ctx *gin.Context) {
	res, stat := u.userService.GetAllUsers()
	resp.Conditional(ctx, stat, res, nil)
}

// AddUser Add new user
//
//	@Summary		Add User
//	@Description	Create user arbitrarily
//	@Tags			users
//	@Accept			json
//
//	@Param			input	body		dto.AddUserInput	true	"add user input"
//
//	@Success		201		{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
//	@Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
//	@Router			/users [post]
func (u *UserController) AddUser(ctx *gin.Context) {
	input := dto.AddUserInput{}
	stat, fieldsErr := httputil.BindJson(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	stat = u.userService.AddUser(&input)
	resp.Conditional(ctx, stat, nil, nil)
}

// Register create new user
//
//	@Summary		Register
//	@Description	Create new user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//
//	@Param			input	body		dto.UserRegisterInput	true	"user registration input"
//
//	@Success		201		{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=dto.UserResponse}}
//	@Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
//	@Router			/users/register [post]
func (u *UserController) Register(ctx *gin.Context) {
	input := dto.UserRegisterInput{}
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

// DeleteUser delete user by the id
//
//	@Summary		Delete User
//	@Description	Delete user
//	@Tags			users
//	@Param			user_id	path		string	true	"user id"
//	@Success		200		{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
//	@Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=common.ParameterError}}
//	@Router			/users/{user_id} [delete]
func (u *UserController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if len(id) == 0 {
		resp.ErrorDetailed(ctx, status.Error(status.BAD_QUERY_ERROR), common.NewNotPresentParameter("id"))
		return
	}

	// Response
	if !util.IsUUID(id) {
		stat := status.Error(status.BAD_PARAMETER_ERROR)
		resp.ErrorDetailed(ctx, stat, common.NewParameterError("id", "should be an UUID type"))
		return
	}

	stat := u.userService.DeleteUser(id)
	resp.Conditional(ctx, stat, nil, nil)
}

// EditUser Edit logged-in user's details
//
//	@Summary		Edit User
//	@Description	Edit logged-in user's details
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			input	body		dto.UserEditInput	true	"user edit input"
//	@Success		200		{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
//	@Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
//	@Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
//	@Router			/users [put]
func (u *UserController) EditUser(ctx *gin.Context) {
	input := dto.UserEditInput{}
	stat, fieldsErr := httputil.BindJson(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	if stat = input.Status(); stat.IsError() {
		resp.Error(ctx, stat)
		return
	}

	claims, stat := common.GetClaims(ctx)
	if stat.IsError() {
		resp.Error(ctx, stat)
		return
	}
	input.UserId = claims.UserId

	stat = u.userService.UpdateUser(&input)
	resp.Conditional(ctx, stat, nil, nil)
}

// EditUserExtended Edit user in extended fields (password) without checking
//
//	@Summary		Edit User Extended
//	@Description	Edit user in extended fields (password) without checking
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			input	body		dto.UserEditExtendedInput	true	"user edit extended input"
//	@Param			user_id	path		uuid.UUID					true	"user id"
//	@Success		200		{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
//	@Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
//	@Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
//	@Failure		500		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
//	@Router			/users/{user_id} [put]
func (u *UserController) EditUserExtended(ctx *gin.Context) {
	input := dto.UserEditExtendedInput{}
	input.ConstructURI(ctx)
	stat, fieldsErr := httputil.BindJson(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	stat = u.userService.UpdateUserExtended(&input)
	resp.Conditional(ctx, stat, nil, nil)
}

// EditUserProfile Edit user logged-in user's profile
//
//	@Summary		Edit User Profile
//	@Description	Edit user logged-in user's profile
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			input	body		dto.ProfileEditInput	true	"profile edit input"
//	@Success		200		{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
//	@Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
//	@Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
//	@Router			/users/profiles [put]
func (u *UserController) EditUserProfile(ctx *gin.Context) {
	input := dto.ProfileEditInput{}
	stat, fieldsErr := httputil.BindJson(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	claims, stat := common.GetClaims(ctx)
	if stat.IsError() {
		resp.Error(ctx, stat)
		return
	}
	input.UserId = claims.UserId

	stat = u.userService.UpdateProfile(&input)
	resp.Conditional(ctx, stat, nil, nil)
}

// EditUserProfileExtended edit user's profile on extended fields
//
//	@Summary		Edit User Profile Extended
//	@Description	Edit user's profile on extended fields
//	@Tags			users
//	@Accept			json
//	@Produce		json
//
//	@Param			input	body		dto.ProfileEditExtendedInput	true	"user profile edit input"
//	@Param			user_id	path		uuid.UUID						true	"user id"
//
//	@Success		200		{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
//	@Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
//	@Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
//	@Router			/users/{user_id}/profiles [put]
func (u *UserController) EditUserProfileExtended(ctx *gin.Context) {
	input := dto.ProfileEditExtendedInput{}
	input.ConstructURI(ctx)
	stat, fieldsErr := httputil.BindJson(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	stat = u.userService.UpdateProfileExtended(&input)
	resp.Conditional(ctx, stat, nil, nil)
}

// UpdateProfileImage Edit logged-in user profile image
//
//	@Summary		Edit Profile Image
//	@Description	Edit logged-in user profile image
//	@Tags			users
//	@Accept			mpfd
//	@Produce		json
//	@Param			image	formData	file	true	"profile's image"
//	@Success		200		{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
//	@Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
//	@Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
//	@Router			/users/profiles/image [put]
//
//	@Router			/users/profiles/image [delete]
func (u *UserController) UpdateProfileImage(ctx *gin.Context) {
	input := dto.ProfileImageUpdateInput{}

	claims, stat := common.GetClaims(ctx)
	if stat.IsError() {
		resp.Error(ctx, stat)
		return
	}

	if ctx.Request.Method == http.MethodDelete {
		stat = u.userService.DeleteProfileImage(claims.UserId)
	} else if ctx.Request.Method == http.MethodPut {
		var fieldErrors []common.FieldError
		input.UserId = claims.UserId
		stat, fieldErrors = httputil.BindMultipartForm(ctx, &input)
		if stat.IsError() {
			resp.ErrorDetailed(ctx, stat, fieldErrors)
			return
		}

		stat = u.userService.UpdateProfileImage(&input)
	}
	resp.Conditional(ctx, stat, nil, nil)
}

// @Summary		Change Password
// @Description	Change logged-in user's password
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			input	body		dto.ChangePasswordInput	true	"change password input"
// @Success		200		{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
// @Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
// @Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Failure		500		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/users/change-password [post]
func (u *UserController) ChangePassword(ctx *gin.Context) {
	input := dto.ChangePasswordInput{}
	stat, fieldsErr := httputil.BindJson(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	claims, stat := common.GetClaims(ctx)
	if stat.IsError() {
		resp.Error(ctx, stat)
		return
	}
	input.UserId = claims.UserId

	stat = u.userService.ChangePassword(&input)
	resp.Conditional(ctx, stat, nil, nil)
}

// @Summary		Request Reset Password
// @Description	Request to reset password by user's email
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			input	body		dto.ChangePasswordInput	true	"change password input"
// @Success		200		{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
// @Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
// @Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Failure		500		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/users/reset-password [post]
func (u *UserController) RequestResetPassword(ctx *gin.Context) {
	input := dto.ResetPasswordRequestInput{}
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

// @Summary		Reset Password
// @Description	Reset password after request
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			input	body		dto.ResetPasswordInput	true	"reset password input"
// @Param			token	path		string					true	"token"
// @Success		200		{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
// @Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
// @Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Failure		500		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/users/reset-password/{token} [post]
// @Router			/users/reset-password/{token} [get]
func (u *UserController) ResetPassword(ctx *gin.Context) {
	input := dto.ResetPasswordInput{}
	input.ConstructURI(ctx)
	stat, fieldsErr := httputil.BindJson(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	// Reset user password
	stat = u.userService.ResetPassword(&input)
	resp.Conditional(ctx, stat, nil, nil)
}

// @Summary		Request Verify Email
// @Description	Request for email verification
// @Tags			users
// @Produce		json
// @Success		200	{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
// @Failure		400	{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
// @Failure		400	{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Failure		500	{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/users/email-verif [post]
func (u *UserController) RequestVerifyEmail(ctx *gin.Context) {
	input := dto.VerifEmailRequestInput{}
	claims, stat := common.GetClaims(ctx)
	if stat.IsError() {
		resp.Error(ctx, stat)
		return
	}
	input.UserId = claims.UserId

	stat = u.userService.RequestEmailVerification(&input)
	resp.Conditional(ctx, stat, nil, nil)
}

// @Summary		Verify Email
// @Description	Verify email based on the token
// @Tags			users
// @Produce		json
// @Param			token	path		string	true	"token"
// @Success		200		{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
// @Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
// @Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Failure		500		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/users/email-verif/{token} [post]
// @Router			/users/email-verif/{token} [get]
func (u *UserController) VerifyEmail(ctx *gin.Context) {
	input := dto.VerifyEmailInput{}
	input.ConstructURI(ctx)

	stat := u.userService.VerifyEmail(&input)
	resp.Conditional(ctx, stat, nil, nil)
}
