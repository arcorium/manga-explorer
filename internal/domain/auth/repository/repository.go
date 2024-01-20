package repository

import "manga-explorer/internal/domain/auth"

type IAuthentication interface {
	Create(credential *auth.Credential) error
	Find(userId, credId string) (*auth.Credential, error)
	FindByAccessTokenId(accessTokenId string) (*auth.Credential, error)
	UpdateAccessTokenId(credentialId, accessTokenId string) error
	FindUserCredentials(userId string) ([]auth.Credential, error)
	Remove(userId, credId string) error
	RemoveByAccessTokenId(userId, accessTokenId string) error
	RemoveUserCredentials(userId string) error
}
