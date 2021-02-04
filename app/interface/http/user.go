package http

import (
	"encoding/json"
	"github.com/VLaroye/gasso-back/app/domain/model"
	"github.com/VLaroye/gasso-back/app/interface/http/response"
	"github.com/VLaroye/gasso-back/app/interface/http/utils"
	"github.com/VLaroye/gasso-back/app/usecase"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterUserHandlers(router *mux.Router, service *userService) {
	router.HandleFunc("/signin", service.Login).Methods("POST")
	router.HandleFunc("/signup", service.Register).Methods("POST")
	router.HandleFunc("/refresh", service.RefreshToken).Methods("GET")
	router.HandleFunc("/users", AuthenticationNeeded(service.ListUsers)).Methods("GET")
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
	type signInRequest struct {
		Email    string
		Password string
	}

	var request signInRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.JSON(
			w,
			http.StatusInternalServerError,
			response.ErrorResponse{Message: err.Error(), Status: http.StatusInternalServerError},
		)
		return
	}

	if request.Email == "" || request.Password == "" {
		response.JSON(
			w,
			http.StatusBadRequest,
			response.ErrorResponse{Message: "email and password fields are required", Status: http.StatusBadRequest},
		)
		return
	}

	if err := u.userUsecase.Login(request.Email, request.Password); err != nil {
		response.JSON(
			w,
			http.StatusUnauthorized,
			response.ErrorResponse{Message: "invalid email or password", Status: http.StatusUnauthorized},
		)
		return
	}

	token, err := utils.NewJWToken(request.Email)
	if err != nil {
		response.JSON(
			w,
			http.StatusInternalServerError,
			response.ErrorResponse{Message: "error generating auth token", Status: http.StatusInternalServerError},
		)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token.SignedString,
		Expires: token.Expires,
	})
}

func (u *userService) Register(w http.ResponseWriter, r *http.Request) {
	type registerUserRequest struct {
		Email    string
		Password string
	}

	var request registerUserRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.JSON(
			w,
			http.StatusBadRequest,
			response.ErrorResponse{Message: "error parsing request", Status: http.StatusBadRequest},
		)
		return
	}

	if request.Email == "" || request.Password == "" {
		response.JSON(
			w,
			http.StatusBadRequest,
			response.ErrorResponse{Message: "email and password fields are required", Status: http.StatusBadRequest},
		)
		return
	}

	if err := u.userUsecase.RegisterUser(request.Email, request.Password); err != nil {
		response.JSON(
			w,
			http.StatusInternalServerError,
			response.ErrorResponse{Message: "error registering user", Status: http.StatusInternalServerError},
		)
	}
}

func (u *userService) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// 1 - GET token from request
	tokenCookie, err := r.Cookie("token")
	if err != nil {
		response.JSON(
			w,
			http.StatusBadRequest,
			response.ErrorResponse{Message: "error getting 'token' cookie from request", Status: http.StatusBadRequest},
		)
		return
	}

	tokenValue := tokenCookie.Value

	// 2 - Validate token
	token, err := utils.ParseJWTToken(tokenValue)
	if err != nil {
		response.JSON(
			w,
			http.StatusBadRequest,
			response.ErrorResponse{Message: "invalid token", Status: http.StatusBadRequest},
		)
		return
	}

	// 3 - Create a new token
	newToken, err := utils.NewJWToken(token.Claims.Email)
	if err != nil {
		response.JSON(
			w,
			http.StatusInternalServerError,
			response.ErrorResponse{Message: "error creating new token", Status: http.StatusInternalServerError},
		)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   newToken.SignedString,
		Expires: newToken.Expires,
	})
}

func (u *userService) ListUsers(w http.ResponseWriter, _ *http.Request) {
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
