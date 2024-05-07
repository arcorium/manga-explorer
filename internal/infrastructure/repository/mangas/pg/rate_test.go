package pg

import (
  "github.com/google/uuid"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/require"
  "manga-explorer/internal/domain/mangas"
  "manga-explorer/internal/util/opt"
  "reflect"
  "testing"
)

func newRateForTest(id opt.Optional[string], userId, mangaId string, rate uint8) *mangas.Rate {
  temp := mangas.NewRate(userId, mangaId, rate)
  if id.HasValue() {
    temp.Id = *id.Value()
  }
  return &temp
}

func Test_mangaRateRepository_FindMangaRatings(t *testing.T) {
  type args struct {
    mangaId string
  }
  tests := []struct {
    name    string
    args    args
    want    []mangas.Rate
    wantErr bool
  }{
    {
      name: "Normal",
      args: args{
        mangaId: "3344af32-3393-4254-ba3a-d4ac03501259",
      },
      want: []mangas.Rate{
        *newRateForTest(opt.New("e1524426-d6a0-41bd-bffb-64f9140e5175"), "dd2166b0-5e62-4b74-b4cb-4be51a5040dc", "3344af32-3393-4254-ba3a-d4ac03501259", 5),
        *newRateForTest(opt.New("c1f84535-debc-4345-8bbd-1c8455c07bbe"), "73141bf0-5e64-4f52-acab-098f1efa3fa7", "3344af32-3393-4254-ba3a-d4ac03501259", 5),
        *newRateForTest(opt.New("949ecce3-99cf-4cba-ad10-992c5ee6d017"), "4afa29b2-d543-4489-b8ef-93f57781c9f6", "3344af32-3393-4254-ba3a-d4ac03501259", 10),
        *newRateForTest(opt.New("0b9c262e-d617-4518-853a-ed7dfb1665d8"), "4d704d17-8900-45d7-83a0-a10e4a4950d9", "3344af32-3393-4254-ba3a-d4ac03501259", 4),
        *newRateForTest(opt.New("d28341f8-576b-4f9d-ad4a-ef4ecf6d2183"), "f8b6a114-8cfe-4e14-b2fb-590c53cec0f1", "3344af32-3393-4254-ba3a-d4ac03501259", 3),
        *newRateForTest(opt.New("29baf1f4-7b54-4ee5-843a-febdb220b60b"), "c7760836-71e7-4664-99e8-a9503482a296", "3344af32-3393-4254-ba3a-d4ac03501259", 5),
        *newRateForTest(opt.New("16e3dd22-d6b8-4db0-b48b-98ddfd7c491d"), "db22a444-41a9-41db-ad8c-cb47759a98a8", "3344af32-3393-4254-ba3a-d4ac03501259", 2),
      },
      wantErr: false,
    },
    {
      name: "Bad UUID",
      args: args{
        mangaId: "asdadasd",
      },
      want:    nil,
      wantErr: true,
    },
    {
      name: "Manga not exists",
      args: args{
        mangaId: uuid.NewString(),
      },
      want:    nil,
      wantErr: true,
    },
    {
      name: "Manga has no ratings",
      args: args{
        mangaId: uuid.NewString(),
      },
      want:    nil,
      wantErr: true,
    },
  }
  for _, tt := range tests {
    m := NewMangaRate(Db)
    t.Run(tt.name, func(t *testing.T) {
      got, err := m.FindMangaRatings(tt.args.mangaId)
      if (err != nil) != tt.wantErr {
        t.Errorf("FindMangaRatings() error = %v, wantErr %v", err, tt.wantErr)
        return
      }

      require.Len(t, got, len(tt.want))

      for i := 0; i < len(got); i++ {
        // Remove time
        got[i].CreatedAt = tt.want[i].CreatedAt
        got[i].UpdatedAt = tt.want[i].UpdatedAt
        // Remove relations
        got[i].Manga = tt.want[i].Manga
        got[i].User = tt.want[i].User
      }

      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("FindMangaRatings() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_mangaRateRepository_Upsert(t *testing.T) {
  type args struct {
    rate *mangas.Rate
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Insert Rate",
      args: args{
        rate: newRateForTest(opt.Null[string](), "afcd4ab0-3190-4d35-885a-1d20eb909bd9", "19382f54-1da7-4cb7-807d-9f6030bb121e", 4),
      },
      wantErr: false,
    },
    {
      name: "Update Rate",
      args: args{
        rate: newRateForTest(opt.Null[string](), "4afa29b2-d543-4489-b8ef-93f57781c9f6", "4cf0e132-e209-4c5e-b47e-862de1fbcd95", 7), // last 3
      },
      wantErr: false,
    },
    {
      name: "Bad Id",
      args: args{
        rate: newRateForTest(opt.New("asdasdasd"), "afcd4ab0-3190-4d35-885a-1d20eb909bd9", "19382f54-1da7-4cb7-807d-9f6030bb121e", 4),
      },
      wantErr: true,
    },
    {
      name: "User not exist",
      args: args{
        rate: newRateForTest(opt.Null[string](), uuid.NewString(), "19382f54-1da7-4cb7-807d-9f6030bb121e", 4),
      },
      wantErr: true,
    },
    {
      name: "Manga not exist",
      args: args{
        rate: newRateForTest(opt.Null[string](), "afcd4ab0-3190-4d35-885a-1d20eb909bd9", uuid.NewString(), 4),
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    tx, err := Db.Begin()
    require.NoError(t, err)

    m := NewMangaRate(tx)
    t.Run(tt.name, func(t *testing.T) {
      defer func() {
        require.NoError(t, tx.Rollback())
      }()

      err := m.Upsert(tt.args.rate)
      if (err != nil) != tt.wantErr {
        t.Errorf("Upsert() error = %v, wantErr %v", err, tt.wantErr)
      }

      if err != nil {
        return
      }
      // Check for upsert and insert
      rate, err := m.FindRating(tt.args.rate.UserId, tt.args.rate.MangaId)
      require.NoError(t, err)
      assert.Equal(t, rate.Rate, tt.args.rate.Rate)
    })
  }
}

func Test_mangaRateRepository_FindRating(t *testing.T) {
  type args struct {
    userId  string
    mangaId string
  }
  tests := []struct {
    name    string
    args    args
    want    *mangas.Rate
    wantErr bool
  }{
    {
      name: "Normal",
      args: args{
        userId:  "4afa29b2-d543-4489-b8ef-93f57781c9f6",
        mangaId: "4cf0e132-e209-4c5e-b47e-862de1fbcd95",
      },
      want:    newRateForTest(opt.New("784c9dc4-5f2c-4c05-9498-17a866214ca5"), "4afa29b2-d543-4489-b8ef-93f57781c9f6", "4cf0e132-e209-4c5e-b47e-862de1fbcd95", 3),
      wantErr: false,
    },
    {
      name: "Bad UserId",
      args: args{
        userId:  "asdasd",
        mangaId: "4cf0e132-e209-4c5e-b47e-862de1fbcd95",
      },
      want:    nil,
      wantErr: true,
    },
    {
      name: "Bad mangaId",
      args: args{
        userId:  "4afa29b2-d543-4489-b8ef-93f57781c9f6",
        mangaId: "asdasdasda",
      },
      want:    nil,
      wantErr: true,
    },
    {
      name: "Rate not found",
      args: args{
        userId:  uuid.NewString(),
        mangaId: uuid.NewString(),
      },
      want:    nil,
      wantErr: true,
    },
  }
  for _, tt := range tests {
    m := NewMangaRate(Db)
    t.Run(tt.name, func(t *testing.T) {
      got, err := m.FindRating(tt.args.userId, tt.args.mangaId)
      if (err != nil) != tt.wantErr {
        t.Errorf("FindRating() error = %v, wantErr %v", err, tt.wantErr)
        return
      }

      if got != nil && tt.want != nil {
        // Remove time
        got.CreatedAt = tt.want.CreatedAt
        got.UpdatedAt = tt.want.UpdatedAt
        // Remove relations
        got.Manga = tt.want.Manga
        got.User = tt.want.User
      }

      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("FindRating() got = %v, want %v", got, tt.want)
      }
    })
  }
}
