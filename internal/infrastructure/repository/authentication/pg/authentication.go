package pg

import (
	"context"
	"manga-explorer/internal/domain/auth"
	"manga-explorer/internal/domain/auth/repository"
	"manga-explorer/internal/util"
	"time"

	"github.com/uptrace/bun"
)

func NewCredential(db bun.IDB) repository.IAuthentication {
	return &credentialRepository{db: db}
}

type credentialRepository struct {
	db bun.IDB
}

func (c credentialRepository) Create(cred *auth.Credential) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := c.db.NewInsert().
		Model(cred).
		Exec(ctx)
	return util.CheckSqlResult(res, err)
}

func (c credentialRepository) Find(userId, credId string) (*auth.Credential, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	cred := new(auth.Credential)
	err := c.db.NewSelect().
		Model(cred).
		Where("id = ? AND user_id = ?", credId, userId).
		Scan(ctx, cred)

	if err != nil {
		return nil, err
	}

	return cred, err
}

func (c credentialRepository) UpdateAccessTokenId(credentialId, accessTokenId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// TODO: Move this into service
	cred := auth.Credential{
		Id:            credentialId,
		AccessTokenId: accessTokenId,
		UpdatedAt:     time.Now(),
	}

	res, err := c.db.NewUpdate().
		Model(&cred).
		WherePK().
		OmitZero().
		Exec(ctx)
	return util.CheckSqlResult(res, err)
}

func (c credentialRepository) FindByAccessTokenId(accessTokenId string) (*auth.Credential, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	cred := new(auth.Credential)
	err := c.db.NewSelect().
		Model(cred).
		Where("access_token_id = ?", accessTokenId).
		Scan(ctx)

	if err != nil {
		return nil, err
	}
	return cred, err
}

func (c credentialRepository) FindUserCredentials(userId string) ([]auth.Credential, error) {
	var creds []auth.Credential

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := c.db.NewSelect().
		Model(&creds).
		Where("user_id = ?", userId).
		Scan(ctx)

	// err is nil even when there is no data returned
	return util.CheckSliceResult(creds, err).Unwrap()
}

func (c credentialRepository) Remove(userId, credId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := c.db.NewDelete().
		Model((*auth.Credential)(nil)).
		Where("id = ? AND user_id = ?", credId, userId).
		Exec(ctx)

	return util.CheckSqlResult(res, err)
}

func (c credentialRepository) RemoveByAccessTokenId(userId, accessTokenId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := c.db.NewDelete().
		Model((*auth.Credential)(nil)).
		Where("user_id = ? AND access_token_id = ?", userId, accessTokenId).
		Exec(ctx)

	return util.CheckSqlResult(res, err)
}

func (c credentialRepository) RemoveUserCredentials(userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := c.db.NewDelete().
		Model((*auth.Credential)(nil)).
		Where("user_id = ?", userId).
		Exec(ctx)

	return util.CheckSqlResult(res, err)
}
