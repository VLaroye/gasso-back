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
}

type Account struct {
	ID   string `json:"id"`
	Name string `json:"email"`
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

func (u *accountService) ListAccounts(w http.ResponseWriter, r *http.Request) {
	type accountResponse struct {
		Accounts []*Account `json:"accounts"`
		Error    error      `json:"error"`
	}

	accounts, err := u.accountUsecase.ListAccounts()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(accountResponse{
		Accounts: toAccounts(accounts),
		Error:    nil,
	})

	return
}

func (u *accountService) CreateAccount(w http.ResponseWriter, r *http.Request) {
	type accountRequest struct {
		Name string `json: name`
	}

	var request accountRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if request.Name == "" {
		http.Error(w, "missing required 'name'", http.StatusBadRequest)
		return
	}

	if err := u.accountUsecase.CreateAccount(request.Name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
