package service

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"manga-explorer/internal/common"
	"manga-explorer/internal/common/constant"
	"manga-explorer/internal/common/status"
	"manga-explorer/internal/domain/users"
	"manga-explorer/internal/domain/users/dto"
	"manga-explorer/internal/domain/users/mapper"
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
		conf = tmp
	}
	return conf
}

func newCredentialServiceForTest(authRepos userRepo.IAuthentication, userRepos userRepo.IUser) credentialService {
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
	userMock.EXPECT().FindUserByEmail(mock.AnythingOfType("string")).Return(nil, sql.ErrNoRows)

	// Auth repo mock
	//errCredential := auth.NewCredential(&temp, "Test", uuid.NewString(), util.GenerateRandomString(40))
	authMock := mockUserRepo.NewAuthenticationMock(t)
	//authMock.EXPECT().Upsert(&errCredential).Return(simpleError)
	authMock.EXPECT().Create(mock.AnythingOfType("*users.Credential")).Return(nil)

	type args struct {
		request dto.LoginInput
	}
	tests := []struct {
		name  string
		args  args
		want1 status.Object
	}{
		{
			name: "Normal",
			args: args{
				request: dto.LoginInput{
					Email:    "arcorium.l@gmail.com",
					Password: "arcorium",
				},
			},
			want1: status.Success(),
		},
		{
			name: "User not found",
			args: args{
				request: dto.LoginInput{
					Email:    util.GenerateRandomString(10) + "@gmail.com",
					Password: "arcorium",
				},
			},
			want1: status.NotFoundError(),
		},
		{
			name: "Wrong password",
			args: args{
				request: dto.LoginInput{
					Email:    "arcorium.l@gmail.com",
					Password: util.GenerateRandomString(10),
				},
			},
			want1: status.Error(status.USER_LOGIN_ERROR),
		},
	}

	c := newCredentialServiceForTest(authMock, userMock)
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
	creds := []users.Credential{
		{
			Id:            uuid.NewString(),
			UserId:        userId,
			AccessTokenId: uuid.NewString(),
			Device: users.Device{
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
			Device: users.Device{
				Name: "Test",
			},
			Token:     util.GenerateRandomString(20),
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
		},
	}
	authMock := mockUserRepo.NewAuthenticationMock(t)
	authMock.EXPECT().FindUserCredentials(userId).Return(creds, nil)
	authMock.EXPECT().FindUserCredentials(mock.AnythingOfType("string")).Return(nil, sql.ErrNoRows)

	type args struct {
		userId string
	}
	tests := []struct {
		name  string
		args  args
		want  []dto.CredentialResponse
		want1 status.Object
	}{
		{
			name: "Normal",
			args: args{
				userId: userId,
			},
			want:  containers.CastSlicePtr(creds, mapper.ToCredentialResponse),
			want1: status.Success(),
		},
		{
			name: "Empty user id",
			args: args{
				userId: "",
			},
			want:  nil,
			want1: status.Error(status.CREDENTIALS_NOT_FOUND),
		},
		{
			name: "User doesn't found or doesn't have credentials",
			args: args{
				userId: uuid.NewString(),
			},
			want:  nil,
			want1: status.Error(status.CREDENTIALS_NOT_FOUND),
		},
	}
	c := newCredentialServiceForTest(authMock, userMock)
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

	authMock := mockUserRepo.NewAuthenticationMock(t)
	authMock.EXPECT().Remove(userId, credId).Return(nil)
	authMock.EXPECT().Remove(mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(sql.ErrNoRows)

	userMock := mockUserRepo.NewUserMock(t)

	type args struct {
		userId string
		credId string
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
				credId: credId,
			},
			want: status.Deleted(),
		},
		{
			name: "User doesn't exists",
			args: args{
				userId: uuid.NewString(),
				credId: credId,
			},
			want: status.Error(status.LOGOUT_CREDENTIAL_NOT_FOUND),
		},
		{
			name: "Credential doesn't exists",
			args: args{
				userId: userId,
				credId: uuid.NewString(),
			},
			want: status.Error(status.LOGOUT_CREDENTIAL_NOT_FOUND),
		},
		{
			name: "Both user and credential doesn't exists",
			args: args{
				userId: uuid.NewString(),
				credId: uuid.NewString(),
			},
			want: status.Error(status.LOGOUT_CREDENTIAL_NOT_FOUND),
		},
	}
	c := newCredentialServiceForTest(authMock, userMock)
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
	authMock := mockUserRepo.NewAuthenticationMock(t)
	authMock.EXPECT().RemoveUserCredentials(userId).Return(nil)
	authMock.EXPECT().RemoveUserCredentials(mock.AnythingOfType("string")).Return(sql.ErrNoRows)

	userMock := mockUserRepo.NewUserMock(t)
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
			want: status.Success(),
		},
		{
			name: "User doesn't exists",
			args: args{
				userId: uuid.NewString(),
			},
			want: status.Error(status.LOGOUT_CREDENTIAL_NOT_FOUND),
		},
	}
	c := newCredentialServiceForTest(authMock, userMock)
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
	accessTokenId := uuid.NewString()
	cred := users.Credential{
		Id:            uuid.NewString(),
		UserId:        userId,
		AccessTokenId: uuid.NewString(),
		Device: users.Device{
			Name: "Test",
		},
		Token:     util.GenerateRandomString(20),
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}

	badAccessTokenId := uuid.NewString()
	badCred := users.Credential{
		Id:            uuid.NewString(),
		UserId:        uuid.NewString(),
		AccessTokenId: uuid.NewString(),
		Device: users.Device{
			Name: "Test",
		},
		Token:     util.GenerateRandomString(20),
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}

	user, err := users.NewUser("arcorium", "arcorium@gmail.com", "arcorium", users.RoleAdmin)
	require.NoError(t, err)

	authMock := mockUserRepo.NewAuthenticationMock(t)
	authMock.EXPECT().FindByAccessTokenId(accessTokenId).Return(&cred, nil)
	authMock.EXPECT().FindByAccessTokenId(mock.AnythingOfType("string")).Return(nil, sql.ErrNoRows)
	authMock.EXPECT().UpdateAccessTokenId(cred.Id, uuid.NewString()).Return(nil)

	authMock.EXPECT().FindByAccessTokenId(badAccessTokenId).Return(&badCred, nil)

	userMock := mockUserRepo.NewUserMock(t)
	userMock.EXPECT().FindUserById(userId).Return(&user, nil)

	type args struct {
		request dto.RefreshTokenInput
	}
	tests := []struct {
		name  string
		args  args
		want1 status.Object
	}{
		{
			name: "Normal",
			args: args{
				request: dto.RefreshTokenInput{
					Type:        constant.TokenType,
					AccessToken: accessTokenId,
				},
			},
			want1: status.Updated(),
		},
		{
			name: "Refresh token without owner",
			args: args{
				request: dto.RefreshTokenInput{
					Type:        constant.TokenType,
					AccessToken: badAccessTokenId,
				},
			},
			want1: status.Updated(),
		},
		{
			name: "Access token not found",
			args: args{
				request: dto.RefreshTokenInput{
					Type:        constant.TokenType,
					AccessToken: uuid.NewString(),
				},
			},
			want1: status.Error(status.ACCESS_TOKEN_WITHOUT_REFRESH_TOKEN),
		},
	}
	c := newCredentialServiceForTest(authMock, userMock)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got1 := c.RefreshToken(&tt.args.request)
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("RefreshToken() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_credentialService_SelfLogout(t *testing.T) {
	userId := uuid.NewString()
	accessTokenId := uuid.NewString()

	authMock := mockUserRepo.NewAuthenticationMock(t)
	authMock.EXPECT().RemoveByAccessTokenId(userId, accessTokenId).Return(nil)
	authMock.EXPECT().RemoveByAccessTokenId(mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(sql.ErrNoRows)

	userMock := mockUserRepo.NewUserMock(t)

	c := newCredentialServiceForTest(authMock, userMock)

	type args struct {
		userId        string
		accessTokenId string
	}
	tests := []struct {
		name string
		args args
		want status.Object
	}{
		{
			name: "Normal",
			args: args{
				userId:        userId,
				accessTokenId: accessTokenId,
			},
			want: status.Deleted(),
		},
		{
			name: "User has no credential",
			args: args{
				userId:        userId,
				accessTokenId: accessTokenId,
			},
			want: status.Deleted(),
		},
		{
			name: "User and access token doesn't exists",
			args: args{
				userId:        uuid.NewString(),
				accessTokenId: uuid.NewString(),
			},
			want: status.Deleted(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.SelfLogout(tt.args.userId, tt.args.accessTokenId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SelfLogout() = %v, want %v", got, tt.want)
			}
		})
	}
}
