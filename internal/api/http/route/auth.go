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

	authRouter.POST("/login", config.Middleware.UserAgentChecker.Handle, controller.Login)
	authRouter.POST("/refresh-token", controller.RefreshToken)

	authRouter.Use(config.Middleware.Authorization.Handle)
	authRouter.POST("/logout/*id", controller.Logout)
	authRouter.POST("/logouts", controller.LogoutAllDevice)
	authRouter.GET("/devices", controller.GetCredentials)
}
