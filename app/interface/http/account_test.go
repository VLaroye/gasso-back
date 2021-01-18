package http

import (
	"encoding/json"
	"github.com/VLaroye/gasso-back/app/domain/model"
	mockUsecase "github.com/VLaroye/gasso-back/app/usecase"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

var mockAccount *model.Account = model.NewAccount("123456", "example1")

func TestListAccounts(t *testing.T) {
	type listAccountsResponse struct {
		Accounts []*Account
	}
	var response listAccountsResponse

	req, err := http.NewRequest("GET", "/accounts", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	ctrl := gomock.NewController(t)
	accountUsecase := mockUsecase.NewMockAccountUsecase(ctrl)
	accountUsecase.EXPECT().ListAccounts().Return([]*model.Account{
		mockAccount,
	}, nil)

	accountService := NewAccountService(accountUsecase)

	handler := http.HandlerFunc(accountService.ListAccounts)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status: got %v want %v",
			status, http.StatusOK)
	}

	json.NewDecoder(rr.Body).Decode(&response)

	if len(response.Accounts) != 1 {
		t.Errorf("wrong number of accounts: got %v want %v",
			len(response.Accounts), 1)
	}
}
