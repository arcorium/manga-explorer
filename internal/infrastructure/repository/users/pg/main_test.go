package pg

import (
	"testing"

	"github.com/uptrace/bun"
	"manga-explorer/database"
	"manga-explorer/internal/app/common"
)

var Conf common.Config
var Db *bun.DB

func TestMain(m *testing.M) {
	var err error
	Conf, err = common.LoadConfig("test", "../../../../../")
	if err != nil {
		panic(err)
	}
	Db = database.Open(&Conf, false)
	defer database.Close(Db)

	m.Run()
}
