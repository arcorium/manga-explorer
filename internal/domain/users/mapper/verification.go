package mapper

import (
	"manga-explorer/internal/domain/users"
	"manga-explorer/internal/domain/users/dto"
)

func ToVerificationResponse(verif *users.Verification) dto.VerificationResponse {
	return dto.VerificationResponse{
		UserId:         verif.UserId,
		Token:          verif.Token,
		Usage:          verif.Usage.String(),
		ExpirationTime: verif.ExpirationTime,
	}
}
