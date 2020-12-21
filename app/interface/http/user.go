package http

import (
	"encoding/json"
	"github.com/VLaroye/gasso-back/app/domain/model"
	"github.com/VLaroye/gasso-back/app/usecase"
	"github.com/gorilla/mux"
	"net/http"
)

type User struct {
	ID string `json:"id"`
	Email string `json:"email"`
}

type userService struct {
	userUsecase usecase.UserUsecase
}

func NewUserService(userUsecase usecase.UserUsecase) *userService {
	return &userService{userUsecase: userUsecase}
}

func toUsers(users []*model.User) []*User {
	result := make([]*User, len(users))

	for i, user := range users {
		result[i] = &User{ID: user.GetId(), Email: user.GetEmail()}
	}
	return result
}

func (u *userService) RegisterHandlers(router *mux.Router) {
	router.HandleFunc("/users", u.ListUsers)
}

func (u *userService) ListUsers(w http.ResponseWriter, r *http.Request) {
	type listUsersResponseType struct {
		Error error `json:"error"`
		Users []*User `json:"users"`
	}
	
	users, err := u.userUsecase.ListUsers()

	if err != nil {
		// TODO: Handle this error
		_ = json.NewEncoder(w).Encode(listUsersResponseType{
			Error: err,
			Users: nil,
		})
		return
	}

	_ = json.NewEncoder(w).Encode(listUsersResponseType{
		Error: nil,
		Users: toUsers(users),
	})
	return
}

func (u *userService) RegisterUser(email string) error {
	return nil
}