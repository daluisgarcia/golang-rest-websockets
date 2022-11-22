package middleware

import (
	"net/http"
	"strings"

	"github.com/daluisgarcia/golang-rest-websockets/models"
	"github.com/daluisgarcia/golang-rest-websockets/server"
	"github.com/golang-jwt/jwt/v4"
)

var NO_AUTH_NEEDED = []string{"login", "signup"}

func shouldCheckAuth(path string) bool {
	for _, p := range NO_AUTH_NEEDED {
		if strings.Contains(path, p) {
			return false
		}
	}
	return true
}

func GetJwtTokenFromHeader(s server.Server, r *http.Request) (*jwt.Token, error) {
	tokenString := strings.TrimSpace(r.Header.Get("Authorization"))

	return jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.Config().JWTSecret), nil
	})
}

func CheckAuthMiddleware(s server.Server) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !shouldCheckAuth(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			_, err := GetJwtTokenFromHeader(s, r)

			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
