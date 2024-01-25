package service

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"manga-explorer/internal/common/status"
	"manga-explorer/internal/domain/mangas"
	"manga-explorer/internal/domain/mangas/dto"
	"manga-explorer/internal/domain/mangas/mapper"
	genreRepoMock "manga-explorer/internal/domain/mangas/repository/mocks"
	"manga-explorer/internal/util"
	"manga-explorer/internal/util/containers"
	"reflect"
	"testing"
	"time"
)

func Test_mangaGenreService_CreateGenre(t *testing.T) {
	input := dto.GenreCreateInput{
		Name: "fantasy",
	}

	badInput := dto.GenreCreateInput{
		Name: "duplicated",
	}

	mockedGenreRepo := genreRepoMock.NewGenreMock(t)
	mockedGenreRepo.EXPECT().CreateGenre(&input).Return(nil)
	mockedGenreRepo.EXPECT().CreateGenre(&badInput).Return(sql.ErrNoRows)
	mockedGenreRepo.EXPECT().CreateGenre(mock.AnythingOfType("*dto.GenreCreateInput")).Return(simpleError)

	m := NewGenreService(mockedGenreRepo)
	type args struct {
		input dto.GenreCreateInput
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
			name: "Duplicate",
			args: args{
				input: badInput,
			},
			want: status.Error(status.GENRE_ALREADY_EXIST),
		},
		{
			name: "Normal",
			args: args{
				input: dto.GenreCreateInput{},
			},
			want: status.InternalError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := m.CreateGenre(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateGenre() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mangaGenreService_DeleteGenre(t *testing.T) {
	genreId := uuid.NewString()
	badGenreId := uuid.NewString()

	mockedGenreRepo := genreRepoMock.NewGenreMock(t)
	mockedGenreRepo.EXPECT().DeleteGenreById(&genreId).Return(nil)
	mockedGenreRepo.EXPECT().DeleteGenreById(&badGenreId).Return(sql.ErrNoRows)

	m := NewGenreService(mockedGenreRepo)
	type args struct {
		genreId string
	}
	tests := []struct {
		name string
		args args
		want status.Object
	}{
		{
			name: "Normal",
			args: args{
				genreId: genreId,
			},
			want: status.Success(),
		},
		{
			name: "Genre not found",
			args: args{
				genreId: badGenreId,
			},
			want: status.Error(status.GENRE_NOT_FOUND),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := m.DeleteGenre(tt.args.genreId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteGenre() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mangaGenreService_ListGenre(t *testing.T) {
	mangaGenres := []mangas.Genre{
		{
			Id:        uuid.NewString(),
			Name:      util.GenerateRandomString(20),
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
		},
		{
			Id:        uuid.NewString(),
			Name:      util.GenerateRandomString(20),
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
		},
	}

	mockedGenreRepo := genreRepoMock.NewGenreMock(t)
	mockedGenreRepo.EXPECT().ListGenres().Return(mangaGenres, nil)

	m := NewGenreService(mockedGenreRepo)

	tests := []struct {
		name  string
		want  []dto.GenreResponse
		want1 status.Object
	}{
		{
			name:  "Normal",
			want:  containers.CastSlicePtr(mangaGenres, mapper.ToGenreResponse),
			want1: status.Success(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := m.ListGenre()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListGenre() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ListGenre() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_mangaGenreService_UpdateGenre(t *testing.T) {
	updatedGenre := mangas.Genre{
		Id:        uuid.NewString(),
		Name:      util.GenerateRandomString(20),
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}
	badGenre := mangas.Genre{
		Id:        uuid.NewString(),
		Name:      util.GenerateRandomString(20),
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}

	mockedGenreRepo := genreRepoMock.NewGenreMock(t)
	mockedGenreRepo.EXPECT().UpdateGenre(&updatedGenre).Return(nil)
	mockedGenreRepo.EXPECT().UpdateGenre(&badGenre).Return(sql.ErrNoRows)
	//mockedGenreRepo.EXPECT().UpdateGenre(mock.AnythingOfType("*mangas.Genre")).Return(simpleError)

	m := NewGenreService(mockedGenreRepo)

	type args struct {
		input *dto.GenreUpdateInput
	}
	tests := []struct {
		name string
		args args
		want status.Object
	}{
		{
			name: "Normal",
			args: args{
				input: &dto.GenreUpdateInput{
					Id:   updatedGenre.Id,
					Name: updatedGenre.Name,
				},
			},
			want: status.Success(),
		},
		{
			name: "Genre not found",
			args: args{
				input: &dto.GenreUpdateInput{
					Id:   badGenre.Id,
					Name: badGenre.Name,
				},
			},
			want: status.Error(status.GENRE_NOT_FOUND),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := m.UpdateGenre(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateGenre() = %v, want %v", got, tt.want)
			}
		})
	}
}
