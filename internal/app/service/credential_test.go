package service

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"manga-explorer/internal/app/common"
	"manga-explorer/internal/app/common/status"
	"manga-explorer/internal/domain/auth"
	"manga-explorer/internal/domain/auth/dto"
	"manga-explorer/internal/domain/auth/mapper"
	"manga-explorer/internal/domain/auth/repository"
	mockAuthRepo "manga-explorer/internal/domain/auth/repository/mocks"
	"manga-explorer/internal/domain/users"
	"manga-explorer/internal/util"
	"manga-explorer/internal/util/containers"
	"time"

	userRepo "manga-explorer/internal/domain/users/repository"
	mockUserRepo "manga-explorer/internal/domain/users/repository/mocks"
	"reflect"
	"testing"
)

var conf *common.Config

var simpleError = errors.New("")

func config() *common.Config {
	if conf == nil {
		tmp, err := common.LoadConfig("test", "../../../")
		if err != nil {
			panic(err) // Should not handle
		}
		conf = &tmp
	}
	return conf
}

func newCredentialService(authRepos repository.IAuthentication, userRepos userRepo.IUser) credentialService {
	return credentialService{
		config:   config(),
		authRepo: authRepos,
		userRepo: userRepos,
	}
}

func Test_credentialService_Authenticate(t *testing.T) {

	// User repo mock
	userMock := mockUserRepo.NewUserMock(t)
	temp, err := users.NewUser("arcorium", "arcorium.l@gmail.com", "arcorium", users.RoleAdmin)
	userMock.EXPECT().FindUserByEmail("arcorium.l@gmail.com").Return(&temp, err)
	userMock.EXPECT().FindUserByEmail(mock.AnythingOfType("string")).Return(nil, simpleError)
	//userMock.On("FindUserByEmail", mock.AnythingOfType("string")).Return(nil, simpleError)

	// Auth repo mock
	//errCredential := auth.NewCredential(&temp, "Test", uuid.NewString(), util.GenerateRandomString(40))
	authMock := mockAuthRepo.NewAuthenticationMock(t)
	//authMock.EXPECT().Create(&errCredential).Return(simpleError)
	authMock.EXPECT().Create(mock.AnythingOfType("*auth.Credential")).Return(nil)

	type args struct {
		request dto.LoginInput
	}
	tests := []struct {
		name  string
		args  args
		want1 common.Status
	}{
		{
			name: "Normal",
			args: args{
				request: dto.LoginInput{
					Email:    "arcorium.l@gmail.com",
					Password: "arcorium",
				},
			},
			want1: common.StatusSuccess(),
		},
		{
			name: "Bad email",
			args: args{
				request: dto.LoginInput{
					Email:    util.GenerateRandomString(20),
					Password: "arcorium",
				},
			},
			want1: common.StatusError(status.BAD_BODY_REQUEST_ERROR),
		},
		{
			name: "User not found",
			args: args{
				request: dto.LoginInput{
					Email:    util.GenerateRandomString(10) + "@gmail.com",
					Password: "arcorium",
				},
			},
			want1: common.StatusError(status.USER_NOT_FOUND),
		},
		{
			name: "Wrong password",
			args: args{
				request: dto.LoginInput{
					Email:    "arcorium.l@gmail.com",
					Password: util.GenerateRandomString(10),
				},
			},
			want1: common.StatusError(status.USER_LOGIN_ERROR),
		},
	}

	c := newCredentialService(authMock, userMock)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := c.Authenticate(&tt.args.request)
			if !tt.want1.IsError() && (len(got.AccessToken) == 0 || len(got.Type) == 0) {
				t.Errorf("Authenticate() got = %v ", got)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Authenticate() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_credentialService_GetCredentials(t *testing.T) {

	// User repo mock
	userMock := mockUserRepo.NewUserMock(t)

	// Auth repo mock
	userId := uuid.NewString()
	creds := []auth.Credential{
		{
			Id:            uuid.NewString(),
			UserId:        userId,
			AccessTokenId: uuid.NewString(),
			Device: auth.Device{
				Name: "Test",
			},
			Token:     util.GenerateRandomString(20),
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
		},
		{
			Id:            uuid.NewString(),
			UserId:        userId,
			AccessTokenId: uuid.NewString(),
			Device: auth.Device{
				Name: "Test",
			},
			Token:     util.GenerateRandomString(20),
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
		},
	}
	authMock := mockAuthRepo.NewAuthenticationMock(t)
	authMock.EXPECT().FindUserCredentials(userId).Return(creds, nil)
	authMock.EXPECT().FindUserCredentials(mock.AnythingOfType("string")).Return(nil, sql.ErrNoRows)

	type args struct {
		userId string
	}
	tests := []struct {
		name  string
		args  args
		want  []dto.CredentialResponse
		want1 common.Status
	}{
		{
			name: "Normal",
			args: args{
				userId: userId,
			},
			want:  containers.CastSlicePtr(creds, mapper.ToCredentialResponse),
			want1: common.StatusSuccess(),
		},
		{
			name: "Empty user id",
			args: args{
				userId: "",
			},
			want:  nil,
			want1: common.StatusError(status.OBJECT_NOT_FOUND),
		},
		{
			name: "User doesn't found or doesn't have credentials",
			args: args{
				userId: uuid.NewString(),
			},
			want:  nil,
			want1: common.StatusError(status.OBJECT_NOT_FOUND),
		},
	}
	c := newCredentialService(authMock, userMock)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, got1 := c.GetCredentials(tt.args.userId)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCredentials() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetCredentials() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_credentialService_Logout(t *testing.T) {
	userId := uuid.NewString()
	credId := uuid.NewString()

	authMock := mockAuthRepo.NewAuthenticationMock(t)
	authMock.EXPECT().Remove(userId, credId).Return(nil)
	authMock.EXPECT().Remove(userId, mock.AnythingOfType("string")).Return(sql.ErrNoRows)
	authMock.EXPECT().Remove(mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(sql.ErrNoRows)

	userMock := mockUserRepo.NewUserMock(t)

	type args struct {
		userId string
		credId string
	}
	tests := []struct {
		name string
		args args
		want common.Status
	}{
		{
			name: "Normal",
			args: args{
				userId: userId,
				credId: credId,
			},
			want: common.StatusSuccess(),
		},
		{
			name: "User doesn't exists",
			args: args{
				userId: uuid.NewString(),
				credId: credId,
			},
			want: common.StatusError(status.OBJECT_NOT_FOUND),
		},
		{
			name: "Credential doesn't exists",
			args: args{
				userId: userId,
				credId: uuid.NewString(),
			},
			want: common.StatusError(status.OBJECT_NOT_FOUND),
		},
		{
			name: "Both user and credential doesn't exists",
			args: args{
				userId: uuid.NewString(),
				credId: uuid.NewString(),
			},
			want: common.StatusError(status.OBJECT_NOT_FOUND),
		},
	}
	c := newCredentialService(authMock, userMock)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.Logout(tt.args.userId, tt.args.credId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Logout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_credentialService_LogoutDevices(t *testing.T) {
	userId := uuid.NewString()
	authMock := mockAuthRepo.NewAuthenticationMock(t)
	authMock.EXPECT().RemoveUserCredentials(userId).Return(nil)
	authMock.EXPECT().RemoveUserCredentials(mock.AnythingOfType("string")).Return(sql.ErrNoRows)

	userMock := mockUserRepo.NewUserMock(t)
	type args struct {
		userId string
	}
	tests := []struct {
		name string
		args args
		want common.Status
	}{
		{
			name: "Normal",
			args: args{
				userId: userId,
			},
			want: common.StatusSuccess(),
		},
		{
			name: "User doesn't exists",
			args: args{
				userId: uuid.NewString(),
			},
			want: common.StatusError(status.OBJECT_NOT_FOUND),
		},
	}
	c := newCredentialService(authMock, userMock)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.LogoutDevices(tt.args.userId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LogoutDevices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_credentialService_RefreshToken(t *testing.T) {
	userId := uuid.NewString()
	authMock := mockAuthRepo.NewAuthenticationMock(t)
	authMock.EXPECT().RemoveUserCredentials(userId).Return(nil)
	authMock.EXPECT().RemoveUserCredentials(mock.AnythingOfType("string")).Return(sql.ErrNoRows)

	userMock := mockUserRepo.NewUserMock(t)
	type args struct {
		request dto.RefreshTokenInput
	}
	tests := []struct {
		name  string
		args  args
		want  dto.RefreshTokenResponse
		want1 common.Status
	}{}
	c := newCredentialService(authMock, userMock)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := c.RefreshToken(&tt.args.request)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RefreshToken() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("RefreshToken() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_credentialService_SelfLogout(t *testing.T) {
	type fields struct {
		config   *common.Config
		authRepo repository.IAuthentication
		userRepo userRepo.IUser
	}
	type args struct {
		userId        string
		accessTokenId string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   common.Status
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := credentialService{
				config:   tt.fields.config,
				authRepo: tt.fields.authRepo,
				userRepo: tt.fields.userRepo,
			}
			if got := c.SelfLogout(tt.args.userId, tt.args.accessTokenId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SelfLogout() = %v, want %v", got, tt.want)
			}
		})
	}
}
