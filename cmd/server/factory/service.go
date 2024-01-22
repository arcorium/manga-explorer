package factory

import (
	"manga-explorer/internal/app/common"
	service3 "manga-explorer/internal/app/service"
	service2 "manga-explorer/internal/domain/mangas/service"
	"manga-explorer/internal/domain/users/service"
	service5 "manga-explorer/internal/infrastructure/file/service"
	service4 "manga-explorer/internal/infrastructure/mail/service"
)

type Service struct {
	// Utility
	Mail service4.IMail
	File service5.IFile

	User           service.IUser
	Authentication service.IAuthentication
	Verification   service.IVerification
	Manga          service2.IManga
	Chapter        service2.IChapter
	Genre          service2.IGenre
}

func CreateServices(config *common.Config, repository *Repository) Service {
	result := Service{
		Mail:           service4.NewSMTPMailService(config),
		File:           service5.NewLocalFileService("./file"),
		Authentication: service3.NewCredential(config, repository.Credential, repository.User),
		Verification:   service3.NewVerification(repository.Verification),
		Chapter:        service3.NewChapterService(repository.Chapter, repository.Comment),
		Genre:          service3.NewGenreService(repository.Genre),
	}

	result.User = service3.NewUser(repository.User, result.Verification, result.Authentication, result.Mail)
	result.Manga = service3.NewMangaService(result.File, repository.Manga, repository.Comment, repository.Rate)
	return result
}
