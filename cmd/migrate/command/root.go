package command

import (
  "database/sql"
  "fmt"
  "github.com/spf13/cobra"
  "github.com/uptrace/bun"
  "github.com/uptrace/bun/dialect/pgdialect"
  "github.com/uptrace/bun/driver/pgdriver"
  "golang.org/x/net/context"
  "log"
  "manga-explorer/database"
  "manga-explorer/database/fixtures"
  "manga-explorer/internal/common"
  "manga-explorer/internal/domain/users"
  "manga-explorer/internal/infrastructure/file"
  "manga-explorer/internal/util"
  "os"
)

var rootCmd = &cobra.Command{
  Use:     "migrate ",
  Example: "migrate -u user -p 123 --host localhost:12345 --database test --special --no-ssl",
  Run:     runRoot,
}

func init() {
  rootCmd.Flags().StringP("user", "u", "", "database username")
  rootCmd.Flags().StringP("pass", "p", "", "database password")
  rootCmd.Flags().String("host", "", "database host url")
  rootCmd.Flags().StringP("database", "d", "manga_explorer", "database name")
  rootCmd.Flags().String("dsn", "", "database dsn")
  rootCmd.Flags().StringP("seed", "s", "", "seed database without migrating")
  rootCmd.Flags().Bool("special", false, "add special record (admin) from env")
  rootCmd.Flags().Bool("no-ssl", false, "Disable SSL on database communication")
  rootCmd.Flags().Bool("env", false, "Use environment variables for database connection")

  rootCmd.MarkFlagsRequiredTogether("user", "pass", "host")
  rootCmd.MarkFlagsMutuallyExclusive("user", "dsn", "env")
  rootCmd.MarkFlagsOneRequired("user", "dsn", "env")
  rootCmd.MarkFlagsMutuallyExclusive("seed", "special")
}

func runRoot(command *cobra.Command, args []string) {
  dsn, _ := command.Flags().GetString("dsn")
  if len(dsn) == 0 {
    env, _ := command.Flags().GetBool("env")
    if env {
      config, err := common.LoadConfig()
      if err != nil {
        panic(err)
      }
      dsn = config.DatabaseDSN()
    } else {
      username, _ := command.Flags().GetString("user")
      pass, _ := command.Flags().GetString("pass")
      host, _ := command.Flags().GetString("host")
      dbName, _ := command.Flags().GetString("database")
      noSSL, _ := command.Flags().GetBool("no-ssl")
      dsn = fmt.Sprintf("postgres://%s:%s@%s/%s", username, pass, host, dbName)
      if noSSL {
        dsn += "?sslmode=disable"
      }
    }
  }

  //log.Println("Open database connection with DSN: ", dsn)

  db, err := openDb(dsn)
  if err != nil {
    log.Fatalln("Failed to open database: ", err)
  }
  defer database.Close(db)
  // Seeding
  seedPath, _ := command.Flags().GetString("seed")
  if len(seedPath) != 0 {
    database.RegisterModels(db)
    err = database.LoadFixtures(db, seedPath, fixtures.UserType, fixtures.MangaType)
    if err != nil {
      log.Fatalln(err)
    }

    fmt.Println("Success!")
    return
  }

  // Migrate Database
  err = migrate(db)
  if err != nil {
    log.Fatalln("Failed to migrate database: ", err)
  }

  // Insert special record
  special, _ := command.Flags().GetBool("special")
  if special {
    adminUsername := os.Getenv("ME_ADMIN_USERNAME")
    adminPass := os.Getenv("ME_ADMIN_PASSWORD")
    adminEmail := os.Getenv("ME_ADMIN_EMAIL")

    util.SetDefaultString(&adminUsername, "admin")
    util.SetDefaultString(&adminPass, "admin123")
    util.SetDefaultString(&adminEmail, "admin@manga-explorer.com")

    user, err := users.NewUser(adminUsername, adminEmail, adminPass, users.RoleAdmin)
    if err != nil {
      log.Fatalln(err)
    }
    profile := users.NewProfile(user.Id, "admin", "", "", file.NullName)

    err = insertSpecialRecord(db, &user, &profile)
    if err != nil {
      log.Fatalln(err)
    }
  }

  fmt.Println("Success!")
}

func Execute() error {
  return rootCmd.Execute()
}

func openDb(dsn string) (*bun.DB, error) {
  sqlDb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
  if err := sqlDb.Ping(); err != nil {
    return nil, err
  }

  db := bun.NewDB(sqlDb, pgdialect.New())

  //db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
  return db, nil
}
func migrate(db *bun.DB) error {
  return database.Migrate(db)
}

func insertSpecialRecord(db *bun.DB, user *users.User, profile *users.Profile) error {
  ctx := context.Background()

  tx, err := db.Begin()
  if err != nil {
    return err
  }

  _, err = tx.NewInsert().
    Model(user).
    Exec(ctx)

  if err != nil {
    return tx.Rollback()
  }

  _, err = tx.NewInsert().
    Model(profile).
    Exec(ctx)

  if err != nil {
    return tx.Rollback()
  }
  return tx.Commit()
}
