package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"manga-explorer/database"
	mangaController "manga-explorer/internal/api/http/controller/v1/mangas"
	userController "manga-explorer/internal/api/http/controller/v1/users"
	"manga-explorer/internal/api/http/middleware"
	"manga-explorer/internal/api/http/route"
	"manga-explorer/internal/app/common"
	"manga-explorer/internal/app/common/constant"
	"manga-explorer/internal/app/service"
	"manga-explorer/internal/domain/users"
	service3 "manga-explorer/internal/infrastructure/file/service"
	service2 "manga-explorer/internal/infrastructure/mail/service"
	mangaRepo "manga-explorer/internal/infrastructure/repository/mangas/pg"
	"manga-explorer/internal/infrastructure/repository/users/pg"
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

	db := database.Open(&config, true)
	if db == nil {
		log.Fatalln("Failed to open database connection")
	}
	defer database.Close(db)
	database.RegisterModels(db)

	engine := gin.Default()

	userRepos := pg.NewUser(db)
	credRepos := pg.NewCredential(db)
	verifRepo := pg.NewVerification(db)
	mangaRepos := mangaRepo.NewManga(db)
	commentRepo := mangaRepo.NewComment(db)
	chapterRepo := mangaRepo.NewMangaChapter(db)
	genreRepo := mangaRepo.NewMangaGenre(db)
	rateRepo := mangaRepo.NewMangaRate(db)

	database.AddAdminUser(userRepos)

	// Utility service
	mailService := service2.NewSMTPMailService(&config)
	fileService := service3.NewLocalFileService("./files")

	mangaService := service.NewMangaService(fileService, mangaRepos, commentRepo, rateRepo)
	authService := service.NewCredential(&config, credRepos, userRepos)
	verifService := service.NewVerification(verifRepo)
	userService := service.NewUser(userRepos, verifService, authService, mailService)
	chapterService := service.NewChapterService(chapterRepo, commentRepo)
	genreService := service.NewGenreService(genreRepo)

	controllerConfig := route.ConfigController{
		Auth:         userController.NewAuthController(authService),
		User:         userController.NewUserController(userService),
		Manga:        mangaController.NewMangaController(mangaService),
		MangaChapter: mangaController.NewChapterController(chapterService),
		MangaGenre:   mangaController.NewGenreController(genreService),
	}

	middlewareConfig := route.ConfigMiddleware{
		Auth: middleware.NewAuthMiddleware(config.JWTSecretKey, &middleware.AuthMiddlewareConfig{
			SigningMethod: config.SigningMethod(),
			ClaimsKey:     constant.ClaimsKey,
		}),
		AdminRestrict: middleware.NewRoleRestrictionMiddleware(&middleware.RoleRestrictionConfig{
			ClaimsKey: constant.ClaimsKey,
		}, []users.Role{users.RoleAdmin}),
	}
	routerConfig := route.Config{
		Controller: controllerConfig,
		Middleware: middlewareConfig,
	}

	rt := route.NewRouter(&routerConfig, &engine.RouterGroup)
	authRoute := route.NewAuthRoute()
	userRoute := route.NewUserRoute()
	mangaRoute := route.NewMangaRoute()
	rt.Routes(authRoute, userRoute, mangaRoute)

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
