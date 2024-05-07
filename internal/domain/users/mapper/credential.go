package mapper

import (
  "manga-explorer/internal/domain/users"
  "manga-explorer/internal/domain/users/dto"
)

func ToCredentialResponse(credential *users.Credential) dto.CredentialResponse {
  return dto.CredentialResponse{
    Id:         credential.Id,
    DeviceName: credential.Device.Name,
  }
}
