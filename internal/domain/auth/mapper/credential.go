package mapper

import (
	"manga-explorer/internal/domain/auth"
	"manga-explorer/internal/domain/auth/dto"
)

func ToCredentialResponse(credential *auth.Credential) dto.CredentialResponse {
	return dto.CredentialResponse{
		Id:         credential.Id,
		DeviceName: credential.Device.Name,
	}
}
