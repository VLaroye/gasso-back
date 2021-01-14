package main

import (
	"log"
	"net/http"

	"github.com/VLaroye/gasso-back/app/domain/service"
	"github.com/VLaroye/gasso-back/app/interface/db"
	httpInterface "github.com/VLaroye/gasso-back/app/interface/http"
	"github.com/VLaroye/gasso-back/app/usecase"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	router := mux.NewRouter()
	database, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	zap, err := zap.NewProduction()
	if err != nil {
		panic("failed to launch logger")
	}

	logger := zap.Sugar()

	router.Use(httpInterface.LoggingMiddleware(logger))

	// Init user repo + service + usecase
	userRepo := db.NewUserRepository(database)
	userService := service.NewUserService(userRepo)
	userUsecase := usecase.NewUserUsecase(userRepo, userService)
	httpUserService := httpInterface.NewUserService(userUsecase)

	// Init account repo + service + usecase
	accountRepo := db.NewAccountRepository(database, logger)
	accountService := service.NewAccountService(accountRepo)
	accountUsecase := usecase.NewAccountUsecase(accountRepo, accountService)
	httpAccountService := httpInterface.NewAccountService(accountUsecase)

	// Handlers
	httpInterface.RegisterUserHandlers(router, httpUserService)
	httpInterface.RegisterAccountHandlers(router, httpAccountService)

	// Run server
	log.Fatal(http.ListenAndServe(":5000", router))
}
