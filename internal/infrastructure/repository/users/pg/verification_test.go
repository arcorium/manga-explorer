package pg

import (
  "github.com/google/uuid"
  "github.com/stretchr/testify/require"
  "manga-explorer/internal/domain/users"
  "reflect"
  "testing"
  "time"
)

func Test_verificationRepository_Create(t *testing.T) {
  type args struct {
    verif *users.users
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "invalid user id",
      args: args{
        verif: &users.Verification{
          Token:          "Svx79uIDYN7NzDFU2EVWgxo3RX",
          UserId:         "asdasdzxcdas",
          Usage:          users.UsageResetPassword,
          ExpirationTime: time.Now().Add(15 * time.Minute),
        },
      },
      wantErr: true,
    },
    {
      name: "no user id",
      args: args{
        verif: &users.Verification{
          Token:          "Svx79uIDYN7NzDFU2EVWgxo3RX",
          UserId:         uuid.NewString(),
          Usage:          0,
          ExpirationTime: time.Now().Add(15 * time.Minute),
        },
      },
      wantErr: true,
    },
    {
      name: "duplicate token",
      args: args{
        verif: &users.Verification{
          Token:          "7EDOK9dOdtRZkRMUTYnaQcu5bn",
          UserId:         "f8b6a114-8cfe-4e14-b2fb-590c53cec0f1",
          Usage:          0,
          ExpirationTime: time.Now().Add(15 * time.Minute),
        },
      },
      wantErr: true,
    },
    {
      name: "duplicate verification usage",
      args: args{
        verif: &users.Verification{
          Token:          "Svx79uIDYN7NzDFU2EVWgxo3RX",
          UserId:         "4afa29b2-d543-4489-b8ef-93f57781c9f6",
          Usage:          1,
          ExpirationTime: time.Now().Add(15 * time.Minute),
        },
      },
      wantErr: true,
    },
    {
      name: "double verification with different usage",
      args: args{
        verif: &users.Verification{
          Token:          "Svx79uIDYN7NzDFU2EVWgxo3RX",
          UserId:         "4afa29b2-d543-4489-b8ef-93f57781c9f6",
          Usage:          0,
          ExpirationTime: time.Now().Add(15 * time.Minute),
        },
      },
      wantErr: false,
    },
    {
      name: "create new one",
      args: args{
        verif: &users.Verification{
          Token:          "Svx79uIDYN7NzDFU2EVWgxo3RX",
          UserId:         "f8b6a114-8cfe-4e14-b2fb-590c53cec0f1",
          Usage:          1,
          ExpirationTime: time.Now().Add(15 * time.Minute),
        },
      },
      wantErr: false,
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      tx, err := Db.Begin()
      require.NoError(t, err)
      defer tx.Rollback()

      v := NewVerification(tx)
      if err := v.Upsert(tt.args.verif); (err != nil) != tt.wantErr {
        t.Errorf("Upsert() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func Test_verificationRepository_Find(t *testing.T) {
  type args struct {
    token string
  }
  tests := []struct {
    name    string
    args    args
    want    *users.Verification
    wantErr bool
  }{
    {
      name: "verification not found",
      args: args{
        token: "UTJBrHGf210pj02S5L9L6kqNUb",
      },
      want:    nil,
      wantErr: true,
    },
    {
      name: "verification found",
      args: args{
        token: "7EDOK9dOdtRZkRMUTYnaQcu5bn",
      },
      want: &users.Verification{
        Token:          "7EDOK9dOdtRZkRMUTYnaQcu5bn",
        UserId:         "c7760836-71e7-4664-99e8-a9503482a296",
        Usage:          users.UsageResetPassword,
        ExpirationTime: time.Now(),
      },
      wantErr: false,
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      tx, err := Db.Begin()
      require.NoError(t, err)
      defer tx.Rollback()

      v := NewVerification(tx)
      got, err := v.Find(tt.args.token)
      if (err != nil) != tt.wantErr {
        t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
        return
      }

      if tt.want == nil {
        return
      }

      got.ExpirationTime = tt.want.ExpirationTime

      if !reflect.DeepEqual(&got, tt.want) {
        t.Errorf("Find() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_verificationRepository_Remove(t *testing.T) {
  type args struct {
    token string
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "non exists verification",
      args: args{
        token: "UTJBrHGf210pj02S5L9L6kqNUb",
      },
      wantErr: true,
    },
    {
      name: "exists verification",
      args: args{
        token: "c4RPI8LTQ3nfWeVrz4I6KY4V8B",
      },
      wantErr: false,
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      tx, err := Db.Begin()
      require.NoError(t, err)
      defer tx.Rollback()

      v := NewVerification(tx)
      if err := v.Remove(tt.args.token); (err != nil) != tt.wantErr {
        t.Errorf("Remove() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}
