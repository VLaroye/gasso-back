package main

import (
	"github.com/VLaroye/gasso-back/app/domain/service"
	"github.com/VLaroye/gasso-back/app/interface/db"
	http2 "github.com/VLaroye/gasso-back/app/interface/http"
	"github.com/VLaroye/gasso-back/app/usecase"
	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	database, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Init user repo + service + usecase
	userRepo := db.NewUserRepository(database)
	userService := service.NewUserService(userRepo)
	userUsecase := usecase.NewUserUsecase(userRepo, userService)
	httpUserService := http2.NewUserService(userUsecase)

	// Handlers
	http2.RegisterUserHandlers(router, httpUserService)

	// Run server
	log.Fatal(http.ListenAndServe(":5000", router))
}
