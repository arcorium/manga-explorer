package factory

import (
  "github.com/uptrace/bun"
  mangaRepo "manga-explorer/internal/domain/mangas/repository"
  userRepo "manga-explorer/internal/domain/users/repository"
  mangaPg "manga-explorer/internal/infrastructure/repository/mangas/pg"
  userPg "manga-explorer/internal/infrastructure/repository/users/pg"
)

type Repository struct {
  User         userRepo.IUser
  Credential   userRepo.IAuthentication
  Verification userRepo.IVerification
  Manga        mangaRepo.IManga
  Chapter      mangaRepo.IChapter
  Comment      mangaRepo.IComment
  Genre        mangaRepo.IGenre
  Rate         mangaRepo.IRate
  Translation  mangaRepo.ITranslation
}

func CreateRepositories(db bun.IDB) Repository {
  return Repository{
    User:         userPg.NewUser(db),
    Credential:   userPg.NewCredential(db),
    Verification: userPg.NewVerification(db),
    Manga:        mangaPg.NewManga(db),
    Chapter:      mangaPg.NewMangaChapter(db),
    Comment:      mangaPg.NewComment(db),
    Genre:        mangaPg.NewMangaGenre(db),
    Rate:         mangaPg.NewMangaRate(db),
    Translation:  mangaPg.NewTranslationRepository(db),
  }
}
