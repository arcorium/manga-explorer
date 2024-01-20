package service

import (
	"manga-explorer/internal/app/common/status"
	"manga-explorer/internal/domain/users"
	"manga-explorer/internal/domain/users/dto"
	"manga-explorer/internal/domain/users/mapper"
	"manga-explorer/internal/domain/users/repository"
	"manga-explorer/internal/domain/users/service"
	"time"
)

func NewVerification(repo repository.IVerification) service.IVerification {
	return &verificationService{repo: repo}
}

type verificationService struct {
	repo repository.IVerification
}

func (v verificationService) Request(userId string, usage users.Usage) (dto.VerificationResponse, status.Object) {
	verif := users.NewVerification(userId, usage)
	err := v.repo.Create(&verif)
	stat := status.FromRepository(err, status.CREATED)
	if stat.IsError() {
		return dto.VerificationResponse{}, stat
	}

	return mapper.ToVerificationResponse(&verif), stat
}

func (v verificationService) Find(token string) (dto.VerificationResponse, status.Object) {
	verif, err := v.repo.Find(token)
	return mapper.ToVerificationResponse(&verif), status.FromRepository(err)
}

func (v verificationService) Remove(token string) status.Object {
	err := v.repo.Remove(token)
	return status.FromRepository(err)
}

func (v verificationService) Validate(response *dto.VerificationResponse) status.Object {
	// Check token expiration
	if response.ExpirationTime.Unix() > time.Now().Unix() {
		return status.Error(status.VERIFICATION_TOKEN_EXPIRED)
	}

	// Check token usage
	if response.Usage != users.UsageResetPassword.String() {
		return status.Error(status.VERIFICATION_TOKEN_MISUSE)
	}

	return status.Success()
}
