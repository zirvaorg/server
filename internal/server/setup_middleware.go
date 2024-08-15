package server

import (
	"net/http"
	"server/middleware"
)

func SetupMiddleware(handler http.Handler) http.Handler {
	panicHandler := middleware.PanicMiddleware(handler)
	corsHandler := middleware.CorsMiddleware(panicHandler)
	return corsHandler
}
