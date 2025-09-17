package middleware

import (
	"net/http"
	"strings"
)

var jwtSecret = []byte("super-secret")

func AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		token = strings.TrimPrefix(token, "Bearer ")

		if token != string(jwtSecret) {
			http.Error(w, "Unauthorized", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)

	})
}
