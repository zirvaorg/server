package middleware

import (
	"net/http"
	"server/internal/logic"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login, err := logic.CheckAuth(r)
		if err != nil || !login {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
