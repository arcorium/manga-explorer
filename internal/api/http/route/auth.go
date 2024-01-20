package route

import (
	"github.com/gin-gonic/gin"
)

func NewAuthRoute() IRoute {
	return &authRoute{}
}

type authRoute struct {
}

func (a authRoute) V1Route(config *Config, router gin.IRouter) {
	authRouter := router.Group("/auth/")
	controller := &config.Controller.Auth

	authRouter.POST("/refresh-token", controller.RefreshToken)
	authRouter.POST("/login", controller.Login)
	authRouter.POST("/register", controller.Register)
	authRouter.POST("/reset-password", controller.RequestResetPassword)
	authRouter.PATCH("/reset-password/:token", controller.ResetPassword)

	authRouter.Use(config.Middleware.Auth.Handle)
	authRouter.GET("/logout/all", controller.LogoutAllDevice)
	authRouter.GET("/logout/:id", controller.Logout)
	authRouter.GET("/devices", controller.GetCredentials)
}
