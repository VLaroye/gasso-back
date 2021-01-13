package http

import (
	"encoding/json"
	"net/http"

	"github.com/VLaroye/gasso-back/app/domain/model"
	"github.com/VLaroye/gasso-back/app/usecase"
	"github.com/gorilla/mux"
)

func RegisterAccountHandlers(router *mux.Router, service *accountService) {
	router.HandleFunc("/accounts", service.ListAccounts).Methods("GET")
	router.HandleFunc("/accounts", service.CreateAccount).Methods("POST")
	router.HandleFunc("/accounts/{id}", service.GetAccountByID).Methods("GET")
	router.HandleFunc("/accounts/{id}", service.UpdateAccount).Methods("POST")
	router.HandleFunc("/accounts/{id}", service.DeleteAccount).Methods("DELETE")
}

type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if account == nil {
		respondError(w, http.StatusNotFound, "account not found")
		return
	}

	respondJSON(w, &Account{ID: account.GetId(), Name: account.GetName()})
	return
}

func (u *accountService) ListAccounts(w http.ResponseWriter, r *http.Request) {
	type accountResponse struct {
		Accounts []*Account `json:"accounts"`
	}

	accounts, err := u.accountUsecase.ListAccounts()

	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := accountResponse{
		Accounts: toAccounts(accounts),
	}

	respondJSON(w, response)
	return
}

func (u *accountService) CreateAccount(w http.ResponseWriter, r *http.Request) {
	type accountRequest struct {
		Name string `json: name`
	}

	var request accountRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if request.Name == "" {
		respondError(w, http.StatusBadRequest, "missing required 'name'")
		return
	}

	if err := u.accountUsecase.CreateAccount(request.Name); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := u.accountUsecase.UpdateAccount(accountID, request.Name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (u *accountService) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID := vars["id"]

	if err := u.accountUsecase.DeleteAccount(accountID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
