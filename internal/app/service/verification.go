package service

import (
	"manga-explorer/internal/common"
	"manga-explorer/internal/common/status"
	"manga-explorer/internal/domain/users"
	"manga-explorer/internal/domain/users/dto"
	"manga-explorer/internal/domain/users/mapper"
	"manga-explorer/internal/domain/users/repository"
	"manga-explorer/internal/domain/users/service"
	"manga-explorer/internal/util/opt"
	"time"
)

func NewVerification(config *common.Config, repo repository.IVerification) service.IVerification {
	return &verificationService{config: config, repo: repo}
}

type verificationService struct {
	config *common.Config
	repo   repository.IVerification
}

func (v verificationService) Request(userId string, usage users.Usage) (dto.VerificationResponse, status.Object) {
	verif := users.NewVerification(userId, usage, v.config.VerificationTokenDuration)
	err := v.repo.Upsert(&verif)
	if err != nil {
		return dto.VerificationResponse{}, status.RepositoryError(err, opt.New(status.VERIFICATION_REQUEST_FAILED))
	}

	return mapper.ToVerificationResponse(&verif), status.Success()
}

func (v verificationService) Response(token string, usage users.Usage) (dto.VerificationResponse, status.Object) {
	// Get the token
	verif, err := v.repo.Find(token)
	if err != nil {
		return dto.VerificationResponse{}, status.RepositoryError(err, opt.New(status.VERIFICATION_TOKEN_NOT_FOUND))
	}

	// Check token expiration
	if verif.ExpirationTime.Before(time.Now()) {
		// Delete token
		v.repo.Remove(token)
		return dto.VerificationResponse{}, status.Error(status.VERIFICATION_TOKEN_EXPIRED)
	}

	// Check token usage
	if verif.Usage != usage {
		// Delete token
		v.repo.Remove(token)
		return dto.VerificationResponse{}, status.Error(status.VERIFICATION_TOKEN_MISUSE)
	}

	// Remove token
	err = v.repo.Remove(verif.Token)
	return mapper.ToVerificationResponse(&verif), status.ConditionalRepository(err, status.SUCCESS, opt.New(status.VERIFICATION_TOKEN_NOT_FOUND))
}
