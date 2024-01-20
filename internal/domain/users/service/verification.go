package service

import (
	"manga-explorer/internal/app/common/status"
	"manga-explorer/internal/domain/users"
	"manga-explorer/internal/domain/users/dto"
)

type IVerification interface {
	// Request Create verification token for the user and send the url which contains the token into user email
	// It is returning the response, message, and error
	Request(userId string, usage users.Usage) (dto.VerificationResponse, status.Object)
	// Find Find the token
	Find(token string) (dto.VerificationResponse, status.Object)
	// Remove Removing token
	Remove(token string) status.Object
	Validate(response *dto.VerificationResponse) status.Object
}
