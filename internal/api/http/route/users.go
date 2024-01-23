package route

import (
	"github.com/gin-gonic/gin"
)

func NewUserRoute() IRoute {
	return &userRoute{}
}

type userRoute struct {
}

func (u userRoute) V1Route(config *Config, router gin.IRouter) {
	user := router.Group("/users")
	admin := user.Group("/", config.Middleware.Authorization.Handle, config.Middleware.AdminRestrict.Handle)

	userController := &config.Controller.User

	user.POST("/register", userController.Register)
	user.POST("/reset-password", userController.RequestResetPassword)
	user.PATCH("/reset-password/:token", userController.ResetPassword)
	user.POST("/change-password", userController.ChangePassword)

	user.Use(config.Middleware.Authorization.Handle)
	user.PUT("/", userController.EditUser)
	// Profiles thingies
	user.GET("/:id/profiles", userController.GetUserProfile)
	user.GET("/profiles", userController.GetUserProfile)
	user.PUT("/profiles", userController.EditUserProfile)
	user.PUT("/profiles/image", userController.UpdateProfileImage)
	user.DELETE("/profiles/image", userController.UpdateProfileImage)
	// Admin
	admin.PUT("/:id", userController.UpdateUserExtended)
	admin.PUT("/:id/profiles", userController.UpdateUserProfileExtended)
	admin.POST("/", userController.AddUser)
	admin.DELETE("/:id", userController.DeleteUser)
	admin.GET("/", userController.GetUsers)
}
