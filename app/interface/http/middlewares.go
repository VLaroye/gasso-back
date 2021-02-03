package http

import (
	"github.com/VLaroye/gasso-back/app/interface/http/response"
	"github.com/VLaroye/gasso-back/app/interface/http/utils"
	"net/http"
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
	return func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("token")
		if err != nil {
			response.JSON(
				w,
				http.StatusBadRequest,
				response.ErrorResponse{Message: "error getting 'token' cookie from request", Status: http.StatusBadRequest},
			)
			return
		}

		isTokenValid, err := utils.ValidateJWTToken(tokenCookie.Value)
		if err != nil || !isTokenValid {
			response.JSON(
				w,
				http.StatusUnauthorized,
				response.ErrorResponse{Message: "invalid token", Status: http.StatusUnauthorized},
			)
			return
		}

		next.ServeHTTP(w, r)
	}
}
