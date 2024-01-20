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
	user := router.Group("/users", config.Middleware.Auth.Handle)
	admin := user.Group("/", config.Middleware.AdminRestrict.Handle)

	userController := &config.Controller.User

	user.PUT("/", userController.EditUser)
	user.POST("/change-password", userController.ChangePassword)

	// Profiles thingies
	user.GET("/:id/profiles", userController.GetUserProfile)
	user.PUT("/profiles", userController.EditUserProfile)
	// Admin
	admin.PUT("/:id", userController.UpdateUserExtended)
	admin.PUT("/:id/profiles", userController.UpdateUserProfileExtended)
	admin.POST("/", userController.AddUser)
	admin.DELETE("/:id", userController.DeleteUser)
	admin.GET("/", userController.GetUserProfiles)
}
