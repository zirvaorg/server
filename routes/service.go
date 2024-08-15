package routes

import (
	"fmt"
	"net/http"
	"server/middleware"
)

func Service(mux *http.ServeMux) {
	mux.HandleFunc("GET /service/{op}", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		op := r.PathValue("op")

		fmt.Println(op)
	}))
}
