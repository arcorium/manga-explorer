package repository

import (
	"manga-explorer/internal/domain/users"
)

type IAuthentication interface {
	Create(credential *users.Credential) error
	Find(userId, credId string) (*users.Credential, error)
	FindByAccessTokenId(accessTokenId string) (*users.Credential, error)
	UpdateAccessTokenId(credentialId, accessTokenId string) error
	FindUserCredentials(userId string) ([]users.Credential, error)
	Remove(userId, credId string) error
	RemoveByAccessTokenId(userId, accessTokenId string) error
	RemoveUserCredentials(userId string) error
}
