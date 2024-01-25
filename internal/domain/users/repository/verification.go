package repository

import (
	"manga-explorer/internal/domain/users"
)

type IVerification interface {
	Upsert(verification *users.Verification) error
	Find(token string) (users.Verification, error)
	Remove(token string) error
	RemoveByUserId(userId string, usage users.Usage) error
}
