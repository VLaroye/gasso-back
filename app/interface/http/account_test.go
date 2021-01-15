package http

import (
	mock_usecase "github.com/VLaroye/gasso-back/app/usecase"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListAccounts(t *testing.T) {
	req, err := http.NewRequest("GET", "/accounts", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	ctrl := gomock.NewController(t)
	accountUsecase := mock_usecase.NewMockAccountUsecase(ctrl)
	accountUsecase.EXPECT().ListAccounts().Return(nil, nil)

	accountService := NewAccountService(accountUsecase)

	handler := http.HandlerFunc(accountService.ListAccounts)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status: got %v want %v",
			status, http.StatusOK)
	}
}
