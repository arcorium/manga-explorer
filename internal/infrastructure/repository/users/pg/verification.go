package pg

import (
  "context"
  "manga-explorer/internal/domain/users"
  "manga-explorer/internal/domain/users/repository"
  "manga-explorer/internal/util"
  "time"

  "github.com/uptrace/bun"
)

func NewVerification(db bun.IDB) repository.IVerification {
  return &verificationRepository{db: db}
}

type verificationRepository struct {
  db bun.IDB
}

func (v verificationRepository) Upsert(verif *users.Verification) error {
  ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
  defer cancel()

  query := v.db.NewInsert().
    Model(verif).
    Returning("NULL").
    On("CONFLICT (user_id, usage) DO UPDATE").
    Set("token = EXCLUDED.token, expiration_time= EXCLUDED.expiration_time")

  res, err := query.
    Exec(ctx)
  return util.CheckSqlResult(res, err)
}

func (v verificationRepository) Find(token string) (users.Verification, error) {
  ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
  defer cancel()

  verif := users.Verification{}
  err := v.db.NewSelect().
    Model(&verif).
    Where("token = ?", token).
    Scan(ctx)
  return verif, err
}

func (v verificationRepository) Remove(token string) error {
  ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
  defer cancel()

  res, err := v.db.NewDelete().
    Model((*users.Verification)(nil)).
    Where("token = ?", token).
    Exec(ctx)
  return util.CheckSqlResult(res, err)
}

func (v verificationRepository) RemoveByUserId(userId string, usage users.Usage) error {
  ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
  defer cancel()

  res, err := v.db.NewDelete().
    Model((*users.Verification)(nil)).
    Where("user_id = ? AND usage = ?", userId, usage).
    Exec(ctx)
  return util.CheckSqlResult(res, err)
}
