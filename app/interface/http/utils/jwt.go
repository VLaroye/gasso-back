package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"os"
	"strconv"
	"time"
)

type JWTToken struct {
	Token *jwt.Token
	SignedString string
	Expires time.Time
	Claims *Claims
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func NewClaims(email string, expirationTime time.Time) *Claims {
	return &Claims{
		Email:          email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
}

func NewJWToken(email string) (*JWTToken, error) {
	jwtExpiryDelay, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION_DELAY"))
	if err != nil {
		return nil, err
	}

	expires := time.Now().Add(time.Duration(jwtExpiryDelay) * time.Minute)

	claims := NewClaims(email, expires)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return nil, err
	}

	return &JWTToken{
		Token:        token,
		SignedString: signedString,
		Expires:      expires,
		Claims:       claims,
	}, nil
}

func ValidateJWTToken(receivedToken string) (bool, error){
	claims := &Claims{}

	jwtToken, err := jwt.ParseWithClaims(receivedToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})
	if err != nil {
		return false, err
	}

	return jwtToken.Valid, err
}

func ParseJWTToken(receivedToken string) (*JWTToken, error){
	claims := &Claims{}

	jwtToken, err := jwt.ParseWithClaims(receivedToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})
	if err != nil {
		return nil, err
	}

	signedString, err := jwtToken.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return nil, err
	}

	if !jwtToken.Valid {
		return nil, errors.New("invalid token")
	}

	token := &JWTToken{
		Token:        jwtToken,
		SignedString: signedString,
		Expires:      time.Unix(claims.ExpiresAt, 0),
		Claims:       claims,
	}

	return token, err
}
