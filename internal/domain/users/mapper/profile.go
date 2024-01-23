package mapper

import (
	"manga-explorer/internal/domain/users"
	"manga-explorer/internal/domain/users/dto"
	"manga-explorer/internal/infrastructure/file"
	fileService "manga-explorer/internal/infrastructure/file/service"
	"time"
)

func toProfileResponse(profile *users.Profile, fl fileService.IFile) dto.InternalProfileResponse {
	return dto.InternalProfileResponse{
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		PhotoURL:  fl.GetFullpath(file.ProfileAsset, profile.PhotoURL),
		Bio:       profile.Bio,
	}
}

func ToProfileResponse(profile *users.Profile, file fileService.IFile) dto.ProfileResponse {
	return dto.ProfileResponse{
		UserResponse:     ToUserResponse(profile.User), // User should be same
		ProfileResponses: toProfileResponse(profile, file),
	}
}

func MapProfileRegisterInput(user *users.User, input *dto.UserRegisterInput) users.Profile {
	return users.Profile{
		UserId:    user.Id,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}
}

func MapAddProfileInput(user *users.User, input *dto.AddUserInput) users.Profile {
	return users.Profile{
		UserId:    user.Id,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		PhotoURL:  file.Name(input.PhotoURL),
		Bio:       input.Bio,
		UpdatedAt: time.Now(),
	}
}

func MapUpdateProfileExtendedInput(input *dto.UpdateProfileExtendedInput) users.Profile {
	return users.Profile{
		UserId:    input.UserId,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		PhotoURL:  file.Name(input.PhotoURL),
		Bio:       input.Bio,
		UpdatedAt: time.Now(),
	}
}

func MapProfileUpdateInput(input *dto.ProfileUpdateInput) users.Profile {
	return users.Profile{
		UserId:    input.UserId,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Bio:       input.Bio,
		UpdatedAt: time.Now(),
	}
}
