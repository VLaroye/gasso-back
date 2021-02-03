package http

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
)

func Logging(logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logger.Errorw(
						"err", err,
					)
				}
			}()

			logger.Infow("http request received",
				"method", r.Method,
				"path", r.URL.EscapedPath(),
			)
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.Infow("http response sent",
				"method", r.Method,
				"path", r.URL.EscapedPath(),
				"duration", time.Since(start),
			)
		}

		return http.HandlerFunc(fn)
	}
}

func AuthenticationNeeded(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("token")
		if err != nil {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}

		tokenValue := tokenCookie.Value
		claims := &JWTClaims{}

		jwtToken, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_KEY")), nil
		})
		if err != nil {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}
		if !jwtToken.Valid {
			respondError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		next.ServeHTTP(w, r)
	})
}
