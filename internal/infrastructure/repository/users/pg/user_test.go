package pg

import (
  "github.com/google/uuid"
  "github.com/stretchr/testify/require"
  "manga-explorer/database/fixtures"
  userEntity "manga-explorer/internal/domain/users"
  "manga-explorer/internal/util"
  "reflect"
  "testing"
  "time"
)

func Test_userRepository_CreateProfile(t *testing.T) {
  type args struct {
    profile *userEntity.users
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Referenced UserId",
      args: args{
        profile: &userEntity.users{
          UserId:    "4d704d17-8900-45d7-83a0-a10e4a4950d9",
          FirstName: "Test1",
          LastName:  "Test2",
          PhotoURL:  "test.png",
          Bio:       "hello",
        },
      },
      wantErr: false,
    },

    {
      name: "Double Referenced UserId",
      args: args{
        profile: &userEntity.users{
          UserId:    "c7760836-71e7-4664-99e8-a9503482a296",
          FirstName: "Test1",
        },
      },
      wantErr: true,
    },
    {
      name: "Non-referenced UserId",
      args: args{
        profile: &userEntity.users{
          UserId:    uuid.NewString(),
          FirstName: "Test1",
          LastName:  "Test2",
          PhotoURL:  "test.png",
          Bio:       "hello",
        },
      },
      wantErr: true,
    },
    {
      name: "Wrong uuid",
      args: args{
        profile: &userEntity.users{
          UserId:    "asdasd",
          FirstName: "Test1",
          LastName:  "Test2",
          PhotoURL:  "test.png",
          Bio:       "hello",
        },
      },
      wantErr: true,
    },
    {
      name: "Null data",
      args: args{
        profile: &userEntity.users{
          UserId: "73141bf0-5e64-4f52-acab-098f1efa3fa7",
        },
      },
      wantErr: true,
    },
  }

  for _, tt := range tests {
    tx, err := Db.Begin()
    require.NoError(t, err)
    //tx := Db

    repo := NewUser(tx)
    t.Run(tt.name, func(t *testing.T) {
      if err := repo.CreateProfile(tt.args.profile); (err != nil) != tt.wantErr {
        t.Errorf("CreateProfile() error = %v, wantErr %v", err, tt.wantErr)
      }
    })

    err = tx.Rollback()
    require.NoError(t, err)
  }
}

func Test_userRepository_CreateUser(t *testing.T) {
  type args struct {
    user    *userEntity.users
    profile *userEntity.users
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Normal",
      args: args{
        user: &userEntity.users{
          Id:       "390846ef-298e-470b-9a7a-1e58becc04c4",
          Username: "test1",
          Email:    "test1@gmail.com",
          Password: util.DropError(util.Hash("asd")),
          Verified: false,
          Role:     userEntity.RoleUser,
        },
        profile: &userEntity.users{
          UserId:    "390846ef-298e-470b-9a7a-1e58becc04c4",
          FirstName: "test",
          LastName:  "1",
        },
      },
      wantErr: false,
    },
    {
      name: "Wrong UUID",
      args: args{
        user: &userEntity.users{
          Id:       "asdasd",
          Username: "test1",
          Email:    "test1@gmail.com",
          Password: util.DropError(util.Hash("asd")),
          Verified: false,
          Role:     userEntity.RoleUser,
        },
        profile: &userEntity.users{
          UserId:    "asdasd",
          FirstName: "test",
          LastName:  "1",
        },
      },
      wantErr: true,
    },
    {
      name: "Wrong referenced profile",
      args: args{
        user: &userEntity.users{
          Id:       uuid.NewString(),
          Username: "test1",
          Email:    "test1@gmail.com",
          Password: util.DropError(util.Hash("asd")),
          Verified: false,
          Role:     userEntity.RoleUser,
        },
        profile: &userEntity.users{
          UserId:    uuid.NewString(),
          FirstName: "test",
          LastName:  "1",
        },
      },
      wantErr: true,
    },
    {
      name: "Email already exist",
      args: args{
        user: &userEntity.users{
          Id:       fixtures.GetConstantUUID(0),
          Username: "test1",
          Email:    "mcrocroft0@ask.com",
          Password: util.DropError(util.Hash("asd")),
          Verified: false,
          Role:     userEntity.RoleUser,
        },
        profile: &userEntity.users{
          UserId:    fixtures.GetConstantUUID(0),
          FirstName: "test",
          LastName:  "1",
        },
      },
      wantErr: true,
    },
    {
      name: "Username already exist",
      args: args{
        user: &userEntity.users{
          Id:       fixtures.GetConstantUUID(0),
          Username: "cmacfaell0",
          Email:    "test1@gmail.com",
          Password: util.DropError(util.Hash("asd")),
          Verified: false,
          Role:     userEntity.RoleUser,
        },
        profile: &userEntity.users{
          UserId:    fixtures.GetConstantUUID(0),
          FirstName: "test",
          LastName:  "1",
        },
      },
      wantErr: true,
    },
  }

  for _, tt := range tests {
    tx, err := Db.Begin()
    require.NoError(t, err)

    t.Run(tt.name, func(t *testing.T) {
      u := NewUser(tx)
      if err := u.CreateUser(tt.args.user, tt.args.profile); (err != nil) != tt.wantErr {
        t.Errorf("RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
      }
    })

    err = tx.Rollback()
    require.NoError(t, err)
  }
}

func Test_userRepository_DeleteUserById(t *testing.T) {
  type args struct {
    userId string
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Present User",
      args: args{
        userId: "4d704d17-8900-45d7-83a0-a10e4a4950d9",
      },
      wantErr: false,
    },
    {
      name: "Non-present User",
      args: args{
        userId: fixtures.GetConstantUUID(0),
      },
      wantErr: true,
    },
    {
      name: "Wrong UUID",
      args: args{
        userId: "c7760ds",
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    tx, err := Db.Begin()
    require.NoError(t, err)
    repo := NewUser(tx)

    t.Run(tt.name, func(t *testing.T) {
      if err := repo.DeleteUser(tt.args.userId); (err != nil) != tt.wantErr {
        t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
      }
    })

    err = tx.Rollback()
    require.NoError(t, err)
  }
}

func Test_userRepository_FindUserByEmail(t *testing.T) {
  type args struct {
    email string
  }
  tests := []struct {
    name    string
    args    args
    want    *userEntity.users
    wantErr bool
  }{
    {
      name: "User present",
      args: args{
        email: "mcrocroft0@ask.com",
      },
      want: &userEntity.users{
        Id:       "c7760836-71e7-4664-99e8-a9503482a296",
        Username: "cmacfaell0",
        Email:    "mcrocroft0@ask.com",
        Password: "$2a$04$e0Fhv.WS9bfP6IZupkh4DeAWx6C9UmUQ1tORctN7WEZaQcaeVfh/6",
        Verified: true,
        Role:     userEntity.RoleAdmin,
      },
      wantErr: false,
    },
    {
      name: "User not present",
      args: args{
        email: "not_present@gmail.com",
      },
      want:    &userEntity.users{},
      wantErr: true,
    },
  }
  for _, tt := range tests {
    tx, err := Db.Begin()
    require.NoError(t, err)
    repo := NewUser(tx)
    t.Run(tt.name, func(t *testing.T) {
      got, err := repo.FindUserByEmail(tt.args.email)
      if (err != nil) != tt.wantErr {
        t.Errorf("FindUserByEmail() error = %v, wantErr %v", err, tt.wantErr)
        return
      }

      // Ignore time fields
      got.CreatedAt = tt.want.CreatedAt
      got.UpdatedAt = tt.want.UpdatedAt

      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("FindUserByEmail() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_userRepository_FindUserById(t *testing.T) {
  type args struct {
    id string
  }
  tests := []struct {
    name    string
    args    args
    want    *userEntity.users
    wantErr bool
  }{
    {
      name: "Present User",
      args: args{
        id: "c7760836-71e7-4664-99e8-a9503482a296",
      },
      want: &userEntity.users{
        Id:       "c7760836-71e7-4664-99e8-a9503482a296",
        Username: "cmacfaell0",
        Email:    "mcrocroft0@ask.com",
        Password: "$2a$04$e0Fhv.WS9bfP6IZupkh4DeAWx6C9UmUQ1tORctN7WEZaQcaeVfh/6",
        Verified: true,
        Role:     userEntity.RoleAdmin,
      },
      wantErr: false,
    },
    {
      name: "User not present",
      args: args{
        id: fixtures.GetConstantUUID(0),
      },
      want:    &userEntity.users{},
      wantErr: true,
    },
    {
      name: "Wrong UUID",
      args: args{
        id: "asdasd",
      },
      want:    &userEntity.users{},
      wantErr: true,
    },
  }
  for _, tt := range tests {
    tx, err := Db.Begin()
    require.NoError(t, err)
    repo := NewUser(tx)

    t.Run(tt.name, func(t *testing.T) {
      got, err := repo.FindUserById(tt.args.id)
      if (err != nil) != tt.wantErr {
        t.Errorf("FindUserById() error = %v, wantErr %v", err, tt.wantErr)
        return
      }

      // Ignore time fields
      got.CreatedAt = tt.want.CreatedAt
      got.UpdatedAt = tt.want.UpdatedAt

      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("FindUserById() got = %v, want %v", got, tt.want)
      }
    })

    err = tx.Rollback()
    require.NoError(t, err)
  }
}

//func Test_userRepository_FindUserMangaFavorites(t *testing.T) {
//	type args struct {
//		userId string
//		offset uint64
//		limit  uint64
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    []mangas.MangaFavorite
//		want1   uint64
//		wantErr bool
//	}{
//		{
//			name: "No Manga Favorites",
//			args: args{
//				userId: "f8b6a114-8cfe-4e14-b2fb-590c53cec0f1",
//				offset: 0,
//				limit:  0,
//			},
//			want:    nil,
//			want1:   0,
//			wantErr: false,
//		},
//		{
//			name: "Out limit",
//			args: args{
//				userId: "4afa29b2-d543-4489-b8ef-93f57781c9f6", // Has 7 manga favorites
//				offset: 0,
//				limit:  10,
//			},
//			want: []mangas.MangaFavorite{
//				{
//					UserId:  "4afa29b2-d543-4489-b8ef-93f57781c9f6",
//					MangaId: "4ab94eda-46ae-4de8-bd3a-734b388e06fc",
//				},
//				{
//					UserId:  "4afa29b2-d543-4489-b8ef-93f57781c9f6",
//					MangaId: "fc1bea74-5fde-4cf0-a332-c957c914d121",
//				},
//				{
//					UserId:  "4afa29b2-d543-4489-b8ef-93f57781c9f6",
//					MangaId: "0afda87c-944e-4cc0-8d51-99eea96844af",
//				},
//				{
//					UserId:  "4afa29b2-d543-4489-b8ef-93f57781c9f6",
//					MangaId: "5b93c8f4-efab-46b1-bae5-130389e0640b",
//				},
//				{
//					UserId:  "4afa29b2-d543-4489-b8ef-93f57781c9f6",
//					MangaId: "cd856f30-853b-45f8-b215-17d23a4f6ffb",
//				},
//				{
//					UserId:  "4afa29b2-d543-4489-b8ef-93f57781c9f6",
//					MangaId: "639e0925-fe96-47c4-8fbf-9ea6abdb4c75",
//				},
//				{
//					UserId:  "4afa29b2-d543-4489-b8ef-93f57781c9f6",
//					MangaId: "24ea4aec-4409-4727-8644-c8e010c56f15",
//				},
//			},
//			want1:   7,
//			wantErr: false,
//		},
//		{
//			name: "Using offset out limit",
//			args: args{
//				userId: "4afa29b2-d543-4489-b8ef-93f57781c9f6", // Has 7 manga favorites
//				offset: 2,
//				limit:  10,
//			},
//			want: []mangas.MangaFavorite{
//				{
//					UserId:  "4afa29b2-d543-4489-b8ef-93f57781c9f6",
//					MangaId: "0afda87c-944e-4cc0-8d51-99eea96844af",
//				},
//				{
//					UserId:  "4afa29b2-d543-4489-b8ef-93f57781c9f6",
//					MangaId: "5b93c8f4-efab-46b1-bae5-130389e0640b",
//				},
//				{
//					UserId:  "4afa29b2-d543-4489-b8ef-93f57781c9f6",
//					MangaId: "cd856f30-853b-45f8-b215-17d23a4f6ffb",
//				},
//				{
//					UserId:  "4afa29b2-d543-4489-b8ef-93f57781c9f6",
//					MangaId: "639e0925-fe96-47c4-8fbf-9ea6abdb4c75",
//				},
//				{
//					UserId:  "4afa29b2-d543-4489-b8ef-93f57781c9f6",
//					MangaId: "24ea4aec-4409-4727-8644-c8e010c56f15",
//				},
//			},
//			want1:   7,
//			wantErr: false,
//		},
//		{
//			name: "Normal pagination",
//			args: args{
//				userId: "4afa29b2-d543-4489-b8ef-93f57781c9f6", // Has 7 manga favorites
//				offset: 2,
//				limit:  3,
//			},
//			want: []mangas.MangaFavorite{
//				{
//					UserId:  "4afa29b2-d543-4489-b8ef-93f57781c9f6",
//					MangaId: "0afda87c-944e-4cc0-8d51-99eea96844af",
//				},
//				{
//					UserId:  "4afa29b2-d543-4489-b8ef-93f57781c9f6",
//					MangaId: "5b93c8f4-efab-46b1-bae5-130389e0640b",
//				},
//				{
//					UserId:  "4afa29b2-d543-4489-b8ef-93f57781c9f6",
//					MangaId: "cd856f30-853b-45f8-b215-17d23a4f6ffb",
//				},
//			},
//			want1:   7,
//			wantErr: false,
//		},
//		{
//			name: "Last index pagination",
//			args: args{
//				userId: "4afa29b2-d543-4489-b8ef-93f57781c9f6", // Has 7 manga favorites
//				offset: 6,
//				limit:  3,
//			},
//			want: []mangas.MangaFavorite{
//				{
//					UserId:  "4afa29b2-d543-4489-b8ef-93f57781c9f6",
//					MangaId: "24ea4aec-4409-4727-8644-c8e010c56f15",
//				},
//			},
//			want1:   7,
//			wantErr: false,
//		},
//		{
//			name: "max index offset",
//			args: args{
//				userId: "4afa29b2-d543-4489-b8ef-93f57781c9f6", // Has 7 manga favorites
//				offset: 7,
//				limit:  2,
//			},
//			want:    nil,
//			want1:   7,
//			wantErr: false,
//		},
//		{
//			name: "Out index offset",
//			args: args{
//				userId: "4afa29b2-d543-4489-b8ef-93f57781c9f6", // Has 7 manga favorites
//				offset: 8,
//				limit:  2,
//			},
//			want:    nil,
//			want1:   7,
//			wantErr: false,
//		},
//		{
//			name: "User not exists",
//			args: args{
//				userId: uuid.NewString(),
//				offset: 0,
//				limit:  0,
//			},
//			want:    nil,
//			want1:   0,
//			wantErr: false,
//		},
//		{
//			name: "Invalid user id",
//			args: args{
//				userId: "asdadasd-123dani",
//				offset: 0,
//				limit:  0,
//			},
//			want:    nil,
//			want1:   0,
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			tx, err := Db.Begin()
//			require.NoError(t, err)
//			repo := NewUser(tx)
//
//			got, got1, err := repo.FindMangaFavorites(tt.args.userId, tt.args.offset, tt.args.limit)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("FindMangaFavorites() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//
//			if got1 != tt.want1 {
//				t.Errorf("FindMangaFavorites() got1 = %v, want %v", got1, tt.want1)
//			}
//
//			// Ignore time fields
//			for i := 0; i < len(got); i++ {
//				got[i].CreatedAt = tt.want[i].CreatedAt
//			}
//
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("FindMangaFavorites() got = %v, want %v", got, tt.want)
//			}
//
//			defer tx.Rollback()
//
//		})
//	}
//}

//func Test_userRepository_FindUserMangaHistories(t *testing.T) {
//	type fields struct {
//		db *bun.DB
//	}
//	type args struct {
//		userId string
//		offset uint64
//		limit  uint64
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    []mangas.MangaHistory
//		want1   uint64
//		wantErr bool
//	}{
//		{
//			name: "No Manga histories",
//			args: args{
//				userId: "4d704d17-8900-45d7-83a0-a10e4a4950d9",
//				offset: 0,
//				limit:  0,
//			},
//			want:    nil,
//			want1:   0,
//			wantErr: false,
//		},
//		{
//			name: "Out limit",
//			args: args{
//				userId: "73141bf0-5e64-4f52-acab-098f1efa3fa7", // Has 5 manga histories
//				offset: 0,
//				limit:  10,
//			},
//			want: []mangas.MangaHistory{
//				{
//					UserId:  "73141bf0-5e64-4f52-acab-098f1efa3fa7",
//					MangaId: "35d1bea2-1a13-45e7-a08c-5d35db26444d",
//				},
//				{
//					UserId:  "73141bf0-5e64-4f52-acab-098f1efa3fa7",
//					MangaId: "0952a563-671c-4a1a-93db-09d6bf64a82b",
//				},
//				{
//					UserId:  "73141bf0-5e64-4f52-acab-098f1efa3fa7",
//					MangaId: "6dd52489-0ba2-4e84-9c0d-1666a63e1699",
//				},
//				{
//					UserId:  "73141bf0-5e64-4f52-acab-098f1efa3fa7",
//					MangaId: "9121bb87-33d8-465c-823a-6e79278eb768",
//				},
//				{
//					UserId:  "73141bf0-5e64-4f52-acab-098f1efa3fa7",
//					MangaId: "e73192c5-b5c4-44ff-9b13-04f4e2983010",
//				},
//			},
//			want1:   5,
//			wantErr: false,
//		},
//		{
//			name: "Using offset out limit",
//			args: args{
//				userId: "73141bf0-5e64-4f52-acab-098f1efa3fa7", // Has 7 manga favorites
//				offset: 2,
//				limit:  10,
//			},
//			want: []mangas.MangaHistory{
//				{
//					UserId:  "73141bf0-5e64-4f52-acab-098f1efa3fa7",
//					MangaId: "6dd52489-0ba2-4e84-9c0d-1666a63e1699",
//				},
//				{
//					UserId:  "73141bf0-5e64-4f52-acab-098f1efa3fa7",
//					MangaId: "9121bb87-33d8-465c-823a-6e79278eb768",
//				},
//				{
//					UserId:  "73141bf0-5e64-4f52-acab-098f1efa3fa7",
//					MangaId: "e73192c5-b5c4-44ff-9b13-04f4e2983010",
//				},
//			},
//			want1:   5,
//			wantErr: false,
//		},
//		{
//			name: "Normal pagination",
//			args: args{
//				userId: "73141bf0-5e64-4f52-acab-098f1efa3fa7", // Has 7 manga favorites
//				offset: 1,
//				limit:  3,
//			},
//			want: []mangas.MangaHistory{
//				{
//					UserId:  "73141bf0-5e64-4f52-acab-098f1efa3fa7",
//					MangaId: "0952a563-671c-4a1a-93db-09d6bf64a82b",
//				},
//				{
//					UserId:  "73141bf0-5e64-4f52-acab-098f1efa3fa7",
//					MangaId: "6dd52489-0ba2-4e84-9c0d-1666a63e1699",
//				},
//				{
//					UserId:  "73141bf0-5e64-4f52-acab-098f1efa3fa7",
//					MangaId: "9121bb87-33d8-465c-823a-6e79278eb768",
//				},
//			},
//			want1:   5,
//			wantErr: false,
//		},
//		{
//			name: "max index offset",
//			args: args{
//				userId: "73141bf0-5e64-4f52-acab-098f1efa3fa7", // Has 7 manga histories
//				offset: 5,
//				limit:  2,
//			},
//			want:    nil,
//			want1:   5,
//			wantErr: false,
//		},
//		{
//			name: "Out index offset",
//			args: args{
//				userId: "73141bf0-5e64-4f52-acab-098f1efa3fa7", // Has 7 manga histories
//				offset: 6,
//				limit:  2,
//			},
//			want:    nil,
//			want1:   5,
//			wantErr: false,
//		},
//		{
//			name: "last item limit",
//			args: args{
//				userId: "73141bf0-5e64-4f52-acab-098f1efa3fa7", // Has 5 manga histories
//				offset: 4,
//				limit:  10,
//			},
//			want: []mangas.MangaHistory{
//				{
//					UserId:  "73141bf0-5e64-4f52-acab-098f1efa3fa7",
//					MangaId: "e73192c5-b5c4-44ff-9b13-04f4e2983010",
//				},
//			},
//			want1:   5,
//			wantErr: false,
//		},
//		{
//			name: "User not exists",
//			args: args{
//				userId: uuid.NewString(),
//				offset: 0,
//				limit:  0,
//			},
//			want:    nil,
//			want1:   0,
//			wantErr: false,
//		},
//		{
//			name: "Invalid user id",
//			args: args{
//				userId: "asdadasd-123dani",
//				offset: 0,
//				limit:  0,
//			},
//			want:    nil,
//			want1:   0,
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			tx, err := Db.Begin()
//			require.NoError(t, err)
//			defer tx.Rollback()
//
//			repo := NewUser(tx)
//
//			got, got1, err := repo.FindMangaHistories(tt.args.userId, tt.args.offset, tt.args.limit)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("FindMangaHistories() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if got1 != tt.want1 {
//				t.Errorf("FindMangaHistories() got1 = %v, want %v", got1, tt.want1)
//			}
//
//			// Ignore time fields
//			for i := 0; i < len(got); i++ {
//				got[i].LastView = tt.want[i].LastView
//			}
//
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("FindMangaHistories() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func Test_userRepository_UpdateProfile(t *testing.T) {
  type args struct {
    profile *userEntity.users
  }
  tests := []struct {
    name    string
    args    args
    want    *userEntity.users
    wantErr bool
  }{
    {
      name: "invalid UUID",
      args: args{
        profile: &userEntity.users{
          UserId:    "asudbaisd-asdjn",
          FirstName: "Lamont",
          LastName:  "Scramage",
          PhotoURL:  "/odio/porttitor/id.html",
          Bio:       "quis libero nullam sit amet turpis elementum ligula vehicula consequat morbi a ipsum integer a nibh in",
          UpdatedAt: time.Now(),
        },
      },
      want:    nil,
      wantErr: true,
    },
    {
      name: "no such user",
      args: args{
        profile: &userEntity.users{
          UserId:    uuid.NewString(),
          FirstName: "Lamont",
          LastName:  "Scramage",
          PhotoURL:  "/odio/porttitor/id.html",
          Bio:       "quis libero nullam sit amet turpis elementum ligula vehicula consequat morbi a ipsum integer a nibh in",
          UpdatedAt: time.Now(),
        },
      },
      want:    &userEntity.users{},
      wantErr: false,
    },
    {
      name: "update first name",
      args: args{
        profile: &userEntity.users{
          UserId:    "c7760836-71e7-4664-99e8-a9503482a296",
          FirstName: "Lamont",
          UpdatedAt: time.Now(),
        },
      },
      want: &userEntity.users{
        Id:        1,
        UserId:    "c7760836-71e7-4664-99e8-a9503482a296",
        FirstName: "Lamont",
        LastName:  "Paternoster",
        PhotoURL:  "/felis/donec/semper/sapien/a/libero.js",
        Bio:       "a odio in hac habitasse platea dictumst maecenas ut massa quis augue luctus tincidunt",
        UpdatedAt: time.Now(),
      },
      wantErr: false,
    },
    {
      name: "update name",
      args: args{
        profile: &userEntity.users{
          UserId:    "c7760836-71e7-4664-99e8-a9503482a296",
          FirstName: "Lamont",
          LastName:  "Scramage",
          UpdatedAt: time.Now(),
        },
      },
      want: &userEntity.users{
        Id:        1,
        UserId:    "c7760836-71e7-4664-99e8-a9503482a296",
        FirstName: "Lamont",
        LastName:  "Scramage",
        PhotoURL:  "/felis/donec/semper/sapien/a/libero.js",
        Bio:       "a odio in hac habitasse platea dictumst maecenas ut massa quis augue luctus tincidunt",
        UpdatedAt: time.Now(),
      },
      wantErr: false,
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      tx, err := Db.Begin()
      require.NoError(t, err)
      defer tx.Rollback()

      repo := NewUser(tx)
      if err := repo.UpdateProfileByUserId(tt.args.profile); (err != nil) != tt.wantErr {
        t.Errorf("UpdateProfileByUserId() error = %v, wantErr %v", err, tt.wantErr)
      }

      if tt.wantErr {
        return
      }

      got, err := repo.FindUserProfiles(tt.args.profile.UserId)

      got.UpdatedAt = tt.want.UpdatedAt

      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("UpdateProfileByUserId() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_userRepository_UpdateUserById(t *testing.T) {
  type args struct {
    user *userEntity.users
  }

  tests := []struct {
    name    string
    args    args
    want    *userEntity.users
    wantErr bool
  }{
    {
      name: "Invalid UUID",
      args: args{
        user: &userEntity.users{
          Id: "asjdbaiw-u12b3k",
        },
      },
      want:    nil,
      wantErr: true,
    },
    {
      name: "No such user",
      args: args{
        user: &userEntity.users{
          Id:        uuid.NewString(),
          Username:  "asd",
          Email:     "asdasd@gmail.com",
          Password:  "$2a$10$9tj4AP3nvqOBHrbfU.jMpucjfFulx7aH/bBEhN6NnWXflkP3.W.w.",
          Verified:  false,
          Role:      0,
          UpdatedAt: time.Now(),
        },
      },
      want:    &userEntity.users{},
      wantErr: true,
    },
    {
      name: "Change username",
      args: args{
        user: &userEntity.users{
          Id:        "a11b349d-59eb-4ebf-9276-eeaec5bdeacc",
          Username:  "asd3",
          UpdatedAt: time.Now(),
        },
      },
      want: &userEntity.users{
        Id:        "a11b349d-59eb-4ebf-9276-eeaec5bdeacc",
        Username:  "asd3",
        Email:     "afawlks2@hao123.com",
        Password:  "$2a$04$VSRUt2udPKrzS3yXKlDF8.zKHJYzbk02h1yO8jMJJAJHOmxT94Iia",
        Verified:  false,
        Role:      userEntity.RoleUser,
        UpdatedAt: time.Now(),
      },
      wantErr: false,
    },
    {
      name: "Change multiple fields",
      args: args{
        user: &userEntity.users{
          Id:        "a11b349d-59eb-4ebf-9276-eeaec5bdeacc",
          Verified:  true,
          Role:      userEntity.RoleAdmin,
          UpdatedAt: time.Now(),
        },
      },
      want: &userEntity.users{
        Id:        "a11b349d-59eb-4ebf-9276-eeaec5bdeacc",
        Username:  "mroan2",
        Email:     "afawlks2@hao123.com",
        Password:  "$2a$04$VSRUt2udPKrzS3yXKlDF8.zKHJYzbk02h1yO8jMJJAJHOmxT94Iia",
        Verified:  true,
        Role:      userEntity.RoleAdmin,
        UpdatedAt: time.Now(),
      },
      wantErr: false,
    },
    {
      name: "Change multiple fields",
      args: args{
        user: &userEntity.users{
          Id:        "a11b349d-59eb-4ebf-9276-eeaec5bdeacc",
          Email:     "asdiasd@gmail.com",
          Password:  "$2a$10$z4DVVjqUKr2PB1twekHo0./5I86mXHR779Ix9ueBXIN5dsYjvKOAu",
          UpdatedAt: time.Now(),
        },
      },
      want: &userEntity.users{
        Id:        "a11b349d-59eb-4ebf-9276-eeaec5bdeacc",
        Username:  "mroan2",
        Email:     "asdiasd@gmail.com",
        Password:  "$2a$10$z4DVVjqUKr2PB1twekHo0./5I86mXHR779Ix9ueBXIN5dsYjvKOAu",
        Verified:  false,
        Role:      userEntity.RoleUser,
        UpdatedAt: time.Now(),
      },
      wantErr: false,
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      tx, err := Db.Begin()
      require.NoError(t, err)
      defer tx.Rollback()

      repo := NewUser(tx)
      if err := repo.UpdateUser(tt.args.user); (err != nil) != tt.wantErr {
        t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
      }

      if tt.wantErr {
        return
      }

      got, err := repo.FindUserById(tt.args.user.Id)

      got.UpdatedAt = tt.want.UpdatedAt
      got.CreatedAt = tt.want.CreatedAt

      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("UpdateUserById() got = %v, want %v", got, tt.want)
      }
    })
  }
}
