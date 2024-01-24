package database

import (
	"context"
	"database/sql"
	"github.com/uptrace/bun/dbfixture"
	"log"
	"manga-explorer/database/fixtures"
	"manga-explorer/internal/domain/mangas"
	"manga-explorer/internal/domain/users"
	"manga-explorer/internal/util"
	"os"
	"slices"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"manga-explorer/internal/app/common"
)

var tables = []interface{}{
	(*users.User)(nil),
	(*users.Profile)(nil),
	(*users.Credential)(nil),
	(*users.Verification)(nil),
	(*mangas.Manga)(nil),
	(*mangas.Volume)(nil),
	(*mangas.MangaFavorite)(nil),
	(*mangas.Chapter)(nil),
	(*mangas.Page)(nil),
	(*mangas.Rate)(nil),
	(*mangas.Genre)(nil),
	(*mangas.Rate)(nil),
	(*mangas.Comment)(nil),
	(*mangas.MangaGenre)(nil),
	(*mangas.Translation)(nil),
	(*mangas.ChapterHistory)(nil),
}

func addDebugLog(db *bun.DB) {
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
}
func Migrate(db *bun.DB) error {
	ctx := context.Background()

	// Registering many-to-many model
	db.RegisterModel(util.Nil[mangas.MangaGenre]())

	for _, table := range tables {
		_, err := db.NewCreateTable().
			Model(table).
			IfNotExists().
			WithForeignKeys().
			Exec(ctx)

		if err != nil {
			return err
		}
	}
	return nil
}

func RegisterModels(db *bun.DB) {
	db.RegisterModel(util.Nil[mangas.MangaGenre]())
	for _, model := range tables {
		db.RegisterModel(model)
	}
}

func Drops(db *bun.DB) []error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	temp := slices.Clone(tables)
	slices.Reverse(temp)
	var errors []error

	for _, table := range temp {
		_, err := db.NewDropTable().
			Model(table).
			Exec(ctx)

		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func Open(config *common.Config, log bool) (*bun.DB, error) {
	dsn := config.DatabaseDSN()
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	if err := sqldb.Ping(); err != nil {
		return nil, err
	}

	db := bun.NewDB(sqldb, pgdialect.New())
	if log {
		addDebugLog(db)
	}
	RegisterModels(db)
	return db, nil
}

func Close(db *bun.DB) {
	if err := db.Close(); err != nil {
		log.Fatalln(err)
	}
}

var fixtureMaps = map[fixtures.Type]string{
	fixtures.UserType:  "user.yaml",
	fixtures.MangaType: "manga.yaml",
}

func LoadFixtures(db *bun.DB, path string, types ...fixtures.Type) error {
	fixtures := dbfixture.New(db)
	for _, tp := range types {
		dir := os.DirFS(path)
		err := fixtures.Load(context.Background(), dir, fixtureMaps[tp])
		if err != nil {
			return err
		}
	}
	return nil
}
