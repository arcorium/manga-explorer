package pg

import (
  "github.com/google/uuid"
  "github.com/stretchr/testify/require"
  "github.com/uptrace/bun"
  "manga-explorer/internal/domain/users"
  "manga-explorer/internal/infrastructure/repository/authentication/pg"
  "manga-explorer/internal/util/opt"
  "reflect"
  "testing"
)

func createCredentialForTest(userId string, deviceName, accessTokenId, token opt.Optional[string]) *users.Credential {
  temp := users.NewCredential2(userId, deviceName.ValueOr("Test"), accessTokenId.ValueOr(uuid.NewString()), token.ValueOr(uuid.NewString()))
  return &temp
}

func Test_credentialRepository_Create(t *testing.T) {
  type args struct {
    cred *users.Credential
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Normal",
      args: args{
        cred: createCredentialForTest("c7760836-71e7-4664-99e8-a9503482a296", opt.NullStr, opt.NullStr, opt.NullStr),
      },
      wantErr: false,
    },
    {
      name: "Duplicate Credential Id",
      args: args{
        cred: &pg.Credentials[0],
      },
      wantErr: true,
    },
    {
      name: "User not exists",
      args: args{
        cred: createCredentialForTest(uuid.NewString(), opt.NullStr, opt.NullStr, opt.NullStr),
      },
      wantErr: true,
    },
    {
      name: "Bad User uuids",
      args: args{
        cred: createCredentialForTest("asdasda0abczu-asd", opt.NullStr, opt.NullStr, opt.NullStr),
      },
      wantErr: true,
    },
    {
      name: "Bad Token Access Ids",
      args: args{
        cred: createCredentialForTest("c7760836-71e7-4664-99e8-a9503482a296", opt.NullStr, opt.New("asdasdasvda-askdabc"), opt.NullStr),
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    tx, err := pg.Db.Begin()
    require.NoError(t, err)
    c := NewCredential(tx)

    t.Run(tt.name, func(t *testing.T) {
      if err := c.Create(tt.args.cred); (err != nil) != tt.wantErr {
        t.Errorf("Upsert() error = %v, wantErr %v", err, tt.wantErr)
      }
      defer func(tx bun.Tx) {
        require.NoError(t, tx.Rollback())
      }(tx)
    })
  }
}

func Test_credentialRepository_Find(t *testing.T) {
  type args struct {
    userId string
    credId string
  }
  tests := []struct {
    name    string
    args    args
    want    *users.Credential
    wantErr bool
  }{
    {
      name: "Normal",
      args: args{
        userId: pg.Credentials[0].UserId,
        credId: pg.Credentials[0].Id,
      },
      want:    &pg.Credentials[0],
      wantErr: false,
    },
    {
      name: "Bad UserId",
      args: args{
        userId: "asdadasdbuasd",
        credId: pg.Credentials[0].Id,
      },
      want:    nil,
      wantErr: true,
    },

    {
      name: "User not found",
      args: args{
        userId: uuid.NewString(),
        credId: pg.Credentials[0].Id,
      },
      want:    nil,
      wantErr: true,
    },
    {
      name: "Bad CredId",
      args: args{
        userId: pg.Credentials[0].UserId,
        credId: "asdasudaisvudasd",
      },
      want:    nil,
      wantErr: true,
    },
    {
      name: "Credential not found",
      args: args{
        userId: pg.Credentials[0].UserId,
        credId: uuid.NewString(),
      },
      want:    nil,
      wantErr: true,
    },
    {
      name: "Bad user and cred id",
      args: args{
        userId: "asdasdasdas",
        credId: "zxcnodnboqiw",
      },
      want:    nil,
      wantErr: true,
    },
  }
  for _, tt := range tests {
    c := NewCredential(pg.Db)
    t.Run(tt.name, func(t *testing.T) {

      got, err := c.Find(tt.args.userId, tt.args.credId)
      if (err != nil) != tt.wantErr {
        t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
        return
      }

      if got != nil && tt.want != nil {
        // Ignore times
        got.UpdatedAt = tt.want.UpdatedAt
        got.CreatedAt = tt.want.CreatedAt
        // Ignore relations
        got.User = tt.want.User
      }

      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("Find() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_credentialRepository_FindByAccessTokenId(t *testing.T) {
  type args struct {
    accessTokenId string
  }
  tests := []struct {
    name    string
    args    args
    want    *users.Credential
    wantErr bool
  }{
    {
      name: "Normal",
      args: args{
        accessTokenId: pg.Credentials[0].AccessTokenId,
      },
      want:    &pg.Credentials[0],
      wantErr: false,
    },
    {
      name: "Bad access token id as uuid",
      args: args{
        accessTokenId: "asdadszxcz0213w-sd",
      },
      want:    nil,
      wantErr: true,
    },
    {
      name: "Credential not found",
      args: args{
        accessTokenId: uuid.NewString(),
      },
      want:    nil,
      wantErr: true,
    },
  }
  for _, tt := range tests {

    c := NewCredential(pg.Db)
    t.Run(tt.name, func(t *testing.T) {
      got, err := c.FindByAccessTokenId(tt.args.accessTokenId)
      if (err != nil) != tt.wantErr {
        t.Errorf("FindByAccessTokenId() error = %v, wantErr %v", err, tt.wantErr)
        return
      }

      if got != nil && tt.want != nil {
        // Ignore times
        got.UpdatedAt = tt.want.UpdatedAt
        got.CreatedAt = tt.want.CreatedAt
        // Ignore relations
        got.User = tt.want.User
      }

      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("FindByAccessTokenId() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_credentialRepository_FindUserCredentials(t *testing.T) {
  type args struct {
    userId string
  }
  tests := []struct {
    name    string
    args    args
    want    []users.Credential
    wantErr bool
  }{
    {
      name: "Normal",
      args: args{
        userId: pg.Credentials[0].UserId,
      },
      want:    pg.Credentials,
      wantErr: false,
    },
    {
      name: "Bad user id as uuid",
      args: args{
        userId: "asdasd-zxcbad",
      },
      want:    nil,
      wantErr: true,
    },
    {
      name: "user has no credentials",
      args: args{
        userId: uuid.NewString(),
      },
      want:    nil,
      wantErr: true,
    },
    {
      name: "empty user id",
      args: args{
        userId: "",
      },
      want:    nil,
      wantErr: true,
    },
  }
  for _, tt := range tests {
    c := NewCredential(pg.Db)
    t.Run(tt.name, func(t *testing.T) {
      got, err := c.FindUserCredentials(tt.args.userId)
      if (err != nil) != tt.wantErr {
        t.Errorf("FindUserCredentials() error = %v, wantErr %v", err, tt.wantErr)
        return
      }

      if got != nil && tt.want != nil {
        require.Len(t, got, len(tt.want))
        for i := 0; i < len(got); i++ {
          // Ignore times
          got[i].UpdatedAt = tt.want[i].UpdatedAt
          got[i].CreatedAt = tt.want[i].CreatedAt
          // Ignore relations
          got[i].User = tt.want[i].User
        }
      }

      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("FindUserCredentials() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_credentialRepository_Remove(t *testing.T) {
  type args struct {
    userId string
    credId string
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Normal",
      args: args{
        userId: pg.Credentials[0].UserId,
        credId: pg.Credentials[0].Id,
      },
      wantErr: false,
    },
    {
      name: "User not found",
      args: args{
        userId: uuid.NewString(),
        credId: pg.Credentials[0].Id,
      },
      wantErr: true,
    },
    {
      name: "Credential of user doesn't exist",
      args: args{
        userId: pg.Credentials[0].UserId,
        credId: uuid.NewString(),
      },
      wantErr: true,
    },
    {
      name: "Bad user id as uuid",
      args: args{
        userId: "asdad0zxczxc-asdad",
        credId: pg.Credentials[0].Id,
      },
      wantErr: true,
    },
    {
      name: "Bad credential id as uuid",
      args: args{
        userId: pg.Credentials[0].UserId,
        credId: "xzczczx-dasdnb",
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {

    tx, err := pg.Db.Begin()
    require.NoError(t, err)
    c := NewCredential(tx)

    t.Run(tt.name, func(t *testing.T) {
      if err := c.Remove(tt.args.userId, tt.args.credId); (err != nil) != tt.wantErr {
        t.Errorf("Remove() error = %v, wantErr %v", err, tt.wantErr)
      }

      defer func(tx bun.Tx) {
        require.NoError(t, tx.Rollback())
      }(tx)
    })
  }
}

func Test_credentialRepository_RemoveByAccessTokenId(t *testing.T) {
  type args struct {
    userId        string
    accessTokenId string
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Normal",
      args: args{
        userId:        pg.Credentials[0].UserId,
        accessTokenId: pg.Credentials[0].AccessTokenId,
      },
      wantErr: false,
    },
    {
      name: "User doesn't exists",
      args: args{
        userId:        uuid.NewString(),
        accessTokenId: pg.Credentials[0].AccessTokenId,
      },
      wantErr: true,
    },
    {
      name: "Access token doesn't exists",
      args: args{
        userId:        pg.Credentials[0].UserId,
        accessTokenId: uuid.NewString(),
      },
      wantErr: true,
    },
    {
      name: "Bad user id as uuid",
      args: args{
        userId:        "asdadadas-zxcznca",
        accessTokenId: pg.Credentials[0].AccessTokenId,
      },
      wantErr: true,
    },
    {
      name: "Bad access token id as uuid",
      args: args{
        userId:        pg.Credentials[0].UserId,
        accessTokenId: "asdasdaszxc-asdnas",
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    tx, err := pg.Db.Begin()
    require.NoError(t, err)
    c := NewCredential(tx)

    t.Run(tt.name, func(t *testing.T) {
      if err := c.RemoveByAccessTokenId(tt.args.userId, tt.args.accessTokenId); (err != nil) != tt.wantErr {
        t.Errorf("RemoveByAccessTokenId() error = %v, wantErr %v", err, tt.wantErr)
      }

      defer func(tx bun.Tx) {
        require.NoError(t, tx.Rollback())
      }(tx)
    })
  }
}

func Test_credentialRepository_RemoveUserCredentials(t *testing.T) {
  type args struct {
    userId string
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Normal",
      args: args{
        userId: pg.Credentials[0].UserId,
      },
      wantErr: false,
    },
    {
      name: "User doesn't have credentials",
      args: args{
        userId: "dc4402e4-0f88-400a-978e-8bb3880ab063",
      },
      wantErr: true,
    },
    {
      name: "User doesn't exists",
      args: args{
        userId: uuid.NewString(),
      },
      wantErr: true,
    },
    {
      name: "Bad user id as uuid",
      args: args{
        userId: "asdasdasd-zxczca",
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    tx, err := pg.Db.Begin()
    require.NoError(t, err)
    c := NewCredential(tx)

    t.Run(tt.name, func(t *testing.T) {
      if err := c.RemoveUserCredentials(tt.args.userId); (err != nil) != tt.wantErr {
        t.Errorf("RemoveUserCredentials() error = %v, wantErr %v", err, tt.wantErr)
      }
      defer func(tx bun.Tx) {
        require.NoError(t, tx.Rollback())
      }(tx)
    })
  }
}

func Test_credentialRepository_UpdateAccessTokenId(t *testing.T) {
  type args struct {
    credentialId  string
    accessTokenId string
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Normal",
      args: args{
        credentialId:  pg.Credentials[0].Id,
        accessTokenId: uuid.NewString(),
      },
      wantErr: false,
    },
    {
      name: "Credential not found",
      args: args{
        credentialId:  uuid.NewString(),
        accessTokenId: uuid.NewString(),
      },
      wantErr: true,
    },
    {
      name: "Bad credential id as uuid",
      args: args{
        credentialId:  "asdadasd-zxczxdasd",
        accessTokenId: uuid.NewString(),
      },
      wantErr: true,
    },
    {
      name: "Bad access token id as uuid",
      args: args{
        credentialId:  pg.Credentials[0].Id,
        accessTokenId: "asdasdas-zxczcased",
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    tx, err := pg.Db.Begin()
    require.NoError(t, err)
    c := NewCredential(tx)

    t.Run(tt.name, func(t *testing.T) {
      if err := c.UpdateAccessTokenId(tt.args.credentialId, tt.args.accessTokenId); (err != nil) != tt.wantErr {
        t.Errorf("UpdateAccessTokenId() error = %v, wantErr %v", err, tt.wantErr)
      }
      defer func(tx bun.Tx) {
        require.NoError(t, tx.Rollback())
      }(tx)
    })
  }
}
