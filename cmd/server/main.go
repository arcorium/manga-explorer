package main

import (
  "errors"
  "github.com/caarlos0/env/v10"
  "github.com/gin-gonic/gin"
  "github.com/gin-gonic/gin/binding"
  "github.com/go-playground/validator/v10"
  "log"
  "manga-explorer/cmd/server/factory"
  "manga-explorer/database"
  "manga-explorer/internal/api/http/route"
  "manga-explorer/internal/common"
  "os"
  "os/signal"
  "syscall"
)

// @title			Manga Explorer API
// @version		1.0
// @description	Simple manga reader app with rate, favorite and comment features
// @contact.name	arcorium
// @contact.url	github.com/arcorium
// @contact.email	arcorium.l@gmail.com
// @BasePath		/api/v1
func main() {
  config, err := common.LoadConfig()
  if err != nil {
    var aggregateError env.AggregateError
    ok := errors.As(err, &aggregateError)
    if ok {
      for _, v := range aggregateError.Errors {
        log.Println(v)
      }
      os.Exit(-1)
    } else {
      log.Fatalln(err)
    }
  }

  db, err := database.Open(config, false)
  if err != nil {
    log.Fatalln("Failed to open database connection: ", err)
  }
  defer database.Close(db)

  repositories := factory.CreateRepositories(db)

  engine := gin.Default()
  if len(config.TrustedProxies) > 0 {
    engine.ForwardedByClientIP = true
    err := engine.SetTrustedProxies(config.TrustedProxies)
    if err != nil {
      log.Fatalln(err)
    }
  }

  common.RegisterValidationTags(binding.Validator.Engine().(*validator.Validate))
  services := factory.CreateServices(config, &repositories, engine)

  router := factory.CreateRouter(config, &services, engine)

  authRoute := route.NewAuthRoute()
  userRoute := route.NewUserRoute()
  mangaRoute := route.NewMangaRoute()
  swaggerRoute := route.NewSwaggerRoute()
  router.Routes(authRoute, userRoute, mangaRoute, swaggerRoute)

  run(config, engine)
}

func run(config *common.Config, engine *gin.Engine) {
  log.Println("Server started!")

  go func() {
    if err := engine.Run(config.Endpoint()); err != nil {
      log.Fatalln(err)
    }
  }()

  quitChan := make(chan os.Signal)
  signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
  <-quitChan
  close(quitChan)

  log.Println("Server stopped!")
}
