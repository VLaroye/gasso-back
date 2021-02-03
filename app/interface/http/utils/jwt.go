package utils

import (
	"github.com/dgrijalva/jwt-go"
	"os"
	"strconv"
	"time"
)

type JWToken struct {
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

func NewJWToken(email string) (*JWToken, error) {
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

	return &JWToken{
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

	return jwtToken.Valid, err
}

