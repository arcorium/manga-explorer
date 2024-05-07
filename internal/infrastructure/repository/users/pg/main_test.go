package pg

import (
  "github.com/google/uuid"
  "log"
  "manga-explorer/internal/domain/users"
  "manga-explorer/internal/util"
  "testing"

  "github.com/uptrace/bun"
  "manga-explorer/database"
  "manga-explorer/internal/app/common"
)

var Conf common.Config
var Db *bun.DB

var Credentials = []users.Credential{
  users.NewCredential2("c7760836-71e7-4664-99e8-a9503482a296", "TEST", uuid.NewString(), util.GenerateRandomString(30)),
  users.NewCredential2("c7760836-71e7-4664-99e8-a9503482a296", "TEST", uuid.NewString(), util.GenerateRandomString(30)),
  users.NewCredential2("c7760836-71e7-4664-99e8-a9503482a296", "TEST", uuid.NewString(), util.GenerateRandomString(30)),
}

func LoadFixtures() bool {
  tx, err := Db.Begin()
  if err != nil {
    return false
  }
  repo := NewCredential(tx)

  for _, cred := range Credentials {
    err = repo.Create(&cred)
    if err != nil {
      tx.Rollback()
      return false
    }
  }
  tx.Commit()
  return true
}

func RemoveFixtures() {
  repo := NewCredential(Db)

  for _, cred := range Credentials {
    err := repo.Remove(cred.UserId, cred.Id)
    if err != nil {
      log.Println(err)
    }
  }
}

func TestMain(m *testing.M) {
  var err error
  Conf, err = common.LoadConfig("test", "../../../../../")
  if err != nil {
    panic(err)
  }
  Db = database.Open(&Conf, false)
  defer database.Close(Db)

  // Add fixtures
  if !LoadFixtures() {
    panic("Error on loading credentials fixtures")
  }

  m.Run()

  RemoveFixtures()
}
