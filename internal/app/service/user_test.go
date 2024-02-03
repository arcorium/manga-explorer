package service

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"manga-explorer/internal/common/status"
	"manga-explorer/internal/domain/users"
	"manga-explorer/internal/domain/users/dto"
	"manga-explorer/internal/domain/users/mapper"
	"manga-explorer/internal/domain/users/repository"
	userRepoMock "manga-explorer/internal/domain/users/repository/mocks"
	"manga-explorer/internal/domain/users/service"
	userServiceMock "manga-explorer/internal/domain/users/service/mocks"
	"manga-explorer/internal/infrastructure/file"
	fileService "manga-explorer/internal/infrastructure/file/service"
	fileServiceMock "manga-explorer/internal/infrastructure/file/service/mocks"
	mailService "manga-explorer/internal/infrastructure/mail/service"
	mailServiceMock "manga-explorer/internal/infrastructure/mail/service/mocks"
	"manga-explorer/internal/util"
	"manga-explorer/internal/util/opt"
	"reflect"
	"testing"
)

func newUserServiceMocked(user repository.IUser, verification service.IVerification, authentication service.IAuthentication, mail mailService.IMail, file fileService.IFile) userService {
	return userService{
		repo:         user,
		verifService: verification,
		authService:  authentication,
		mailService:  mail,
		fileService:  file,
	}
}

func Test_userService_AddUser(t *testing.T) {
	input := &dto.AddUserInput{
		Username:  util.GenerateRandomString(10),
		Email:     util.GenerateRandomString(10) + "@gmail.com",
		Password:  util.GenerateRandomString(10),
		Verified:  false,
		Role:      users.RoleUser.String(),
		FirstName: util.GenerateRandomString(10),
		LastName:  util.GenerateRandomString(10),
		Bio:       util.GenerateRandomString(10),
		PhotoURL:  util.GenerateRandomString(10),
	}

	userInput, err := mapper.MapAddUserInput(input)
	require.NoError(t, err)
	profileInput := mapper.MapAddProfileInput(&userInput, input)

	mockedUserRepo := userRepoMock.NewUserMock(t)
	mockedUserRepo.On("CreateUser", mock.Anything, mock.Anything).Return(func(user *users.User, profile *users.Profile) error {
		// Reset generated id
		user.Id = userInput.Id
		user.Password = userInput.Password
		user.UpdatedAt = userInput.UpdatedAt
		user.CreatedAt = userInput.CreatedAt
		user.BannedUntil = userInput.BannedUntil
		user.DeletedAt = userInput.DeletedAt
		profile.UserId = user.Id
		profile.UpdatedAt = profileInput.UpdatedAt

		if reflect.DeepEqual(user, &userInput) && reflect.DeepEqual(profile, &profileInput) {
			return nil
		}
		return sql.ErrNoRows
	})

	mockedVerifService := userServiceMock.NewVerificationMock(t)
	mockedAuthService := userServiceMock.NewAuthenticationMock(t)
	mockedMailService := mailServiceMock.NewMailMock(t)
	mockedFileService := fileServiceMock.NewFileMock(t)

	u := newUserServiceMocked(mockedUserRepo, mockedVerifService, mockedAuthService, mockedMailService, mockedFileService)
	type args struct {
		input *dto.AddUserInput
	}
	tests := []struct {
		name string
		args args
		want status.Object
	}{
		{
			name: "Normal",
			args: args{
				input: input,
			},
			want: status.Created(),
		},
		{
			name: "Wrong Role",
			args: args{
				input: &dto.AddUserInput{
					Username:  "",
					Email:     "asdadsa@gmail.com",
					Password:  "asd",
					Verified:  false,
					Role:      "unknown",
					FirstName: "",
					LastName:  "",
					Bio:       "",
					PhotoURL:  "",
				},
			},
			want: status.BadRequestError(),
		},
		{
			name: "Duplicate",
			args: args{
				input: &dto.AddUserInput{
					Username:  "",
					Email:     "asdadsa@gmail.com",
					Password:  "asd",
					Verified:  false,
					Role:      users.RoleUser.String(),
					FirstName: "",
					LastName:  "",
					Bio:       "",
					PhotoURL:  "",
				},
			},
			want: status.Error(status.USER_CREATION_ALREADY_EXIST),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := u.AddUser(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_ChangePassword(t *testing.T) {
	user, err := users.NewUser(util.GenerateRandomString(10), "asdads@gmail.com", "something", users.RoleUser)
	require.NoError(t, err)
	mockedUserRepo := userRepoMock.NewUserMock(t)
	mockedUserRepo.EXPECT().FindUserById(user.Id).Return(&user, nil)
	mockedUserRepo.EXPECT().FindUserById(mock.Anything).Return(nil, sql.ErrNoRows)

	mockedUserRepo.EXPECT().UpdateUser(mock.Anything).Return(nil)

	mockedVerifService := userServiceMock.NewVerificationMock(t)
	mockedAuthService := userServiceMock.NewAuthenticationMock(t)
	mockedMailService := mailServiceMock.NewMailMock(t)
	mockedFileService := fileServiceMock.NewFileMock(t)

	u := newUserServiceMocked(mockedUserRepo, mockedVerifService, mockedAuthService, mockedMailService, mockedFileService)
	type args struct {
		input *dto.ChangePasswordInput
	}
	tests := []struct {
		name string
		args args
		want status.Object
	}{
		{
			name: "Normal",
			args: args{
				input: &dto.ChangePasswordInput{
					UserId:       user.Id,
					LastPassword: "something",
					NewPassword:  "123",
				},
			},
			want: status.Updated(),
		},
		{
			name: "Different last password",
			args: args{
				input: &dto.ChangePasswordInput{
					UserId:       user.Id,
					LastPassword: util.GenerateRandomString(10),
					NewPassword:  util.GenerateRandomString(10),
				},
			},
			want: status.Error(status.USER_LOGIN_ERROR),
		},
		{
			name: "No Such user",
			args: args{
				input: &dto.ChangePasswordInput{
					UserId:       uuid.NewString(),
					LastPassword: "",
					NewPassword:  "",
				},
			},
			want: status.Error(status.USER_NOT_FOUND),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := u.ChangePassword(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChangePassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_DeleteProfileImage(t *testing.T) {
	userId := uuid.NewString()
	badUserId := uuid.NewString()

	profile := users.NewProfile(userId, util.GenerateRandomString(10), util.GenerateRandomString(10), util.GenerateRandomString(10), opt.New(file.Name("asadasd.jpg")))
	profile2 := users.NewProfile(badUserId, util.GenerateRandomString(10), util.GenerateRandomString(10), util.GenerateRandomString(10), file.NullName)

	mockedUserRepo := userRepoMock.NewUserMock(t)
	mockedUserRepo.EXPECT().FindUserProfiles(userId).Return(&profile, nil)
	mockedUserRepo.EXPECT().FindUserProfiles(badUserId).Return(&profile2, nil)
	mockedUserRepo.EXPECT().UpdateProfile(mock.Anything).Return(nil)

	mockedVerifService := userServiceMock.NewVerificationMock(t)
	mockedAuthService := userServiceMock.NewAuthenticationMock(t)
	mockedMailService := mailServiceMock.NewMailMock(t)
	mockedFileService := fileServiceMock.NewFileMock(t)
	mockedFileService.EXPECT().Delete(mock.Anything, mock.Anything).Return(status.Deleted()).Once()

	u := newUserServiceMocked(mockedUserRepo, mockedVerifService, mockedAuthService, mockedMailService, mockedFileService)
	type args struct {
		userId string
	}
	tests := []struct {
		name string
		args args
		want status.Object
	}{
		{
			name: "Normal",
			args: args{
				userId: userId,
			},
			want: status.Deleted(),
		},
		{
			name: "Profile has no photo",
			args: args{
				userId: badUserId,
			},
			want: status.Deleted(),
		},
		{
			name: "No user",
			args: args{
				userId: uuid.NewString(),
			},
			want: status.Error(status.PROFILE_UPDATE_FAILED),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := u.DeleteProfileImage(tt.args.userId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteProfileImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_DeleteUser(t *testing.T) {
	userId := uuid.NewString()
	badUserId := uuid.NewString()

	mockedUserRepo := userRepoMock.NewUserMock(t)
	mockedUserRepo.EXPECT().DeleteUser(userId).Return(nil)
	mockedUserRepo.EXPECT().DeleteUser(badUserId).Return(sql.ErrNoRows)
	mockedUserRepo.EXPECT().DeleteUser(mock.Anything).Return(simpleError)

	mockedVerifService := userServiceMock.NewVerificationMock(t)
	mockedAuthService := userServiceMock.NewAuthenticationMock(t)
	mockedMailService := mailServiceMock.NewMailMock(t)
	mockedFileService := fileServiceMock.NewFileMock(t)

	u := newUserServiceMocked(mockedUserRepo, mockedVerifService, mockedAuthService, mockedMailService, mockedFileService)
	type args struct {
		userId string
	}
	tests := []struct {
		name string
		args args
		want status.Object
	}{
		{
			name: "Normal",
			args: args{
				userId: userId,
			},
			want: status.Deleted(),
		},
		{
			name: "User not found",
			args: args{
				userId: badUserId,
			},
			want: status.Error(status.USER_NOT_FOUND),
		},
		{
			name: "Another error",
			args: args{
				userId: util.GenerateRandomString(10),
			},
			want: status.InternalError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := u.DeleteUser(tt.args.userId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_FindUserByEmail(t *testing.T) {
	user, err := users.NewUser(util.GenerateRandomString(10), "asdad@gmail.com", "something", users.RoleUser)
	require.NoError(t, err)

	badEmail := util.GenerateRandomString(10)
	mockedUserRepo := userRepoMock.NewUserMock(t)
	mockedUserRepo.EXPECT().FindUserByEmail(user.Email).Return(&user, nil)
	mockedUserRepo.EXPECT().FindUserByEmail(badEmail).Return(nil, sql.ErrNoRows)
	mockedUserRepo.EXPECT().FindUserByEmail(mock.Anything).Return(nil, simpleError)

	mockedVerifService := userServiceMock.NewVerificationMock(t)
	mockedAuthService := userServiceMock.NewAuthenticationMock(t)
	mockedMailService := mailServiceMock.NewMailMock(t)
	mockedFileService := fileServiceMock.NewFileMock(t)
	u := newUserServiceMocked(mockedUserRepo, mockedVerifService, mockedAuthService, mockedMailService, mockedFileService)

	type args struct {
		email string
	}
	tests := []struct {
		name  string
		args  args
		want  dto.UserResponse
		want1 status.Object
	}{
		{
			name: "Normal",
			args: args{
				email: user.Email,
			},
			want:  mapper.ToUserResponse(&user),
			want1: status.Success(),
		},
		{
			name: "User not found",
			args: args{
				email: badEmail,
			},
			want:  mapper.ToUserResponse(&users.User{}),
			want1: status.Error(status.USER_NOT_FOUND),
		},
		{
			name: "Another error",
			args: args{
				email: util.GenerateRandomString(10),
			},
			want:  mapper.ToUserResponse(&users.User{}),
			want1: status.InternalError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := u.FindUserByEmail(tt.args.email)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindUserByEmail() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("FindUserByEmail() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_userService_FindUserById(t *testing.T) {
	user, err := users.NewUser(util.GenerateRandomString(10), "asdad@gmail.com", "something", users.RoleUser)
	require.NoError(t, err)

	badId := uuid.NewString()
	mockedUserRepo := userRepoMock.NewUserMock(t)
	mockedUserRepo.EXPECT().FindUserById(user.Id).Return(&user, nil)
	mockedUserRepo.EXPECT().FindUserById(badId).Return(nil, sql.ErrNoRows)
	mockedUserRepo.EXPECT().FindUserById(mock.Anything).Return(nil, simpleError)

	mockedVerifService := userServiceMock.NewVerificationMock(t)
	mockedAuthService := userServiceMock.NewAuthenticationMock(t)
	mockedMailService := mailServiceMock.NewMailMock(t)
	mockedFileService := fileServiceMock.NewFileMock(t)

	u := newUserServiceMocked(mockedUserRepo, mockedVerifService, mockedAuthService, mockedMailService, mockedFileService)

	type args struct {
		id string
	}
	tests := []struct {
		name  string
		args  args
		want  dto.UserResponse
		want1 status.Object
	}{
		{
			name: "Normal",
			args: args{
				id: user.Id,
			},
			want:  mapper.ToUserResponse(&user),
			want1: status.Success(),
		},
		{
			name: "User not found",
			args: args{
				id: badId,
			},
			want:  mapper.ToUserResponse(&users.User{}),
			want1: status.Error(status.USER_NOT_FOUND),
		},
		{
			name: "Another error",
			args: args{
				id: util.GenerateRandomString(10),
			},
			want:  mapper.ToUserResponse(&users.User{}),
			want1: status.InternalError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := u.FindUserById(tt.args.id)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindUserById() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("FindUserById() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_userService_FindUserProfileById(t *testing.T) {
	profile := users.NewProfile(uuid.NewString(), util.GenerateRandomString(10), util.GenerateRandomString(10), util.GenerateRandomString(10), file.NullName)
	badId := uuid.NewString()

	mockedUserRepo := userRepoMock.NewUserMock(t)
	mockedUserRepo.EXPECT().FindUserProfiles(profile.UserId).Return(&profile, nil)
	mockedUserRepo.EXPECT().FindUserProfiles(badId).Return(nil, sql.ErrNoRows)
	mockedUserRepo.EXPECT().FindUserProfiles(mock.Anything).Return(nil, simpleError)

	mockedVerifService := userServiceMock.NewVerificationMock(t)
	mockedAuthService := userServiceMock.NewAuthenticationMock(t)
	mockedMailService := mailServiceMock.NewMailMock(t)
	mockedFileService := fileServiceMock.NewFileMock(t)

	mockedFileService.EXPECT().GetFullpath(mock.Anything, mock.Anything).Return(util.GenerateRandomString(10))

	u := newUserServiceMocked(mockedUserRepo, mockedVerifService, mockedAuthService, mockedMailService, mockedFileService)

	type args struct {
		userId string
	}
	tests := []struct {
		name  string
		args  args
		want  dto.ProfileResponse
		want1 status.Object
	}{
		{
			name: "Normal",
			args: args{
				userId: profile.UserId,
			},
			want:  mapper.ToProfileResponse(&profile, mockedFileService),
			want1: status.Success(),
		},
		{
			name: "User not found / User has no profile",
			args: args{
				userId: badId,
			},
			want:  mapper.ToProfileResponse(&users.Profile{}, mockedFileService),
			want1: status.Error(status.PROFILE_NOT_FOUND),
		},
		{
			name: "Another error",
			args: args{
				userId: util.GenerateRandomString(10),
			},
			want:  mapper.ToProfileResponse(&users.Profile{}, mockedFileService),
			want1: status.InternalError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := u.FindUserProfileById(tt.args.userId)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindUserProfileById() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("FindUserProfileById() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_userService_GetAllUsers(t *testing.T) {
	mockedUserRepo := userRepoMock.NewUserMock(t)
	mockedUserRepo.EXPECT().GetAllUsers().Return(nil, sql.ErrNoRows)

	mockedVerifService := userServiceMock.NewVerificationMock(t)
	mockedAuthService := userServiceMock.NewAuthenticationMock(t)
	mockedMailService := mailServiceMock.NewMailMock(t)
	mockedFileService := fileServiceMock.NewFileMock(t)

	mockedFileService.EXPECT().GetFullpath(mock.Anything, mock.Anything).Return(util.GenerateRandomString(10))

	u := newUserServiceMocked(mockedUserRepo, mockedVerifService, mockedAuthService, mockedMailService, mockedFileService)
	tests := []struct {
		name  string
		want  []dto.UserResponse
		want1 status.Object
	}{
		{
			name:  "Normal",
			want:  nil,
			want1: status.Success(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := u.GetAllUsers()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllUsers() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetAllUsers() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_userService_RegisterUser(t *testing.T) {
	input := &dto.UserRegisterInput{
		Username:  util.GenerateRandomString(10),
		Email:     util.GenerateRandomString(10) + "@gmail.com",
		Password:  util.GenerateRandomString(10),
		FirstName: util.GenerateRandomString(10),
		LastName:  util.GenerateRandomString(10),
		Bio:       util.GenerateRandomString(30),
	}

	user, err := users.NewUser(input.Username, input.Email, input.Password, users.RoleUser)
	require.NoError(t, err)
	profile := users.NewProfile(user.Id, input.FirstName, input.LastName, input.Bio, file.NullName)

	mockedUserRepo := userRepoMock.NewUserMock(t)
	mockedUserRepo.On("CreateUser", mock.Anything, mock.Anything).Return(func(user1 *users.User, profile1 *users.Profile) error {
		user1.Id = user.Id
		user1.Password = user.Password
		user1.BannedUntil = user.BannedUntil
		user1.CreatedAt = user.CreatedAt
		user1.DeletedAt = user.DeletedAt
		user1.UpdatedAt = user.UpdatedAt

		profile1.Id = profile.Id
		profile1.UserId = profile.UserId
		profile1.UpdatedAt = profile.UpdatedAt

		if reflect.DeepEqual(&user, user1) && reflect.DeepEqual(&profile, profile1) {
			return nil
		}
		return sql.ErrNoRows
	})

	mockedVerifService := userServiceMock.NewVerificationMock(t)
	mockedVerifService.EXPECT().Request(mock.Anything, mock.Anything).Return(dto.VerificationResponse{}, status.Success())

	mockedAuthService := userServiceMock.NewAuthenticationMock(t)
	mockedMailService := mailServiceMock.NewMailMock(t)

	mockedMailService.EXPECT().SendEmail(mock.Anything).Return(nil)

	mockedFileService := fileServiceMock.NewFileMock(t)

	u := newUserServiceMocked(mockedUserRepo, mockedVerifService, mockedAuthService, mockedMailService, mockedFileService)

	type args struct {
		input *dto.UserRegisterInput
	}
	tests := []struct {
		name  string
		args  args
		want  dto.UserResponse
		want1 status.Object
	}{
		{
			name: "Normal",
			args: args{
				input: input,
			},
			want:  mapper.ToUserResponse(&user),
			want1: status.Created(),
		},
		{
			name: "User already exist",
			args: args{
				input: &dto.UserRegisterInput{
					Username: "",
					Email:    "asdad@gmail.com",
					Password: "asd",
				},
			},
			want:  dto.UserResponse{},
			want1: status.Error(status.USER_CREATION_ALREADY_EXIST),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := u.RegisterUser(tt.args.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegisterUser() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("RegisterUser() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_userService_RequestEmailVerification(t *testing.T) {
	user, err := users.NewUser(util.GenerateRandomString(10), "asdad@gmail.com", "something", users.RoleUser)
	require.NoError(t, err)
	badUserId := uuid.NewString()

	mockedUserRepo := userRepoMock.NewUserMock(t)
	mockedUserRepo.EXPECT().FindUserById(user.Id).Return(&user, nil)
	mockedUserRepo.EXPECT().FindUserById(badUserId).Return(nil, sql.ErrNoRows)
	mockedUserRepo.EXPECT().FindUserById(mock.Anything).Return(nil, simpleError)

	mockedVerifService := userServiceMock.NewVerificationMock(t)
	mockedVerifService.EXPECT().Request(user.Id, mock.Anything).Return(dto.VerificationResponse{}, status.Success())

	mockedAuthService := userServiceMock.NewAuthenticationMock(t)
	mockedMailService := mailServiceMock.NewMailMock(t)
	mockedMailService.EXPECT().SendEmail(mock.Anything).Return(nil)

	mockedFileService := fileServiceMock.NewFileMock(t)

	u := newUserServiceMocked(mockedUserRepo, mockedVerifService, mockedAuthService, mockedMailService, mockedFileService)
	type args struct {
		input *dto.VerifEmailRequestInput
	}
	tests := []struct {
		name string
		args args
		want status.Object
	}{
		{
			name: "Normal",
			args: args{
				input: &dto.VerifEmailRequestInput{
					UserId: user.Id,
				},
			},
			want: status.Success(),
		},
		{
			name: "User not found",
			args: args{
				input: &dto.VerifEmailRequestInput{
					UserId: badUserId,
				},
			},
			want: status.Error(status.VERIFICATION_USER_NOT_EXISTS),
		},
		{
			name: "Another Error",
			args: args{
				input: &dto.VerifEmailRequestInput{
					UserId: util.GenerateRandomString(10),
				},
			},
			want: status.InternalError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := u.RequestEmailVerification(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RequestEmailVerification() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_RequestResetPassword(t *testing.T) {
	user, err := users.NewUser(util.GenerateRandomString(10), "asdad@gmail.com", "something", users.RoleUser)
	require.NoError(t, err)
	badEmail := "asdadasdas@gmail.com"

	mockedUserRepo := userRepoMock.NewUserMock(t)
	mockedUserRepo.EXPECT().FindUserByEmail(user.Email).Return(&user, nil)
	mockedUserRepo.EXPECT().FindUserByEmail(badEmail).Return(nil, sql.ErrNoRows)
	mockedUserRepo.EXPECT().FindUserByEmail(mock.Anything).Return(nil, simpleError)

	mockedVerifService := userServiceMock.NewVerificationMock(t)
	mockedVerifService.EXPECT().Request(user.Id, mock.Anything).Return(dto.VerificationResponse{}, status.Success())

	mockedAuthService := userServiceMock.NewAuthenticationMock(t)
	mockedMailService := mailServiceMock.NewMailMock(t)
	mockedMailService.EXPECT().SendEmail(mock.Anything).Return(nil)

	mockedFileService := fileServiceMock.NewFileMock(t)

	u := newUserServiceMocked(mockedUserRepo, mockedVerifService, mockedAuthService, mockedMailService, mockedFileService)
	type args struct {
		input *dto.ResetPasswordRequestInput
	}
	tests := []struct {
		name string
		args args
		want status.Object
	}{
		{
			name: "Normal",
			args: args{
				input: &dto.ResetPasswordRequestInput{Email: user.Email},
			},
			want: status.Success(),
		},
		{
			name: "User not found",
			args: args{
				input: &dto.ResetPasswordRequestInput{Email: badEmail},
			},
			want: status.Error(status.USER_NOT_FOUND),
		},
		{
			name: "Another error",
			args: args{
				input: &dto.ResetPasswordRequestInput{Email: "macmamc@mgial.com"},
			},
			want: status.InternalError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := u.RequestResetPassword(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RequestResetPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_ResetPassword(t *testing.T) {
	userId := uuid.NewString()
	token := util.GenerateRandomString(20)
	badToken := util.GenerateRandomString(20)

	mockedUserRepo := userRepoMock.NewUserMock(t)
	mockedUserRepo.EXPECT().UpdateUser(mock.Anything).Return(nil)

	mockedVerifService := userServiceMock.NewVerificationMock(t)
	mockedVerifService.EXPECT().Response(token, mock.Anything).Return(dto.VerificationResponse{UserId: userId}, status.Success())
	//mockedVerifService.EXPECT().Response(token, users.UsageResetPassword).Return(dto.VerificationResponse{}, status.Success())
	//mockedVerifService.EXPECT().Response(token, users.UsageVerifyEmail).Return(dto.VerificationResponse{}, status.Error(status.VERIFICATION_TOKEN_MISUSE))
	mockedVerifService.EXPECT().Response(badToken, mock.Anything).Return(dto.VerificationResponse{}, status.Error(status.VERIFICATION_TOKEN_NOT_FOUND))

	mockedAuthService := userServiceMock.NewAuthenticationMock(t)
	mockedAuthService.EXPECT().LogoutDevices(userId).Return(status.Success())

	mockedMailService := mailServiceMock.NewMailMock(t)
	mockedFileService := fileServiceMock.NewFileMock(t)

	u := newUserServiceMocked(mockedUserRepo, mockedVerifService, mockedAuthService, mockedMailService, mockedFileService)
	type args struct {
		input *dto.ResetPasswordInput
	}
	tests := []struct {
		name string
		args args
		want status.Object
	}{
		{
			name: "Normal",
			args: args{
				input: &dto.ResetPasswordInput{
					Token:       token,
					NewPassword: util.GenerateRandomString(20),
				},
			},
			want: status.Success(),
		},
		{
			name: "Token not found",
			args: args{
				input: &dto.ResetPasswordInput{
					Token:       badToken,
					NewPassword: util.GenerateRandomString(20),
				},
			},
			want: status.Error(status.VERIFICATION_TOKEN_NOT_FOUND),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := u.ResetPassword(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResetPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_UpdateProfile(t *testing.T) {
	input := &dto.ProfileEditInput{
		UserId:    uuid.NewString(),
		FirstName: util.GenerateRandomString(10),
		LastName:  util.GenerateRandomString(10),
		Bio:       util.GenerateRandomString(20),
	}
	profile := mapper.MapProfileUpdateInput(input)

	mockedUserRepo := userRepoMock.NewUserMock(t)
	mockedUserRepo.EXPECT().UpdateProfile(&profile).Return(nil)
	mockedUserRepo.EXPECT().UpdateProfile(mock.Anything).Return(sql.ErrNoRows)

	mockedVerifService := userServiceMock.NewVerificationMock(t)
	mockedAuthService := userServiceMock.NewAuthenticationMock(t)
	mockedMailService := mailServiceMock.NewMailMock(t)
	mockedFileService := fileServiceMock.NewFileMock(t)

	u := newUserServiceMocked(mockedUserRepo, mockedVerifService, mockedAuthService, mockedMailService, mockedFileService)
	type args struct {
		input *dto.ProfileEditInput
	}
	tests := []struct {
		name string
		args args
		want status.Object
	}{
		{
			name: "Normal",
			args: args{
				input: input,
			},
			want: status.Updated(),
		},
		{
			name: "Failed",
			args: args{
				input: &dto.ProfileEditInput{},
			},
			want: status.Error(status.PROFILE_UPDATE_FAILED),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := u.UpdateProfile(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateProfileByUserId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_UpdateProfileExtended(t *testing.T) {
	input := &dto.ProfileEditExtendedInput{
		UserId:    uuid.NewString(),
		FirstName: util.GenerateRandomString(10),
		LastName:  util.GenerateRandomString(10),
		Bio:       util.GenerateRandomString(20),
		PhotoURL:  util.GenerateRandomString(10),
	}
	profile := mapper.MapProfileUpdateExtendedInput(input)

	mockedUserRepo := userRepoMock.NewUserMock(t)
	mockedUserRepo.EXPECT().UpdateProfile(&profile).Return(nil)
	mockedUserRepo.EXPECT().UpdateProfile(mock.Anything).Return(sql.ErrNoRows)

	mockedVerifService := userServiceMock.NewVerificationMock(t)
	mockedAuthService := userServiceMock.NewAuthenticationMock(t)
	mockedMailService := mailServiceMock.NewMailMock(t)
	mockedFileService := fileServiceMock.NewFileMock(t)

	u := newUserServiceMocked(mockedUserRepo, mockedVerifService, mockedAuthService, mockedMailService, mockedFileService)
	type args struct {
		input *dto.ProfileEditExtendedInput
	}
	tests := []struct {
		name string
		args args
		want status.Object
	}{
		{
			name: "Normal",
			args: args{
				input: input,
			},
			want: status.Updated(),
		},
		{
			name: "Failed",
			args: args{
				input: &dto.ProfileEditExtendedInput{},
			},
			want: status.Error(status.PROFILE_UPDATE_FAILED),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := u.UpdateProfileExtended(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateProfileExtended() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_UpdateProfileImage(t *testing.T) {
	profile := users.NewProfile(uuid.NewString(), util.GenerateRandomString(10), util.GenerateRandomString(10), "", opt.New(file.Name("asdasdada")))
	badProfile := users.NewProfile(uuid.NewString(), util.GenerateRandomString(10), util.GenerateRandomString(10), "", file.NullName)

	mockedUserRepo := userRepoMock.NewUserMock(t)
	mockedUserRepo.EXPECT().FindUserProfiles(profile.UserId).Return(&profile, nil)
	mockedUserRepo.EXPECT().FindUserProfiles(badProfile.UserId).Return(&badProfile, nil)
	mockedUserRepo.EXPECT().FindUserProfiles(mock.Anything).Return(nil, sql.ErrNoRows)
	mockedUserRepo.EXPECT().UpdateProfile(mock.Anything).Return(nil)

	mockedVerifService := userServiceMock.NewVerificationMock(t)
	mockedAuthService := userServiceMock.NewAuthenticationMock(t)
	mockedMailService := mailServiceMock.NewMailMock(t)
	mockedFileService := fileServiceMock.NewFileMock(t)

	mockedFileService.EXPECT().Upload(file.ProfileAsset, mock.Anything).Return("something", status.Success())
	mockedFileService.EXPECT().Delete(file.ProfileAsset, mock.Anything).Return(status.Success()).Once()
	mockedFileService.EXPECT().Delete(mock.Anything, mock.Anything).Return(status.InternalError()).Once()

	u := newUserServiceMocked(mockedUserRepo, mockedVerifService, mockedAuthService, mockedMailService, mockedFileService)
	type args struct {
		input *dto.ProfileImageUpdateInput
	}
	tests := []struct {
		name string
		args args
		want status.Object
	}{
		{
			name: "Normal",
			args: args{
				input: &dto.ProfileImageUpdateInput{
					UserId: profile.UserId,
					Image:  nil,
				},
			},
			want: status.Updated(),
		},
		{
			name: "Profile with no image",
			args: args{
				input: &dto.ProfileImageUpdateInput{
					UserId: badProfile.UserId,
				},
			},
			want: status.Updated(),
		},
		{
			name: "User/Profile doesn't found",
			args: args{
				input: &dto.ProfileImageUpdateInput{
					UserId: uuid.NewString(),
				},
			},
			want: status.Error(status.PROFILE_NOT_FOUND),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := u.UpdateProfileImage(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateProfileImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_UpdateUser(t *testing.T) {
	input := &dto.UserEditInput{
		UserId:   uuid.NewString(),
		Username: util.GenerateRandomString(10),
		Email:    "asdadaw@gmail.com",
	}
	user := mapper.MapUserUpdateInput(input)

	mockedUserRepo := userRepoMock.NewUserMock(t)
	mockedUserRepo.EXPECT().UpdateUser(&user).Return(nil)
	mockedUserRepo.EXPECT().UpdateUser(mock.Anything).Return(sql.ErrNoRows)

	mockedVerifService := userServiceMock.NewVerificationMock(t)
	mockedAuthService := userServiceMock.NewAuthenticationMock(t)
	mockedMailService := mailServiceMock.NewMailMock(t)
	mockedFileService := fileServiceMock.NewFileMock(t)

	u := newUserServiceMocked(mockedUserRepo, mockedVerifService, mockedAuthService, mockedMailService, mockedFileService)

	type args struct {
		input *dto.UserEditInput
	}
	tests := []struct {
		name string
		args args
		want status.Object
	}{
		{
			name: "Normal",
			args: args{
				input: input,
			},
			want: status.Updated(),
		},
		{
			name: "Update failed",
			args: args{
				input: &dto.UserEditInput{},
			},
			want: status.Error(status.USER_UPDATE_FAILED),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := u.UpdateUser(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_UpdateUserExtended(t *testing.T) {
	input := &dto.UserEditExtendedInput{
		UserId:   uuid.NewString(),
		Username: util.GenerateRandomString(10),
		Email:    "asdadaw@gmail.com",
		Password: "pasdpasd",
	}
	user, err := mapper.MapUserUpdateExtendedInput(input)
	require.NoError(t, err)

	mockedUserRepo := userRepoMock.NewUserMock(t)
	mockedUserRepo.EXPECT().UpdateUser(&user).Return(nil)
	mockedUserRepo.EXPECT().UpdateUser(mock.Anything).Return(sql.ErrNoRows)

	mockedVerifService := userServiceMock.NewVerificationMock(t)
	mockedAuthService := userServiceMock.NewAuthenticationMock(t)
	mockedMailService := mailServiceMock.NewMailMock(t)
	mockedFileService := fileServiceMock.NewFileMock(t)

	u := newUserServiceMocked(mockedUserRepo, mockedVerifService, mockedAuthService, mockedMailService, mockedFileService)
	type args struct {
		input *dto.UserEditExtendedInput
	}
	tests := []struct {
		name string
		args args
		want status.Object
	}{
		{
			name: "Normal",
			args: args{
				input: input,
			},
			want: status.Updated(),
		},
		{
			name: "Update failed",
			args: args{
				input: &dto.UserEditExtendedInput{},
			},
			want: status.Error(status.PROFILE_UPDATE_FAILED),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := u.UpdateUserExtended(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EditUserExtended() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_VerifyEmail(t *testing.T) {
	userId := uuid.NewString()
	token := util.GenerateRandomString(20)
	badToken := util.GenerateRandomString(20)

	mockedUserRepo := userRepoMock.NewUserMock(t)
	mockedUserRepo.EXPECT().UpdateUser(mock.Anything).Return(nil)

	mockedVerifService := userServiceMock.NewVerificationMock(t)
	mockedVerifService.EXPECT().Response(token, mock.Anything).Return(dto.VerificationResponse{UserId: userId}, status.Success())
	//mockedVerifService.EXPECT().Response(token, users.UsageResetPassword).Return(dto.VerificationResponse{}, status.Success())
	//mockedVerifService.EXPECT().Response(token, users.UsageVerifyEmail).Return(dto.VerificationResponse{}, status.Error(status.VERIFICATION_TOKEN_MISUSE))
	mockedVerifService.EXPECT().Response(badToken, mock.Anything).Return(dto.VerificationResponse{}, status.Error(status.VERIFICATION_TOKEN_NOT_FOUND))

	mockedAuthService := userServiceMock.NewAuthenticationMock(t)

	mockedMailService := mailServiceMock.NewMailMock(t)
	mockedFileService := fileServiceMock.NewFileMock(t)

	u := newUserServiceMocked(mockedUserRepo, mockedVerifService, mockedAuthService, mockedMailService, mockedFileService)
	type args struct {
		input *dto.VerifyEmailInput
	}
	tests := []struct {
		name string
		args args
		want status.Object
	}{
		{
			name: "Normal",
			args: args{
				input: &dto.VerifyEmailInput{
					UserId: userId,
					Token:  token,
				},
			},
			want: status.Updated(),
		},
		{
			name: "Token not found",
			args: args{
				input: &dto.VerifyEmailInput{
					UserId: userId,
					Token:  badToken,
				},
			},
			want: status.Error(status.VERIFICATION_TOKEN_NOT_FOUND),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := u.VerifyEmail(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VerifyEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
