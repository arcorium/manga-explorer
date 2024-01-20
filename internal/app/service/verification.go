package service

import (
	"manga-explorer/internal/app/common/status"
	"manga-explorer/internal/domain/users"
	"manga-explorer/internal/domain/users/dto"
	"manga-explorer/internal/domain/users/mapper"
	"manga-explorer/internal/domain/users/repository"
	"manga-explorer/internal/domain/users/service"
	"time"

	"manga-explorer/internal/app/common"
)

func NewVerification(repo repository.IVerification) service.IVerification {
	return &verificationService{repo: repo}
}

type verificationService struct {
	repo repository.IVerification
}

func (v verificationService) Request(userId string, usage users.Usage) (dto.VerificationResponse, common.Status) {
	verif := users.New(userId, usage)
	err := v.repo.Create(&verif)
	return mapper.ToVerificationResponse(&verif), common.ConditionalStatus(err, status.INTERNAL_SERVER_ERROR)
}

func (v verificationService) Find(token string) (dto.VerificationResponse, common.Status) {
	verif, err := v.repo.Find(token)
	return mapper.ToVerificationResponse(&verif), common.ConditionalStatus(err, status.INTERNAL_SERVER_ERROR)
}

func (v verificationService) Remove(token string) common.Status {
	err := v.repo.Remove(token)
	return common.ConditionalStatus(err, status.INTERNAL_SERVER_ERROR)
}

func (v verificationService) Validate(response *dto.VerificationResponse) common.Status {
	// Check token expiration
	if response.ExpirationTime.Unix() > time.Now().Unix() {
		return common.StatusError(status.VERIFICATION_TOKEN_EXPIRED)
	}

	// Check token usage
	if response.Usage != users.UsageResetPassword.String() {
		return common.StatusError(status.VERIFICATION_TOKEN_MISUSE)
	}

	return common.StatusSuccess()
}
