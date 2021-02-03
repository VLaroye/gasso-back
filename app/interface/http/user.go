package http

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/VLaroye/gasso-back/app/domain/model"
	"github.com/VLaroye/gasso-back/app/usecase"
	"github.com/gorilla/mux"
)

func RegisterUserHandlers(router *mux.Router, service *userService) {
	router.HandleFunc("/signin", service.Login).Methods("POST")
	router.HandleFunc("/signup", service.Register).Methods("POST")
	router.HandleFunc("/users", AuthenticationNeeded(service.ListUsers)).Methods("GET")
}

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type JWTClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
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
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := u.userUsecase.Login(request.Email, request.Password); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	jwtExpiracyDelay, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION_DELAY"))
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	tokenExpirationTime := time.Now().Add(time.Duration(jwtExpiracyDelay) * time.Minute)

	claims := &JWTClaims{
		Email: request.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpirationTime.Unix(),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jwtToken.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		fmt.Println(err, os.Getenv("JWT_KEY"))
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: tokenExpirationTime,
	})

	respondJSON(w, nil)
}

func (u *userService) Register(w http.ResponseWriter, r *http.Request) {
	type registerUserRequest struct {
		Email    string
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
