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

  user.POST("/reset-password/:token", userController.ResetPassword)
  user.GET("/reset-password/:token", userController.ResetPassword)

  user.POST("/email-verif/:token", userController.VerifyEmail)
  user.GET("/email-verif/:token", userController.VerifyEmail)

  user.Use(config.Middleware.Authorization.Handle)
  user.PUT("/", userController.EditUser)
  user.POST("/change-password", userController.ChangePassword)
  user.POST("/email-verif", userController.RequestVerifyEmail)
  // Profiles thingies
  user.GET("/:id/profiles", userController.GetUserProfile)
  user.GET("/profiles", userController.GetUserProfile)
  user.PUT("/profiles", userController.EditUserProfile)
  user.PATCH("/profiles/image", userController.UpdateProfileImage)
  user.DELETE("/profiles/image", userController.UpdateProfileImage)
  // Admin
  admin.PUT("/:id", userController.EditUserExtended)
  admin.PUT("/:id/profiles", userController.EditUserProfileExtended)
  admin.POST("/", userController.AddUser)
  admin.DELETE("/:id", userController.DeleteUser)
  admin.GET("/", userController.GetUsers)
}
