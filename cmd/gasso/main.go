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
	dbt, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	userRepo := db.NewUserRepository(dbt)
	userService := service.NewUserService(userRepo)
	userUsecase := usecase.NewUserUsecase(userRepo, userService)

	httpUserService := http2.NewUserService(userUsecase)
	httpUserService.RegisterHandlers(router)

	log.Fatal(http.ListenAndServe(":5000", router))
}
