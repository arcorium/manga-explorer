package pg

import (
	"context"
	"database/sql"
	"manga-explorer/internal/domain/users"
	"manga-explorer/internal/domain/users/repository"
	"manga-explorer/internal/util"
	"time"

	"github.com/uptrace/bun"
)

func NewUser(db bun.IDB) repository.IUser {
	return &UserRepository{db: db}
}

type UserRepository struct {
	db bun.IDB
}

func (u UserRepository) GetAllUsers() ([]users.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var result []users.User

	err := u.db.NewSelect().
		Model(&result).
		Scan(ctx)

	return result, err
}

func (u UserRepository) CreateProfile(profile *users.Profile) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := u.db.NewInsert().
		Model(profile).
		Returning("NULL").
		Exec(ctx)
	return util.CheckSqlResult(res, err)
}

func (u UserRepository) UpdateProfile(profile *users.Profile) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := u.db.NewUpdate().
		Model(profile).
		OmitZero().
		Where("user_id = ?", profile.UserId).
		Exec(ctx)

	return err
}

func (u UserRepository) CreateUser(user *users.User, profile *users.Profile) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := u.db.BeginTx(ctx, &sql.TxOptions{})

	_, err = tx.NewInsert().
		Model(user).
		Returning("NULL").
		Exec(ctx)

	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.NewInsert().
		Model(profile).
		Returning("NULL").
		Exec(ctx)
	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return err
}

func (u UserRepository) FindUserById(id string) (*users.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	usr := new(users.User)
	err := u.db.NewSelect().
		Model(usr).
		Where("id = ?", id).
		Scan(ctx)
	return usr, err
}

func (u UserRepository) FindUserProfiles(userId string) (*users.Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	profile := new(users.Profile)
	err := u.db.NewSelect().
		Model(profile).
		Where("user_id = ?", userId).
		Scan(ctx)
	return profile, err
}

func (u UserRepository) FindUserByEmail(email string) (*users.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	usr := new(users.User)
	err := u.db.NewSelect().
		Model(usr).
		Where("email = ?", email).
		Scan(ctx, usr)
	return usr, err
}

func (u UserRepository) UpdateUser(user *users.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := u.db.NewUpdate().
		Model(user).
		OmitZero().
		WherePK().
		Exec(ctx)

	return util.CheckSqlResult(res, err)
}

func (u UserRepository) DeleteUser(userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := u.db.NewDelete().
		Model((*users.User)(nil)).
		Where("id = ?", userId).
		Exec(ctx)

	return util.CheckSqlResult(res, err)
}