package route

import (
	"github.com/gin-gonic/gin"
	"manga-explorer/internal/api/http/controller/v1/mangas"
	"manga-explorer/internal/api/http/controller/v1/users"
	"manga-explorer/internal/api/http/middleware"
)

type ConfigController struct {
	Auth         users.AuthController
	User         users.UserController
	Manga        mangas.MangaController
	MangaChapter mangas.ChapterController
	MangaGenre   mangas.GenreController
}

type ConfigMiddleware struct {
	Authorization    middleware.AuthMiddleware
	UserAgentChecker middleware.UserAgentCheckerMiddleware
	AdminRestrict    middleware.RoleRestrictionMiddleware
}

type Config struct {
	Controller ConfigController
	Middleware ConfigMiddleware
}

func NewRouter(config *Config, router gin.IRouter) Router {
	return Router{Config: config, Router: router}
}

const (
	RouterVersion1 = "/api/v1"
)

type Router struct {
	Config *Config
	Router gin.IRouter
}

func (r *Router) Routes(routes ...IRoute) {
	router := r.Router.Group(RouterVersion1)

	for _, route := range routes {
		route.V1Route(r.Config, router)
	}
}

func (r *Router) ArbitraryRoutes(routes ...IArbitraryRoute) {
	router := r.Router.Group("/api")

	for _, route := range routes {
		route.Route(r.Config, router)
	}
}

type IRoute interface {
	IV1Route
}

type IV1Route interface {
	V1Route(config *Config, router gin.IRouter)
}

type IArbitraryRoute interface {
	Route(config *Config, router gin.IRouter)
}
