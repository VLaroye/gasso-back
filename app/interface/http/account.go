package http

import (
	"encoding/json"
	"github.com/VLaroye/gasso-back/app/interface/http/response"
	"net/http"

	"github.com/VLaroye/gasso-back/app/domain/model"
	"github.com/VLaroye/gasso-back/app/usecase"
	"github.com/gorilla/mux"
)

func RegisterAccountHandlers(router *mux.Router, service *accountService) {
	router.HandleFunc("/accounts", AuthenticationNeeded(service.ListAccounts)).Methods("GET")
	router.HandleFunc("/accounts", AuthenticationNeeded(service.CreateAccount)).Methods("POST")
	router.HandleFunc("/accounts/{id}", AuthenticationNeeded(service.GetAccountByID)).Methods("GET")
	router.HandleFunc("/accounts/{id}", AuthenticationNeeded(service.UpdateAccount)).Methods("PUT")
	router.HandleFunc("/accounts/{id}", AuthenticationNeeded(service.DeleteAccount)).Methods("DELETE")
}

type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewAccount(id, name string) *Account {
	return &Account{ID: id, Name: name}
}

type accountService struct {
	accountUsecase usecase.AccountUsecase
}

func NewAccountService(accountUsecase usecase.AccountUsecase) *accountService {
	return &accountService{accountUsecase: accountUsecase}
}

func toAccounts(accounts []*model.Account) []*Account {
	result := make([]*Account, len(accounts))

	for i, account := range accounts {
		result[i] = &Account{ID: account.GetId(), Name: account.GetName()}
	}
	return result
}

func (u *accountService) GetAccountByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	account, err := u.accountUsecase.GetAccountByID(id)

	if err != nil {
		response.JSON(
			w,
			http.StatusInternalServerError,
			response.ErrorResponse{Message: "error getting account from db", Status: http.StatusInternalServerError},
		)
		return
	}

	if account == nil {
		response.JSON(
			w,
			http.StatusNotFound,
			response.ErrorResponse{Message: "account not found", Status: http.StatusNotFound},
		)
		return
	}

	response.JSON(w, http.StatusOK, &Account{ID: account.GetId(), Name: account.GetName()})
	return
}

func (u *accountService) ListAccounts(w http.ResponseWriter, r *http.Request) {
	type accountResponse struct {
		Accounts []*Account `json:"accounts"`
	}

	accounts, err := u.accountUsecase.ListAccounts()

	if err != nil {
		response.JSON(
			w,
			http.StatusInternalServerError,
			response.ErrorResponse{Message: "error getting accounts from db", Status: http.StatusInternalServerError},
		)
		return
	}

	resp := accountResponse{
		Accounts: toAccounts(accounts),
	}

	response.JSON(w, http.StatusOK, resp)
	return
}

func (u *accountService) CreateAccount(w http.ResponseWriter, r *http.Request) {
	type accountRequest struct {
		Name string `json: name`
	}

	var request accountRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.JSON(
			w,
			http.StatusBadRequest,
			response.ErrorResponse{Message: "error decoding request", Status: http.StatusBadRequest},
		)
		return
	}

	if request.Name == "" {
		response.JSON(
			w,
			http.StatusBadRequest,
			response.ErrorResponse{Message: "missing required 'name'", Status: http.StatusBadRequest},
		)
		return
	}

	if err := u.accountUsecase.CreateAccount(request.Name); err != nil {
		response.JSON(
			w,
			http.StatusInternalServerError,
			response.ErrorResponse{Message: "error creating account", Status: http.StatusInternalServerError},
		)
		return
	}
}

func (u *accountService) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	type accountRequest struct {
		Name string `json:"name"`
	}

	vars := mux.Vars(r)
	accountID := vars["id"]

	var request accountRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.JSON(
			w,
			http.StatusBadRequest,
			response.ErrorResponse{Message: "error decoding request", Status: http.StatusBadRequest},
		)
		return
	}

	if err := u.accountUsecase.UpdateAccount(accountID, request.Name); err != nil {
		response.JSON(
			w,
			http.StatusInternalServerError,
			response.ErrorResponse{Message: "error updating account", Status: http.StatusInternalServerError},
		)
		return
	}
}

func (u *accountService) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID := vars["id"]

	if err := u.accountUsecase.DeleteAccount(accountID); err != nil {
		response.JSON(
			w,
			http.StatusInternalServerError,
			response.ErrorResponse{Message: "error deleting account", Status: http.StatusInternalServerError},
		)
		return
	}
}
