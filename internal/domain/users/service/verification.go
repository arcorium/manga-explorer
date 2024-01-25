package service

import (
	"manga-explorer/internal/common/status"
	"manga-explorer/internal/domain/users"
	"manga-explorer/internal/domain/users/dto"
)

type IVerification interface {
	// Request Upsert verification token for the user and send the url which contains the token into user email
	// It is returning the response, message, and error
	Request(userId string, usage users.Usage) (dto.VerificationResponse, status.Object)
	// Response the verification token by validating and removing the verification, returned verification response should be already removed from data source
	Response(token string, usage users.Usage) (dto.VerificationResponse, status.Object)
}
