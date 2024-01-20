package authentication

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"manga-explorer/internal/app/common"
	"manga-explorer/internal/app/common/status"
	"manga-explorer/internal/app/service/utility/mail"
	"manga-explorer/internal/domain/auth/dto"
	authService "manga-explorer/internal/domain/auth/service"
	"manga-explorer/internal/domain/users"
	dto2 "manga-explorer/internal/domain/users/dto"
	userService "manga-explorer/internal/domain/users/service"
	"manga-explorer/internal/util/httputil"
)

func NewAuthController(userService userService.IUser, credService authService.IAuthentication, verifService userService.IVerification, mailService mail.IService) AuthController {
	return AuthController{
		userService:  userService,
		authService:  credService,
		verifService: verifService,
		mailService:  mailService,
	}
}

type AuthController struct {
	userService  userService.IUser
	authService  authService.IAuthentication
	verifService userService.IVerification

	mailService mail.IService
}

// Login sign in account
// @Summary login
// @Description login to get token
// @Tags auth, account
// @Accept json
// @Produce json
// @Success 200 {object} common.Response
// @Failure 400 {object} common.Response
// @Router /login [post]
func (a AuthController) Login(ctx *gin.Context) {
	var req dto.LoginInput
	if err := ctx.BindJSON(&req); err != nil {
		verr := err.(validator.ValidationErrors)
		httputil.ErrorResponse(ctx, common.StatusError(status.BAD_BODY_REQUEST_ERROR), common.GetFieldsError(verr))
		return
	}
	resp, cerr := a.authService.Authenticate(&req)
	httputil.Response(ctx, cerr, resp)
}

func (a AuthController) Register(ctx *gin.Context) {
	var input dto2.UserRegisterInput
	if err := ctx.BindJSON(&input); err != nil {
		httputil.ErrorResponse(ctx, common.StatusError(status.BAD_BODY_REQUEST_ERROR))
		return
	}

	userResponse, cerr := a.userService.RegisterUser(&input)
	if cerr.IsError() {
		httputil.ErrorResponse(ctx, cerr)
		return
	}

	// Verify Email
	verifResponse, cerr := a.verifService.Request(userResponse.Id, users.UsageVerifyEmail)
	if cerr.IsError() {
		httputil.ErrorResponse(ctx, cerr)
		return
	}

	go func() {
		ml := mail.Mail{
			From:    "Me",
			To:      userResponse.Email,
			Subject: "Email-Verification",
			Body:    verifResponse.Token,
		}
		err := a.mailService.SendEmail(&ml)
		if err != nil {
			log.Println("Error Sending Email: ", err)
		}
	}()

	httputil.SuccessResponse(ctx, cerr, userResponse)
}

// Logout Handle user logout, when the uri provide id parameter, those credential will be removed otherwise
// current credential will be removed
func (a AuthController) Logout(ctx *gin.Context) {
	// Check parameter
	data, ok := ctx.Get("claims")
	if !ok {
		httputil.ErrorResponse(ctx, common.StatusError(status.AUTH_UNAUTHORIZED))
		return
	}
	token := data.(*common.AccessTokenClaims)

	var status common.Status
	credId := ctx.Param("id")
	if len(credId) == 0 {
		status = a.authService.SelfLogout(token.UserId, token.Id)
	} else {
		status = a.authService.Logout(token.UserId, credId)
	}
	httputil.Response(ctx, status, nil)
}

func (a AuthController) LogoutAllDevice(ctx *gin.Context) {
	data, ok := ctx.Get("claims")
	if !ok {
		httputil.ErrorResponse(ctx, common.StatusError(status.AUTH_UNAUTHORIZED))
		return
	}
	token := data.(*common.AccessTokenClaims)
	cerr := a.authService.LogoutDevices(token.UserId)
	httputil.Response(ctx, cerr, nil)
}

func (a AuthController) RequestResetPassword(ctx *gin.Context) {
	var input dto.ResetPasswordRequestInput
	if err := ctx.BindJSON(&input); err != nil {
		httputil.ErrorResponse(ctx, common.StatusError(status.BAD_BODY_REQUEST_ERROR))
		return
	}

	userResponse, cerr := a.userService.FindUserByEmail(input.Email)
	if cerr.IsError() {
		httputil.ErrorResponse(ctx, cerr)
		return
	}

	verifResponse, cerr := a.verifService.Request(userResponse.Id, users.UsageResetPassword)
	if cerr.IsError() {
		httputil.ErrorResponse(ctx, cerr)
		return
	}

	go func() {
		ml := mail.Mail{
			From:    "Me",
			To:      userResponse.Email,
			Subject: "Reset-Password",
			Body:    verifResponse.Token,
		}
		a.mailService.SendEmail(&ml)
	}()

	httputil.SuccessResponseMessage(ctx, cerr, common.MSG_SUCCESS_REQUEST_RESET_PASSWORD)
}

func (a AuthController) ResetPassword(ctx *gin.Context) {
	token := ctx.Param("token")
	if len(token) == 0 {
		httputil.ErrorResponse(ctx, common.StatusError(status.BAD_PARAMETER_ERROR))
		return
	}

	var input dto2.ResetPasswordInput
	if err := ctx.BindJSON(&input); err != nil {
		httputil.ErrorResponse(ctx, common.StatusError(status.BAD_BODY_REQUEST_ERROR))
		return
	}

	// Find token
	verifResponse, cerr := a.verifService.Find(token)
	if cerr.IsError() {
		httputil.ErrorResponse(ctx, cerr)
		return
	}
	input.UserId = verifResponse.UserId

	// Validate token
	cerr = a.verifService.Validate(&verifResponse)
	if cerr.IsError() {
		httputil.ErrorResponse(ctx, cerr)
		return
	}

	// Immediately revoke verify token
	cerr = a.verifService.Remove(token)
	if cerr.IsError() {
		httputil.ErrorResponse(ctx, cerr)
		return
	}

	// Reset user password
	cerr = a.userService.ResetPassword(&input)
	if cerr.IsError() {
		httputil.ErrorResponse(ctx, cerr)
		return
	}

	// SelfLogout user devices
	cerr = a.authService.LogoutDevices(verifResponse.UserId)
	httputil.Response(ctx, cerr, nil)
}

func (a AuthController) GetCredentials(ctx *gin.Context) {
	data, ok := ctx.Get("claims")
	if !ok {
		httputil.ErrorResponse(ctx, common.StatusError(status.AUTH_UNAUTHORIZED))
		return
	}
	token := data.(*common.AccessTokenClaims)
	creds, cerr := a.authService.GetCredentials(token.UserId)
	httputil.Response(ctx, cerr, creds)
}

func (a AuthController) RefreshToken(ctx *gin.Context) {
	var req dto.RefreshTokenInput
	if err := ctx.BindJSON(&req); err != nil {
		log.Println(err)
		httputil.ErrorResponse(ctx, common.StatusError(status.BAD_BODY_REQUEST_ERROR))
		return
	}

	token, cerr := a.authService.RefreshToken(&req)
	httputil.Response(ctx, cerr, token)
}
