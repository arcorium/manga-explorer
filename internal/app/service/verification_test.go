package service

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"manga-explorer/internal/common/status"
	"manga-explorer/internal/domain/users"
	"manga-explorer/internal/domain/users/dto"
	"manga-explorer/internal/domain/users/mapper"
	verifRepoMock "manga-explorer/internal/domain/users/repository/mocks"
	"manga-explorer/internal/util"
	"reflect"
	"testing"
	"time"
)

func Test_verificationService_Request(t *testing.T) {
	userId := uuid.NewString()
	verif := users.Verification{
		Token:          util.GenerateRandomString(20),
		UserId:         userId,
		Usage:          users.UsageResetPassword,
		ExpirationTime: time.Now().Add(time.Hour * 1),
		User:           nil,
	}
	verifMock := verifRepoMock.NewVerificationMock(t)
	verifMock.EXPECT().Create(&verif).Return(nil)
	verifMock.EXPECT().Create(mock.AnythingOfType("*users.Verification")).Return(simpleError)

	v := NewVerification(verifMock)

	type args struct {
		userId string
		usage  users.Usage
	}
	tests := []struct {
		name  string
		args  args
		want  dto.VerificationResponse
		want1 status.Object
	}{
		{
			name: "Normal",
			args: args{
				userId: userId,
				usage:  users.UsageResetPassword,
			},
			want:  mapper.ToVerificationResponse(&verif),
			want1: status.Created(),
		},
		{
			name: "User doesn't have verification tokens",
			args: args{
				userId: uuid.NewString(),
				usage:  users.UsageVerifyEmail,
			},
			want:  dto.VerificationResponse{},
			want1: status.Error(status.VERIFICATION_REQUEST_FAILED),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := v.Request(tt.args.userId, tt.args.usage)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Request() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Request() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_verificationService_Response(t *testing.T) {
	userId := uuid.NewString()
	verif := users.Verification{
		Token:          util.GenerateRandomString(20),
		UserId:         userId,
		Usage:          users.UsageResetPassword,
		ExpirationTime: time.Now().Add(time.Hour * 1),
		User:           nil,
	}

	verif2 := users.Verification{
		Token:          util.GenerateRandomString(20),
		UserId:         userId,
		Usage:          users.UsageResetPassword,
		ExpirationTime: time.Now(),
		User:           nil,
	}

	verifMock := verifRepoMock.NewVerificationMock(t)
	verifMock.EXPECT().Find(verif.Token).Return(verif, nil)
	verifMock.EXPECT().Find(verif2.Token).Return(verif, nil)
	verifMock.EXPECT().Find(mock.AnythingOfType("string")).Return(users.Verification{}, sql.ErrNoRows)

	verifMock.EXPECT().Remove(verif.Token).Return(nil)
	verifMock.EXPECT().Remove(mock.AnythingOfType("string")).Return(sql.ErrNoRows)

	v := NewVerification(verifMock)
	type args struct {
		token string
		usage users.Usage
	}
	tests := []struct {
		name  string
		args  args
		want  dto.VerificationResponse
		want1 status.Object
	}{
		{
			name: "Normal",
			args: args{
				token: verif.Token,
				usage: verif.Usage,
			},
			want:  mapper.ToVerificationResponse(&verif),
			want1: status.Success(),
		},
		{
			name: "Different Usage",
			args: args{
				token: verif.Token,
				usage: users.UsageVerifyEmail,
			},
			want:  dto.VerificationResponse{},
			want1: status.Error(status.VERIFICATION_TOKEN_MISUSE),
		},
		{
			name: "Verification expired",
			args: args{
				token: verif2.Token,
				usage: verif2.Usage,
			},
			want:  dto.VerificationResponse{},
			want1: status.Error(status.VERIFICATION_TOKEN_EXPIRED),
		},
		{
			name: "Verification doesn't exists",
			args: args{
				token: util.GenerateRandomString(20),
				usage: verif2.Usage,
			},
			want:  dto.VerificationResponse{},
			want1: status.Error(status.VERIFICATION_TOKEN_NOT_FOUND),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := v.Response(tt.args.token, tt.args.usage)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Response() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Response() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
