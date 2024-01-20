package repository

import (
	"manga-explorer/internal/domain/users"
)

type IVerification interface {
	Create(verification *users.Verification) error
	Find(token string) (users.Verification, error)
	Remove(token string) error
}
