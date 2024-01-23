package main

import (
	"log"
	"manga-explorer/database"
	"manga-explorer/database/fixtures"
	"manga-explorer/internal/app/common"
)

func main() {
	config, err := common.LoadConfig("test")
	if err != nil {
		log.Fatalln(err)
	}

	db, err := database.Open(config, true)
	if err != nil {
		log.Fatalln("Failed to open database connection: ", err)
	}
	defer database.Close(db)
	database.Drops(db)

	database.RegisterModels(db)
	err = database.Migrate(db)
	if err != nil {
		database.Drops(db)
		log.Fatalln("Failed to migrate database: ", err)
	}

	err = database.InsertSpecialRecords(db)
	if err != nil {
		database.Drops(db)
		log.Fatalln("Failed to insert special records: ", err)
	}

	err = database.LoadFixtures(db, "./database/fixtures", fixtures.UserType, fixtures.MangaType)
	if err != nil {
		database.Drops(db)
		log.Fatalln(err)
	}
}
