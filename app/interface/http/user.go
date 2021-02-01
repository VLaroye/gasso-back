package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/VLaroye/gasso-back/app/domain/model"
	"github.com/VLaroye/gasso-back/app/usecase"
	"github.com/gorilla/mux"
)

func RegisterUserHandlers(router *mux.Router, service *userService) {
	router.HandleFunc("/signin", service.Login).Methods("POST")
	router.HandleFunc("/signup", service.RegisterUser).Methods("POST")
	router.HandleFunc("/users", service.ListUsers).Methods("GET")
}

type User struct {
	ID    string `json:"id"`
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

func (u *userService) Login(w http.ResponseWriter, r *http.Request) {
	type registerUserRequest struct {
		Email string
		Password string
	}

	var request registerUserRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := u.userUsecase.Login(request.Email, request.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Logged in")

	respondJSON(w, nil)
}

func (u *userService) ListUsers(w http.ResponseWriter, r *http.Request) {
	type listUsersResponseType struct {
		Error error   `json:"error"`
		Users []*User `json:"users"`
	}

	users, err := u.userUsecase.ListUsers()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	_ = json.NewEncoder(w).Encode(listUsersResponseType{
		Error: nil,
		Users: toUsers(users),
	})
	return
}

func (u *userService) RegisterUser(w http.ResponseWriter, r *http.Request) {
	type registerUserRequest struct {
		Email string
		Password string
	}

	var request registerUserRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := u.userUsecase.RegisterUser(request.Email, request.Password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
