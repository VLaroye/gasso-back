package main

import (
	"github.com/VLaroye/gasso-back/app/interface/db"
	"log"
	"net/http"

	"github.com/VLaroye/gasso-back/app/domain/service"
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
	userUsecase := initUserUsecase(database)
	httpUserService := httpInterface.NewUserService(userUsecase)

	// Init account repo + service + usecase
	accountUsecase := initAccountUsecase(database, logger)
	httpAccountService := httpInterface.NewAccountService(accountUsecase)

	// Init invoice repo + service + usecase
	invoiceUsecase := initInvoiceUsecase(database, logger)
	httpInvoiceService := httpInterface.NewInvoiceService(invoiceUsecase, accountUsecase)

	// Handlers
	httpInterface.RegisterUserHandlers(router, httpUserService)
	httpInterface.RegisterAccountHandlers(router, httpAccountService)
	httpInterface.RegisterInvoiceHandlers(router, httpInvoiceService)

	// Run server
	log.Fatal(http.ListenAndServe(":5000", router))
}

func initUserUsecase(database *gorm.DB) usecase.UserUsecase {
	userRepo := db.NewUserRepository(database)
	userService := service.NewUserService(userRepo)

	return usecase.NewUserUsecase(userRepo, userService)
}

func initAccountUsecase(database *gorm.DB, logger *zap.SugaredLogger) usecase.AccountUsecase {
	accountRepo := db.NewAccountRepository(database, logger)
	accountService := service.NewAccountService(accountRepo)

	return usecase.NewAccountUsecase(accountRepo, accountService)
}

func initInvoiceUsecase(database *gorm.DB, logger *zap.SugaredLogger) usecase.InvoiceUsecase {
	invoiceRepo := db.NewInvoiceRepository(database, logger)
	invoiceService := service.NewInvoiceService(invoiceRepo)

	return usecase.NewInvoiceUsecase(invoiceService, invoiceRepo)
}
