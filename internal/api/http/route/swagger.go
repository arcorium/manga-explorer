package route

import (
  "github.com/gin-gonic/gin"
  swaggerFiles "github.com/swaggo/files"
  ginSwagger "github.com/swaggo/gin-swagger"
  "manga-explorer/docs"
)

func NewSwaggerRoute() IRoute { return &swaggerRoute{} }

type swaggerRoute struct{}

func (s swaggerRoute) V1Route(config *Config, router gin.IRouter) {
  docs.SwaggerInfo.BasePath = "api/v1"
  router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
