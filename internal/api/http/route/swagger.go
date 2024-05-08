package route

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "manga-explorer/docs"
)

func NewSwaggerRoute() IArbitraryRoute { return &swaggerRoute{} }

type swaggerRoute struct{}

func (s swaggerRoute) Route(config *Config, router gin.IRouter) {
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, func(config *ginSwagger.Config) {
		config.DefaultModelsExpandDepth = -1
	}))
}
