package mapper

import (
	"manga-explorer/internal/domain/users"
	"manga-explorer/internal/domain/users/dto"
	"manga-explorer/internal/infrastructure/file"
	"time"
)

func ToProfileResponse(profile *users.Profile) dto.ProfileResponse {
	return dto.ProfileResponse{
		UserResponse: ToUserResponse(profile.User),
		FirstName:    profile.FirstName,
		LastName:     profile.LastName,
		PhotoURL:     profile.PhotoURL.HostnameFullpath(file.ProfileAsset),
		Bio:          profile.Bio,
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
