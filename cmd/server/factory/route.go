package factory

import (
  "github.com/gin-gonic/gin"
  mangaController "manga-explorer/internal/api/http/controller/v1/mangas"
  userController "manga-explorer/internal/api/http/controller/v1/users"
  "manga-explorer/internal/api/http/middleware"
  "manga-explorer/internal/api/http/route"
  "manga-explorer/internal/common"
  "manga-explorer/internal/common/constant"
  "manga-explorer/internal/domain/users"
)

func CreateRouter(config *common.Config, service *Service, router gin.IRouter) route.Router {
  controllerConfig := route.ConfigController{
    Auth:         userController.NewAuthController(service.Authentication),
    User:         userController.NewUserController(service.User),
    Manga:        mangaController.NewMangaController(service.Manga),
    MangaChapter: mangaController.NewChapterController(service.Chapter),
    MangaGenre:   mangaController.NewGenreController(service.Genre),
  }

  middlewareConfig := route.ConfigMiddleware{
    Authorization: middleware.NewAuthMiddleware(config.JWTSecretKey, &middleware.AuthMiddlewareConfig{
      SigningMethod: config.SigningMethod(),
      ClaimsKey:     constant.ClaimsKey,
    }),
    AdminRestrict: middleware.NewRoleRestrictionMiddleware(&middleware.RoleRestrictionConfig{
      ClaimsKey: constant.ClaimsKey,
    }, []users.Role{users.RoleAdmin}),
    UserAgentChecker: middleware.NewUserAgentChecker(true, middleware.UserAgentConfig{Key: constant.UserAgentKey}),
  }
  routerConfig := route.Config{
    Controller: controllerConfig,
    Middleware: middlewareConfig,
  }

  return route.NewRouter(&routerConfig, router)
}
