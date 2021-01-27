package main

import (
	"github.com/VLaroye/gasso-back/app/interface/http"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	port := 5000
	database, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	http.Start(port, database)
}