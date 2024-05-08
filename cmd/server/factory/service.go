package factory

import (
	"github.com/gin-gonic/gin"
	"manga-explorer/internal/app/service"
	"manga-explorer/internal/common"
	"manga-explorer/internal/common/constant"
	mangaService "manga-explorer/internal/domain/mangas/service"
	userService "manga-explorer/internal/domain/users/service"
	fileService "manga-explorer/internal/infrastructure/file/service"
	mailService "manga-explorer/internal/infrastructure/mail/service"
)

type Service struct {
	// Utility
	Mail mailService.IMail
	File fileService.IFile

	User           userService.IUser
	Authentication userService.IAuthentication
	Verification   userService.IVerification
	Manga          mangaService.IManga
	Chapter        mangaService.IChapter
	Genre          mangaService.IGenre
}

func CreateServices(config *common.Config, repository *Repository, router gin.IRouter) Service {
	result := Service{
		Mail: mailService.NewSMTPMailService(constant.SenderEmail, mailService.SMTPMailerConfig{
			Host: config.SMTPHost,
			Port: config.SMTPPort,
			User: config.SMTPUser,
			Pass: config.SMTPPass,
		}),
		File:           fileService.NewLocalFileService(config, config.Endpoint(), "/static", "./files", router), // Used for both user profile and manga chapter images
		Authentication: service.NewCredential(config, repository.Credential, repository.User),
		Verification:   service.NewVerification(config, repository.Verification),
		Genre:          service.NewGenreService(repository.Genre),
	}

	result.User = service.NewUser(config, repository.User, result.Verification, result.Authentication, result.Mail, result.File)
	result.Manga = service.NewMangaService(result.File, repository.Manga, repository.Translation, repository.Comment, repository.Rate)
	result.Chapter = service.NewChapterService(result.File, repository.Chapter, repository.Comment)

	return result
}
