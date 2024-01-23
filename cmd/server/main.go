package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"manga-explorer/cmd/server/factory"
	"manga-explorer/database"
	"manga-explorer/internal/api/http/route"
	"manga-explorer/internal/app/common"
	"os"
	"os/signal"
	"syscall"
)

// @title           Manga Explorer API
// @version         1.0
// @description     Simple manga explorer used to read manga and save the manga

// @contact.name    arcorium
// @contact.url     github.com/arcorium
// @contact.email   arcorium.l@gmail.com

// @BasePath        /api/v1
func main() {
	config, err := common.LoadConfig("test")
	if err != nil {
		log.Fatalln(err)
	}

	db, err := database.Open(config, true)
	if err != nil {
		log.Fatalln("Failed to open database connection: ", err)
	}
	defer database.Close(db)

	repositories := factory.CreateRepositories(db)
	database.AddAdminUser(repositories.User)
	services := factory.CreateServices(config, &repositories)

	engine := gin.Default()
	router := factory.CreateRouter(config, &services, engine)

	authRoute := route.NewAuthRoute()
	userRoute := route.NewUserRoute()
	mangaRoute := route.NewMangaRoute()
	router.Routes(authRoute, userRoute, mangaRoute)

	run(config, engine)
}

func run(config *common.Config, engine *gin.Engine) {
	go func() {
		if err := engine.Run(config.Endpoint()); err != nil {
			log.Fatalln(err)
		}
	}()

	quitChan := make(chan os.Signal)
	signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quitChan
	close(quitChan)
}
