package main

import (
	"fmt"
	"github.com/VLaroye/gasso-back/app/interface/db"
	"github.com/VLaroye/gasso-back/app/interface/http"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// TODO: Handle migrations on a subcommand ? Here only for dev (?)
	if err := database.AutoMigrate(&db.User{}, &db.Account{}, &db.Invoice{}); err != nil {
		panic(err)
	}

	http.Start(5000, database)
}
