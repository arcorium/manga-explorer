package pg

import (
  "github.com/google/uuid"
  "github.com/stretchr/testify/require"
  "github.com/uptrace/bun"
  "manga-explorer/internal/domain/mangas"
  "manga-explorer/internal/util/opt"
  "testing"
)

func newGenreForTest(id opt.Optional[string], name string) *mangas.Genre {
  temp := mangas.NewGenre(name)
  if id.HasValue() {
    temp.Id = *id.Value()
  }
  return &temp
}

func Test_mangaGenreRepository_CreateGenre(t *testing.T) {

  type args struct {
    genre *mangas.Genre
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Normal",
      args: args{
        genre: newGenreForTest(opt.Null[string](), "asodniad"),
      },
      wantErr: false,
    },
    {
      name: "Duplicate Id",
      args: args{
        genre: newGenreForTest(opt.New("d1dc898b-4849-4b06-b613-5880bcf07cff"), "nasidnasd"),
      },
      wantErr: true,
    },
    {
      name: "Bad UUID",
      args: args{
        genre: newGenreForTest(opt.New("asdiasdasd"), "ncizxcnia"),
      },
      wantErr: true,
    },
    {
      name: "Duplicate Name",
      args: args{
        genre: newGenreForTest(opt.Null[string](), "comedy"),
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    tx, err := Db.Begin()
    require.NoError(t, err)

    m := NewMangaGenre(tx)
    t.Run(tt.name, func(t *testing.T) {
      defer func(tx2 bun.Tx) {
        require.NoError(t, tx2.Rollback())
      }(tx)

      if err := m.CreateGenre(tt.args.genre); (err != nil) != tt.wantErr {
        t.Errorf("CreateGenre() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func Test_mangaGenreRepository_DeleteGenreById(t *testing.T) {

  type args struct {
    genreId string
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Normal",
      args: args{
        genreId: "d1dc898b-4849-4b06-b613-5880bcf07cff",
      },
      wantErr: false,
    },
    {
      name: "Genre not exists",
      args: args{
        genreId: uuid.NewString(),
      },
      wantErr: true,
    },
    {
      name: "Bad UUID",
      args: args{
        genreId: "asdadsasd",
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    tx, err := Db.Begin()
    require.NoError(t, err)

    m := NewMangaGenre(tx)
    t.Run(tt.name, func(t *testing.T) {
      defer func(tx2 bun.Tx) {
        require.NoError(t, tx2.Rollback())
      }(tx)

      if err := m.DeleteGenreById(tt.args.genreId); (err != nil) != tt.wantErr {
        t.Errorf("DeleteGenreById() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}
