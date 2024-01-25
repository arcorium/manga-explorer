package repository

import (
	"manga-explorer/internal/domain/users"
)

type IUser interface {
	CreateUser(user *users.User, profile *users.Profile) error
	CreateProfile(profile *users.Profile) error
	GetAllUsers() ([]users.User, error)
	FindUserById(userId string) (*users.User, error)
	FindUserByEmail(email string) (*users.User, error)
	FindUserProfiles(userId string) (*users.Profile, error)
	UpdateUser(user *users.User) error
	UpdateProfileByUserId(profile *users.Profile) error
	UpdateProfile(profile *users.Profile) error
	DeleteUser(userId string) error
}
