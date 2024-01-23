package route

import (
	"github.com/gin-gonic/gin"
)

func NewMangaRoute() IRoute {
	return &_mangaRoute{}
}

type _mangaRoute struct {
}

func (m _mangaRoute) V1Route(config *Config, router gin.IRouter) {
	m.MangaRoute(config, router)
	m.ChapterRoute(config, router)
	m.GenreRoute(config, router)
}

func (m _mangaRoute) MangaRoute(config *Config, router gin.IRouter) {
	mangaRoute := router.Group("/mangas")
	// Manga IRoute
	mangaController := &config.Controller.Manga
	chapterController := &config.Controller.MangaChapter

	mangaRoute.GET("/", mangaController.ListManga)
	mangaRoute.GET("/search", mangaController.Search)
	mangaRoute.GET("/random", mangaController.Random)
	mangaRoute.GET("/:manga_id", mangaController.FindMangaById)
	mangaRoute.GET("/:manga_id/comments", mangaController.FindMangaComments)
	mangaRoute.GET("/:manga_id/ratings", mangaController.FindMangaRatings)
	mangaRoute.GET("/:manga_id/translates/*language", mangaController.FindMangaTranslations)
	// Login user
	mangaRoute.Use(config.Middleware.Authorization.Handle)
	mangaRoute.POST("/:manga_id/comments", mangaController.CreateMangaComments)
	mangaRoute.POST("/:manga_id/ratings", mangaController.CreateMangaRatings)
	mangaRoute.GET("/favorite", mangaController.GetMangaFavorites)
	mangaRoute.GET("/history", mangaController.GetMangaHistories)
	// Admin
	mangaRoute.Use(config.Middleware.AdminRestrict.Handle)

	mangaRoute.POST("/:manga_id/translates", mangaController.InsertMangaTranslate)
	mangaRoute.DELETE("/:manga_id/translates", mangaController.DeleteMangaTranslations)
	mangaRoute.PUT("/translates", mangaController.UpdateTranslation)
	mangaRoute.DELETE("/translates", mangaController.DeleteTranslations)

	mangaRoute.POST("/", mangaController.CreateManga)
	mangaRoute.POST("/:manga_id", mangaController.EditManga)
	mangaRoute.POST("/:manga_id/volumes", mangaController.CreateVolume)
	mangaRoute.DELETE("/:manga_id/volumes/:volume", mangaController.DeleteVolume)
	mangaRoute.POST("/:manga_id/chapters", chapterController.CreateChapter)

	mangaRoute.PUT("/:manga_id/cover", mangaController.UpdateMangaCover)
}
func (m _mangaRoute) ChapterRoute(config *Config, router gin.IRouter) {
	chapterController := &config.Controller.MangaChapter
	chapterRoute := router.Group("/chapters")
	chapterRoute.GET("/:chapter_id/comments", chapterController.FindChapterComments)
	chapterRoute.GET("/:chapter_id", chapterController.FindChapterPages)
	// Login user
	chapterRoute.Use(config.Middleware.Authorization.Handle)
	chapterRoute.POST("/:chapter_id/comments", chapterController.CreateChapterComments)
	// Admin
	chapterRoute.Use(config.Middleware.AdminRestrict.Handle)
	chapterRoute.PUT("/:chapter_id", chapterController.EditChapter)
	chapterRoute.DELETE("/:chapter_id", chapterController.DeleteChapter)

	chapterRoute.POST("/:chapter_id/pages", chapterController.InsertChapterPage)
	chapterRoute.DELETE("/:chapter_id/pages", chapterController.DeleteChapterPage)

	// Page IRoute
	pageRoute := router.Group("/pages")
	pageRoute.GET("/:page_id/comments", chapterController.FindPageComments)
	// Login user
	pageRoute.Use(config.Middleware.Authorization.Handle)
	pageRoute.POST("/:page_id/comments", chapterController.CreatePageComments)

	volumeRoute := router.Group("/volumes")
	volumeRoute.GET("/:volume_id", chapterController.FindVolumeChapters)
}
func (m _mangaRoute) GenreRoute(config *Config, router gin.IRouter) {
	genreController := &config.Controller.MangaGenre
	genreRoute := router.Group("/genres", config.Middleware.Authorization.Handle, config.Middleware.AdminRestrict.Handle)
	genreRoute.GET("/", genreController.ListGenre)
	genreRoute.POST("/", genreController.CreateGenre)
	genreRoute.DELETE("/:genre_id", genreController.DeleteGenre)
}
