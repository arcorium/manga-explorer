package factory

import (
  "github.com/gin-gonic/gin"
  service3 "manga-explorer/internal/app/service"
  "manga-explorer/internal/common"
  "manga-explorer/internal/common/constant"
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

func CreateServices(config *common.Config, repository *Repository, router gin.IRouter) Service {
  result := Service{
    Mail: service4.NewSMTPMailService(constant.SenderEmail, service4.SMTPMailerConfig{
      Host: config.SMTPHost,
      Port: config.SMTPPort,
      User: config.SMTPUser,
      Pass: config.SMTPPass,
    }),
    File:           service5.NewLocalFileService(config, config.Endpoint(), "/static", "./files", router), // Used for both user profile and manga chapter images
    Authentication: service3.NewCredential(config, repository.Credential, repository.User),
    Verification:   service3.NewVerification(config, repository.Verification),
    Genre:          service3.NewGenreService(repository.Genre),
  }

  result.User = service3.NewUser(config, repository.User, result.Verification, result.Authentication, result.Mail, result.File)
  result.Manga = service3.NewMangaService(result.File, repository.Manga, repository.Translation, repository.Comment, repository.Rate)
  result.Chapter = service3.NewChapterService(result.File, repository.Chapter, repository.Comment)

  return result
}
