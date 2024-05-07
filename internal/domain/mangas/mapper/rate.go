package mapper

import (
  "github.com/google/uuid"
  "manga-explorer/internal/domain/mangas"
  "manga-explorer/internal/domain/mangas/dto"
  "manga-explorer/internal/domain/users/mapper"
  "time"
)

func ToRatingResponse(rate *mangas.Rate) dto.RateResponse {
  return dto.RateResponse{
    User: mapper.ToUserResponse(rate.User),
    Rate: rate.Rate,
    Time: rate.CreatedAt,
  }
}

func MapRateUpsertInput(input *dto.RateUpsertInput) mangas.Rate {
  now := time.Now()
  return mangas.Rate{
    Id:        uuid.NewString(),
    UserId:    input.UserId,
    MangaId:   input.MangaId,
    Rate:      input.Rate,
    CreatedAt: now,
    UpdatedAt: now,
  }
}
