package http

import (
	"fmt"
	"github.com/VLaroye/gasso-back/app/domain/repository"
	"github.com/VLaroye/gasso-back/app/domain/service"
	"github.com/VLaroye/gasso-back/app/interface/db"
	"github.com/VLaroye/gasso-back/app/usecase"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Repository struct {
	UserRepo    repository.UserRepository
	AccountRepo repository.AccountRepository
	InvoiceRepo repository.InvoiceRepository
}

type Service struct {
	UserService    *service.UserService
	AccountService *service.AccountService
	InvoiceService *service.InvoiceService
}

type Usecase struct {
	UserUsecase    usecase.UserUsecase
	AccountUsecase usecase.AccountUsecase
	InvoiceUsecase usecase.InvoiceUsecase
}

func Start(port int, database *gorm.DB) {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("failed to launch logger")
	}

	sugaredLogger := logger.Sugar()
	router := mux.NewRouter()

	router.Use(Logging(sugaredLogger))

	repo := initRepository(database, sugaredLogger)
	sv := initService(repo)
	uc := initUsecase(repo, sv)

	httpUserService := NewUserService(uc.UserUsecase)
	httpAccountService := NewAccountService(uc.AccountUsecase)
	httpInvoiceService := NewInvoiceService(uc.InvoiceUsecase, uc.AccountUsecase)

	// Handlers
	RegisterUserHandlers(router, httpUserService)
	RegisterAccountHandlers(router, httpAccountService)
	RegisterInvoiceHandlers(router, httpInvoiceService)

	// Run server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), router))
}

func initRepository(database *gorm.DB, logger *zap.SugaredLogger) *Repository {
	userRepo := db.NewUserRepository(database)
	accountRepo := db.NewAccountRepository(database, logger)
	invoiceRepo := db.NewInvoiceRepository(database, logger)

	return &Repository{
		UserRepo:    userRepo,
		AccountRepo: accountRepo,
		InvoiceRepo: invoiceRepo,
	}
}

func initService(repo *Repository) *Service {
	userService := service.NewUserService(repo.UserRepo)
	accountService := service.NewAccountService(repo.AccountRepo, repo.InvoiceRepo)
	invoiceService := service.NewInvoiceService(repo.InvoiceRepo)

	return &Service{
		UserService:    userService,
		AccountService: accountService,
		InvoiceService: invoiceService,
	}
}

func initUsecase(repo *Repository, service *Service) *Usecase {
	return &Usecase{
		UserUsecase:    usecase.NewUserUsecase(repo.UserRepo, service.UserService),
		AccountUsecase: usecase.NewAccountUsecase(repo.AccountRepo, service.AccountService),
		InvoiceUsecase: usecase.NewInvoiceUsecase(repo.InvoiceRepo, service.InvoiceService),
	}
}
